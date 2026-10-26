# Remain Faithful for Android

Native app (`com.remainfaithful.app`) for the same API as the iOS app: `https://remain-faithful-api.fly.dev`.

Colors live in `app/src/main/kotlin/com/remainfaithful/app/theme/RfTheme.kt` and `app/src/main/res/values/colors.xml`. Change both when the gold moves to a brighter yellow.

## Build

```bash
cd android
./gradlew :core:test :app:assembleDebug
```

The debug APK is `app/build/outputs/apk/debug/app-debug.apk`. GitHub Actions uploads that file as the artifact `remain-faithful-android-debug` on pull requests that touch `android/`.

`app/google-services.json` in this public repo is a placeholder. Replace it on a release machine with the file from Firebase. Do not commit a real service-account private key.

## Google sign-in

The login and create-account screens send a Google ID token to `POST /auth/google`, the same call the iOS app uses. The token audience is the Web client id.

Set `GOOGLE_WEB_CLIENT_ID` in `android/local.properties` or in the environment. Do not commit a real client id. On the server, add that same Web client id to `GOOGLE_CLIENT_IDS` (comma-separated, next to the iOS client id). Until the id is set, the button says Google Sign In is not available.

## Deep Scan

Screen capture is off for v1. `DEEP_SCAN_ENABLED` defaults to false. The service code stays in `monitoring/DeepScanService.kt` and `src/deepScan/AndroidManifest.xml`. Those permissions are merged only when you build with `DEEP_SCAN_ENABLED=true`. The v1 manifest does not request media projection, and the screens do not mention Deep Scan.

## Permissions

| Permission | Why |
| --- | --- |
| Internet | Sign in, groups, partners, alerts, and heartbeat use the existing API. |
| Notifications | Partner alerts, the urgent request, and the monitoring reminder. |
| Usage access | Time the apps the person chose. The package name stays on the phone. |
| Display over other apps | Draw the blocking shield the person turned on. |
| Foreground service (special use) | Keep Alerts and Blocking running, with a visible notification. |
| Run at startup | After reboot, remind the person to open the app so monitoring can continue. |
| Media projection | Not in the v1 manifest. Used only if Deep Scan is turned on for a later build. |

Usage access and the shield each have an in-app disclosure before the system screen.

## What leaves the phone

An alert is category, severity, summary, and time. A heartbeat is `active`, `idle`, or `monitoring`, plus an opaque device id and the device type `android`. The id is created on the phone and stored in app preferences. It is not the push token, an advertising id, or a package name. Package names, web addresses, and screen contents stay on the phone.

If this phone was removed from the account, the server answers `409` with code `device_removed`. The app stops Alerts and Blocking on this phone and shows: "This phone was removed from the account. Monitoring on this phone has stopped. Turn Alerts or Blocking on again to add it back." Turning monitoring on again calls reclaim, then heartbeats resume. Settings lists My devices for this account.

## Publishing

1. Sign up for Google Play Console. The registration fee is $25 once.
2. In Firebase, add an Android app with package name `com.remainfaithful.app`. Download `google-services.json` and use it for the release build. Keep it out of this public repo.
3. Create a Firebase service account that can send Cloud Messaging. On Fly, set `FCM_PROJECT_ID`, `FCM_CLIENT_EMAIL`, and `FCM_PRIVATE_KEY`. Leave the Apple APNs key as it is.
4. Create an upload keystore and store it somewhere safe, not in git. Turn on Play App Signing. In GitHub Actions secrets, add `ANDROID_KEYSTORE_BASE64` (the keystore file, base64), `ANDROID_KEYSTORE_PASSWORD`, `ANDROID_KEY_ALIAS`, and `ANDROID_KEY_PASSWORD`. Also add `GOOGLE_WEB_CLIENT_ID` when the Firebase Web client exists. The release job builds `app-release.aab` only after those four signing secrets are set, and uploads it as `remain-faithful-android-release`.
5. Play listing: adults 18 and up, self-chosen accountability, privacy policy `https://remainfaithful.com/privacy`. Data safety answers for this v1 build:
   - Account info: name and email, from the form or from Google sign-in. Required for the account. Sent to the Remain Faithful server.
   - Device or other IDs: the Firebase push token, so partner alerts can arrive, and an opaque device id so this phone has its own heartbeat. That id is not an advertising id.
   - The app does not collect browsing history, photos, contacts, screen contents, or the names of other apps. Usage access stays on the phone and becomes an alert category, a fixed summary, and a heartbeat.
   - Data is not sold.
   - Do not declare screen capture. Deep Scan is not in this build.
6. Ship the signed Android App Bundle to the internal testing track first. The CI debug APK is for sideloading a test build, not the Play upload. Move to production after the listing, content rating, and Data safety form are done.
