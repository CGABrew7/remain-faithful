//go:build darwin

package monitor

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"remain-faithful/desktop/core"
)

func Platform() string { return "mac" }

func Ready() bool {
	path := helperPath()
	if path == "" {
		return false
	}
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func Sample() (Front, error) {
	out, err := runHelper("front")
	if err != nil {
		return Front{}, err
	}
	if out == "" || core.IsSelf(Platform(), out) {
		return Front{}, nil
	}
	return Front{Key: out}, nil
}

func Quit(key string) error {
	if !core.AllowQuit(Platform(), key) {
		return errors.New("that app is not chosen for blocking")
	}
	_, err := runHelper("quit", key)
	return err
}

func helperPath() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Join(filepath.Dir(exe), "RFMonitorHelper")
}

func runHelper(args ...string) (string, error) {
	path := helperPath()
	if path == "" {
		return "", errors.New("monitor helper is missing")
	}
	if _, err := os.Stat(path); err != nil {
		return "", errors.New("monitor helper is missing")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path, args...)
	raw, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(raw)), nil
}
