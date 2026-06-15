# Overnight Audit — Remain Faithful
Date: 2026-06-13

## Critical

### Missing PrivacyInfo.xcprivacy for Main App
**File:** `/Users/jbhome/remain-faithful/RemainFaithful/`  
**Issue:** The main app target lacks a PrivacyInfo.xcprivacy file. This is now required by Apple for all apps on the App Store (mandatory since Q1 2024). The app uses:
- ReplayKit (screen recording via broadcast extension)
- SensitiveContentAnalysis (SCA API for image analysis)
- Vision OCR (text recognition from screen)
- Network access (API calls to backend and Anthropic Claude API)
- None of these are declared in a privacy manifest

**Risk:** App Store rejection. Build will fail with required privacy manifest error.  
**Fix:** Create `RemainFaithful/PrivacyInfo.xcprivacy` with declarations for:
- NSPrivacyTracking: false
- NSPrivacyTrackingDomains: []
- NSPrivacyAccessedAPITypes:
  - ReplayKit (for broadcast screen recording)
  - SensitiveContentAnalyzer (for image analysis)
  - Vision OCR (for text recognition)

---

### Unvalidated JWT Claims in Auth Middleware
**File:** `/Users/jbhome/remain-faithful/backend/internal/auth/jwt.go:43-47`  
**Issue:** The `Parse` function validates the token signature and checks `token.Valid` but does **NOT explicitly validate expiry, issuer, or audience claims**. While `token.Valid` implicitly checks expiry, there is no issuer validation, which could allow tokens from compromised or forged sources if the secret is leaked.

**Risk:** Potential for forged JWT tokens to be accepted if the JWT_SECRET is ever compromised. No issuer binding means tokens could be accepted from unintended services.  
**Fix:** Add explicit issuer validation in Parse():
```go
if claims.Issuer != "remain-faithful" {
  return nil, errors.New("invalid issuer")
}
```

---

### Classify Endpoint Allows Text Logging to Cloud (Privacy Risk)
**File:** `/Users/jbhome/remain-faithful/RemainFaithful/RemainFaithfulBroadcast/SampleHandler.swift:182-198`  
**Issue:** The Tier 3 `/classify` endpoint sends extracted OCR text (up to 500 characters) to the backend's `/classify` handler, which forwards it to Anthropic Claude. The backend `/classify` handler logs the request with `writeError()` and error messages, but no explicit logging of the text content itself. However:
1. The text IS sent to an external service (Anthropic)
2. The Anthropic API will log requests per their data retention policy
3. Users may not realize raw text extracted from screen (even if truncated to 500 chars) is being sent to a third-party cloud service

**Risk:** Privacy disclosure gap. Broadcast extension sends OCR-extracted text to Anthropic Claude API without a clear data processing agreement statement visible to users. This should be disclosed prominently in privacy policy and app privacy manifest.  
**Fix:**
1. Update privacy manifest to declare Anthropic API calls
2. Update Privacy Policy to explicitly state "OCR text is sent to Anthropic Claude for classification"
3. Consider adding a user-facing notice before first Tier 3 classification
4. Add explicit request logging in `/classify` handler with text prefix for audit trail

---

## High

### Swallowed Error on Email Sending (Partner Invite)
**File:** `/Users/jbhome/remain-faithful/backend/internal/handler/invites.go:70`  
**Issue:** When inviting a partner whose email is already registered, the code sends an invite email but **silently ignores errors**:
```go
_ = h.Email.SendPartnerInvite(req.Email, inviterName, acceptURL)
```
If the email service is down or fails, the user receives a 201 Created response but no notification email reaches the invitee. The inviter has no way to know the email failed.

**Risk:** Silent failure leading to user confusion. Invitees never receive invite notifications.  
**Fix:** Return error or at least log with `log.Printf()` like the other invite path (line 95).

---

### Swallowed Error on Group Email Invite
**File:** `/Users/jbhome/remain-faithful/backend/internal/handler/invites.go:236`  
**Issue:** Same pattern — email errors are silently swallowed:
```go
_ = h.Email.SendGroupInvite(req.Email, inviterName, groupName, acceptURL)
```

**Risk:** Group invitees never receive notifications due to email failures with no feedback to the user.  
**Fix:** Log or return error.

---

### No Rate Limiting on Authentication Endpoints
**File:** `/Users/jbhome/remain-faithful/backend/cmd/server/main.go`  
**Issue:** No rate limiting middleware is configured for:
- `/auth/register` — unauthenticated, allows account creation spam
- `/auth/login` — unauthenticated, allows brute-force attempts
- `/auth/forgot-password` — unauthenticated, allows email enumeration and spam
- `/auth/reset-password` — unauthenticated, allows token guessing

**Risk:** Brute-force attacks, email list enumeration, spam/DoS of password reset email system.  
**Fix:** Add middleware like `github.com/noelyoo/ratelimit` or similar:
```go
r.Use(ratelimit.Middleware(100 * time.Minute)) // per IP
```

---

### No Rate Limiting on /classify Endpoint
**File:** `/Users/jbhome/remain-faithful/backend/cmd/server/main.go:166`  
**Issue:** The `/classify` endpoint is only protected by an optional shared secret (CLASSIFY_SECRET env var). If the secret is:
1. Empty (default) — endpoint is completely open
2. Weak or public — endpoint is open to DoS/abuse

No rate limiting is applied per IP or API key, allowing unlimited calls to the expensive Anthropic Claude API.

**Risk:** Denial of service via unlimited classification requests. Massive AWS bills if endpoint is abused.  
**Fix:** Apply rate limiting + require a strong, unique CLASSIFY_SECRET:
```go
ratelimit.Middleware(10 * time.Minute) // strict limit per IP
```

---

### Missing Error Check on `rows.Close()`
**File:** `/Users/jbhome/remain-faithful/backend/internal/handler/heartbeat.go:102`  
**Issue:** `rows.Close()` is called twice — once deferred at line 75 and manually at line 102. The second call returns no error checked, which is minor but redundant.

**Risk:** Low. Minor code quality issue.  
**Fix:** Remove the redundant `rows.Close()` at line 102; defer is sufficient.

---

### Unchecked Error When Rendering CGImage in Broadcast Extension
**File:** `/Users/jbhome/remain-faithful/RemainFaithfulBroadcast/SampleHandler.swift:358-360`  
**Issue:** `ciContext.createCGImage()` is called but if it fails (returns nil), the frame analysis bails out silently. No logging of what went wrong. Could mask memory pressure or CI context issues.

**Risk:** Silent failures during analysis; memory issues undetected.  
**Fix:** Add logging:
```swift
guard let cgImage: CGImage = autoreleasepool(invoking: {
  ciContext.createCGImage(renderImage, from: renderImage.extent)
}) else { 
  logger.warning("Failed to create CGImage from CIImage")
  return 
}
```

---

### Missing Input Validation: Negative Stream Bitrates
**File:** `/Users/jbhome/remain-faithful/backend/internal/handler/donations.go:27`  
**Issue:** `body.AmountDollars < 1` check exists, but nothing prevents negative values:
```go
if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.AmountDollars < 1
```
If the JSON contains `"amount_dollars": -50`, it passes the `< 1` check and a refund/credit could be applied.

**Risk:** Users could create negative donations (refunds/credits) without authorization.  
**Fix:** Add explicit non-negative check:
```go
if body.AmountDollars < 1 || body.AmountDollars > 10000 {
  writeError(w, http.StatusBadRequest, "invalid amount")
  return
}
```

---

### Try? Silently Swallows Errors on Token Refresh
**File:** `/Users/jbhome/remain-faithful/RemainFaithful/APIClient.swift:364,373,398,384`  
**Issue:** Multiple places use `try?` on `refreshTokenIfNeeded()`:
```swift
try? await refreshTokenIfNeeded()
```
If token refresh fails (network, server error, expired refresh), it's silently ignored. The stale/expired token continues to be used, causing silent auth failures downstream.

**Risk:** Users encounter unexplained "unauthorized" errors on API calls after token expiry.  
**Fix:** Log failures:
```swift
do {
  try await refreshTokenIfNeeded()
} catch {
  print("[API] Token refresh failed: \(error)")
}
```

---

## Medium

### No X-Frame-Options Header (Click-jacking Risk on Website)
**File:** `/Users/jbhome/remain-faithful/website/src/app/layout.tsx`  
**Issue:** Website does not set `X-Frame-Options: DENY` header. If the website is served from a Vercel deployment, Vercel may not set it by default.

**Risk:** Clickjacking attacks on the website; CORS misconfiguration.  
**Fix:** Add headers in Next.js (next.config.js or via middleware):
```js
headers: [
  { key: 'X-Frame-Options', value: 'DENY' },
  { key: 'X-Content-Type-Options', value: 'nosniff' },
  { key: 'X-XSS-Protection', value: '1; mode=block' },
]
```

---

### API Key Hardcoded in Device Analytics
**File:** `/Users/jbhome/remain-faithful/RemainFaithful/RemainFaithfulApp.swift:128-131`  
**Issue:** Google Sign-In client ID is loaded from `GoogleService-Info.plist` bundle resource. While not a hardcoded string, it IS embedded in the app binary and could be reverse-engineered. No issue per se, but worth noting for completeness.

**Risk:** Low. Google client IDs are not sensitive secrets; they're meant to be public.  
**Fix:** No fix needed; this is correct practice.

---

### Donation Notification Background Fetch Not Awaited
**File:** `/Users/jbhome/remain-faithful/backend/internal/handler/donations.go:120-132`  
**Issue:** APNs push notifications are sent in a loop without awaiting or checking errors thoroughly:
```go
for rows.Next() {
  var token string
  if rows.Scan(&token) != nil {
    continue
  }
  _ = h.APNS.Send(r.Context(), &apns.Notification{...})
}
```
If a send fails, it's silently discarded. If the device token is invalid, the token should be marked inactive (as is done in other notification paths), but here it is not.

**Risk:** Stale device tokens accumulate and continue to fail silently. No cleanup.  
**Fix:** Apply the same invalidation pattern:
```go
if err := h.APNS.Send(ctx, n); err != nil {
  var invalidErr *apns.ErrInvalidToken
  if errors.As(err, &invalidErr) {
    h.markTokenInactive(ctx, token)
  }
}
```

---

### Missing Content-Security-Policy Header
**File:** `/Users/jbhome/remain-faithful/website/src/app/layout.tsx`  
**Issue:** No CSP header is set to restrict script/style sources. External CDNs or injected scripts could execute.

**Risk:** XSS vulnerabilities if input sanitization is bypassed.  
**Fix:** Add CSP header:
```
Content-Security-Policy: default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'
```

---

### Overly Permissive CORS Configuration
**File:** `/Users/jbhome/remain-faithful/backend/cmd/server/main.go:126-149`  
**Issue:** CORS middleware allows `GET, POST, PUT, DELETE, OPTIONS` from **any origin in the allowed list**:
```go
w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
```
If the ALLOWED_ORIGINS contains `*` or a malicious domain, all methods are allowed. No credential control.

**Risk:** Moderate. If origins list is misconfigured, CORS could expose sensitive endpoints.  
**Fix:** Separate method allowlists per endpoint type:
```go
if strings.HasPrefix(r.URL.Path, "/admin") {
  w.Header().Set("Access-Control-Allow-Methods", "GET")
} else {
  w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
}
```

---

### No SQL Injection Validation on URL Slugs
**File:** `/Users/jbhome/remain-faithful/backend/internal/handler/groups.go` (not fully reviewed, but general observation)  
**Issue:** While parameterized queries are used throughout, no explicit length/character validation on email or name fields before inserting into database. Email regex is not validated with RFC-compliant parser.

**Risk:** Low. Go's database/sql driver protects against SQL injection via parameterization. However, no length limits on names/emails could cause truncation issues.  
**Fix:** Add length checks:
```go
if len(req.Email) > 254 {
  writeError(w, http.StatusBadRequest, "email too long")
  return
}
```

---

### Device Token Upsert Logs Entire Token
**File:** `/Users/jbhome/remain-faithful/backend/internal/handler/notifications.go:56`  
**Issue:** Logs the full token (first 8 chars):
```go
log.Printf("[push] registered token %.8s... user=%d env=%s", req.Token, userID, req.Environment)
```
While prefixed, it still includes a partial token in logs. If logs are shipped externally, tokens could be reconstructed.

**Risk:** Low to medium. Depends on log storage and access controls.  
**Fix:** Log only token count, not partial tokens:
```go
log.Printf("[push] registered 1 token for user=%d env=%s", userID, req.Environment)
```

---

## Low

### Unused Import in APIClient
**File:** `/Users/jbhome/remain-faithful/RemainFaithful/APIClient.swift`  
**Issue:** No unused imports detected; code is clean.

**Risk:** None. Code quality is good.  
**Fix:** None needed.

---

### Hardcoded Default Contact Email
**File:** `/Users/jbhome/remain-faithful/backend/internal/handler/contact.go:39`  
**Issue:** Default contact email is hardcoded to `jeff@hanokventures.co`:
```go
toEmail := "jeff@hanokventures.co"
```
If the email address changes, the code must be redeployed.

**Risk:** Low. Contact form becomes unresponsive if env var is not set.  
**Fix:** Require `CONTACT_TO_EMAIL` env var with no fallback, or document the default clearly.

---

### Potential Memory Leak in Broadcast Extension
**File:** `/Users/jbhome/remain-faithful/RemainFaithfulBroadcast/SampleHandler.swift:286,321-324`  
**Issue:** Task objects are retained (`heartbeatTask`, `retryTask`). While they're cancelled in `broadcastFinished()`, if the extension crashes, these tasks could leak memory.

**Risk:** Very low. Extensions run in isolated processes with memory limits.  
**Fix:** No fix needed; memory is reclaimed when the extension process exits.

---

### Inconsistent Error Messaging in Contact Form
**File:** `/Users/jbhome/remain-faithful/website/src/components/ContactForm.tsx:36`  
**Issue:** Error messages sometimes reference `jeff@hanokventures.co` but backend env var may differ. If the backend uses a different email, users see stale contact info.

**Risk:** Very low. User confusion about who to contact.  
**Fix:** Backend should return contact email in error response, or website should fetch it from config.

---

### No Timeout on Donate Checkout Session
**File:** `/Users/jbhome/remain-faithful/RemainFaithful/APIClient.swift:249-256`  
**Issue:** The Stripe checkout session URL fetch has a hardcoded 15-second timeout in `makeRequest()`. If the backend is slow, users see "timeout" instead of a proper error.

**Risk:** Very low. User sees a generic timeout message.  
**Fix:** Increase timeout to 30s or return a more helpful error message.

---

### Empty Error Messages in Error Responses
**File:** `/Users/jbhome/remain-faithful/backend/cmd/server/main.go:117`  
**Issue:** Some error responses return just `{"error":""}` if error.Error() is empty or missing.

**Risk:** Very low. Poor UX but not a security risk.  
**Fix:** Always include a human-readable error message.

---

## Summary of Findings

**Critical (3):**
1. Missing PrivacyInfo.xcprivacy — **App Store rejection risk**
2. JWT claims not fully validated — **token forgery risk**
3. OCR text sent to Anthropic without disclosure — **privacy risk**

**High (5):**
1. Email errors swallowed (2 places)
2. No rate limiting on auth endpoints
3. No rate limiting on /classify endpoint
4. Duplicate rows.Close() call
5. Unchecked CGImage rendering errors

**Medium (4):**
1. Missing X-Frame-Options header
2. Invalid device tokens not cleaned up on donation notifications
3. Missing CSP header
4. Overly permissive CORS methods

**Low (4):**
1. Hardcoded contact email fallback
2. Potential memory leak in broadcast extension (very unlikely)
3. Error messaging inconsistency
4. Timeout handling in checkout

**Priority Actions:**
1. **Immediate:** Create PrivacyInfo.xcprivacy and update privacy policy
2. **This sprint:** Add rate limiting to auth and classify endpoints
3. **This sprint:** Fix email error handling in invites
4. **This sprint:** Add issuer validation to JWT parsing
5. **Next sprint:** Add CSP and X-Frame-Options headers
