package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	apiURL  = "https://api.anthropic.com/v1/messages"
	version = "2023-06-01"
	// Use Haiku for low-latency moderation; swap to sonnet for higher accuracy.
	defaultModel = "claude-haiku-4-5-20251001"
)

// ClassificationResult is returned by Classify.
type ClassificationResult struct {
	Category   string  `json:"category"`   // "adult_content" | "gambling" | "violence" | "self_harm" | "clean"
	Confidence float64 `json:"confidence"` // 0.0 – 1.0
	Severity   string  `json:"severity"`   // "informational" | "concerning" | "severe"
}

// Client calls the Anthropic Messages API.
// When ANTHROPIC_API_KEY is unset, Classify returns a local keyword fallback
// rather than erroring, so the pipeline degrades gracefully without credentials.
type Client struct {
	apiKey  string
	model   string
	http    *http.Client
}

// New returns a Client. Set ANTHROPIC_API_KEY in the environment to enable
// actual Claude calls; omit it to use the built-in keyword fallback only.
func New() *Client {
	return &Client{
		apiKey: os.Getenv("ANTHROPIC_API_KEY"),
		model:  defaultModel,
		http:   &http.Client{},
	}
}

// Classify sends text to Claude for content moderation and returns the result.
func (c *Client) Classify(ctx context.Context, text string) (*ClassificationResult, error) {
	if c.apiKey == "" {
		return keywordFallback(text), nil
	}
	return c.callAPI(ctx, text)
}

func (c *Client) callAPI(ctx context.Context, text string) (*ClassificationResult, error) {
	systemPrompt := `You are a content moderation engine for a Christian accountability app.
Classify the following screen text into exactly one category.

Categories and their meanings:
- adult_content: pornographic or sexually explicit material
- gambling: online betting, casino games, or sports wagering
- violence: graphic violent content, gore, or harmful imagery
- self_harm: suicidal ideation, self-harm instructions, or eating disorder content
- clean: safe, non-harmful content

Respond with ONLY valid JSON in this exact format (no markdown, no explanation):
{"category":"<category>","confidence":<0.0-1.0>,"severity":"<informational|concerning|severe>"}`

	userMsg := fmt.Sprintf("Classify this screen text:\n\n%s", strings.TrimSpace(text))
	if len(userMsg) > 2000 {
		userMsg = userMsg[:2000]
	}

	payload := map[string]any{
		"model":      c.model,
		"max_tokens": 120,
		"system":     systemPrompt,
		"messages": []map[string]string{
			{"role": "user", "content": userMsg},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("anthropic: marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("anthropic: build request: %w", err)
	}
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", version)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("anthropic: request: %w", err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("anthropic: HTTP %d: %s", resp.StatusCode, respBody)
	}

	// Parse the Messages API response envelope.
	var envelope struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(respBody, &envelope); err != nil {
		return nil, fmt.Errorf("anthropic: parse envelope: %w", err)
	}
	if len(envelope.Content) == 0 {
		return nil, errors.New("anthropic: empty content array")
	}

	// Extract the JSON classification from the text block.
	raw := strings.TrimSpace(envelope.Content[0].Text)
	var result ClassificationResult
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return nil, fmt.Errorf("anthropic: parse result %q: %w", raw, err)
	}

	// Sanitize
	valid := map[string]bool{"adult_content": true, "gambling": true, "violence": true, "self_harm": true, "clean": true}
	if !valid[result.Category] {
		result.Category = "clean"
	}
	if result.Confidence < 0 {
		result.Confidence = 0
	} else if result.Confidence > 1 {
		result.Confidence = 1
	}
	if result.Severity == "" {
		switch {
		case result.Confidence >= 0.75:
			result.Severity = "severe"
		case result.Confidence >= 0.4:
			result.Severity = "concerning"
		default:
			result.Severity = "informational"
		}
	}

	return &result, nil
}

// keywordFallback provides a dependency-free classification when no API key is set.
// Mirrors the FallbackTextClassifier in SampleHandler.swift.
func keywordFallback(text string) *ClassificationResult {
	lower := strings.ToLower(text)

	type rule struct {
		terms    []string
		category string
	}
	rules := []rule{
		{[]string{"porn", "xxx", "nude", "onlyfans", "chaturbate", "erotic", "nsfw"}, "adult_content"},
		{[]string{"casino", "sportsbook", "draftkings", "fanduel", "bet365", "gambling", "wager"}, "gambling"},
		{[]string{"gore", "beheading", "snuff", "murder video", "graphic violence"}, "violence"},
		{[]string{"suicide", "self harm", "self-harm", "kill myself", "suicidal", "overdose on"}, "self_harm"},
	}

	for _, r := range rules {
		count := 0
		for _, t := range r.terms {
			if strings.Contains(lower, t) {
				count++
			}
		}
		if count > 0 {
			conf := min(1.0, 0.40+float64(count)*0.15)
			sev := "informational"
			if conf >= 0.75 {
				sev = "severe"
			} else if conf >= 0.50 {
				sev = "concerning"
			}
			return &ClassificationResult{Category: r.category, Confidence: conf, Severity: sev}
		}
	}
	return &ClassificationResult{Category: "clean", Confidence: 0.95, Severity: "informational"}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
