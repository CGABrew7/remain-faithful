# Remain Faithful

Remain Faithful is a privacy-first accountability app for adults 18 and older. You choose a partner or a group of up to 12 people. It is self monitoring. It is not a tool for watching someone else's phone.

Only event metadata and a heartbeat leave the device. Screenshots, screen text, and raw screen content stay on the device.

## Privacy model

A partner, or someone in a group you join, can receive:

- the name and email on the account
- the time, category, severity, and a short summary of an alert
- a heartbeat while protection is on

If Blocking is on and a blocked app is opened, the partner gets an alert that a blocked app was opened. That alert does not include the app name.

Sign in can use email, Apple, or Google. The API runs on Fly. The website runs on Vercel. Account email is sent with Resend. Website analytics stay off until a visitor allows them.

## How to verify the privacy claims

Read these files. They are the behavior, not a separate promise.

- `RemainFaithfulShieldConfig/ShieldConfigurationExtension.swift` sends a blocked app event with category, severity, summary, and time. The summary is "Blocked app opened".
- `RemainFaithfulBroadcast/SampleHandler.swift` keeps Deep Scan classification on the phone.
- `RemainFaithfulDeviceActivity/DeviceActivityMonitorExtension.swift` sends a usage threshold event without an app name.
- `backend/internal/handler/events.go` stores the event metadata the apps upload.
- `desktop/core/privacy.go` lists what the desktop app is allowed to send.
- `website/src/app/privacy/page.tsx` is the public policy. It should match the files above.

## Repo layout

- `RemainFaithful/` iPhone app
- `RemainFaithfulShieldConfig/` Screen Time shield
- `RemainFaithfulBroadcast/` optional Deep Scan
- `RemainFaithfulDeviceActivity/` usage alerts
- `RemainFaithfulTests/` unit tests
- `android/` Android app
- `desktop/` desktop app
- `backend/` Go API
- `website/` public site
- `docs/` product notes

## Checks

- Backend: `cd backend && go test ./...`
- Website: `cd website && npm test`
- iPhone: the `iOS CI` GitHub Action builds the app and runs `RemainFaithfulTests` on a simulator. It does not sign for distribution and it does not upload a build.
