# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| Latest (`main`) | ✅ |

## Reporting a Vulnerability

**Please do not open a public GitHub issue for security vulnerabilities.**

Email **jeff@cornerstonelicensing.com** with the subject line `[SECURITY] Remain Faithful — <brief description>`.

Include:
- A description of the vulnerability and its potential impact
- Steps to reproduce (proof-of-concept code or screenshots welcome)
- Any suggested mitigations you have in mind

### What to expect

| Milestone | Target |
|-----------|--------|
| Initial acknowledgement | Within 48 hours |
| Triage and severity assessment | Within 5 business days |
| Status update / fix timeline | Within 10 business days |
| Public disclosure (coordinated) | After fix is deployed |

We follow coordinated disclosure: we ask that you give us a reasonable window to patch before publishing details publicly.

## Scope

**In scope:**
- Authentication and authorization bypass
- Data exposure (user events, alerts, account details)
- Server-side injection (SQL, command, etc.)
- Push notification spoofing or hijacking
- Stripe payment flow manipulation

**Out of scope:**
- Denial-of-service attacks
- Social engineering of staff or users
- Vulnerabilities in third-party dependencies already tracked by their own advisories
- Issues requiring physical access to a device

## Bug Bounty

This is an indie app with no formal bounty program. We can't offer financial rewards, but we will credit you (by name or handle, your choice) in the release notes for any confirmed vulnerability.

## Preferred Languages

Reports in English are preferred.
