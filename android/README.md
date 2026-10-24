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

## Permissions

| Permission | Why |
| --- | --- |
| Internet | Sign in, groups, partners, alerts, and heartbeat use the existing API. |
| Notifications | Partner alerts, the urgent request, and the monitoring reminder. |
| Usage access | Time the apps the person chose. The package name stays on the phone. |
| Display over other apps | Draw the blocking shield the person turned on. |
| Foreground service (special use) | Keep Alerts and Blocking running, with a visible notification. |
| Foreground service (media projection) | Deep Scan, only after the person starts it and approves screen capture. |
| Run at startup | After reboot, remind the person to open the app so monitoring can continue. |

Usage access and the shield each have an in-app disclosure before the system screen.

## What leaves the phone

An alert is category, severity, summary, and time. A heartbeat is `active`, `idle`, or `monitoring`. Package names, web addresses, and screen contents stay on the phone.

## Publishing

1. Sign up for Google Play Console. The registration fee is $25 once.
2. In Firebase, add an Android app with package name `com.remainfaithful.app`. Download `google-services.json` and use it for the release build. Keep it out of this public repo.
3. Create a Firebase service account that can send Cloud Messaging. On Fly, set `FCM_PROJECT_ID`, `FCM_CLIENT_EMAIL`, and `FCM_PRIVATE_KEY`. Leave the Apple APNs key as it is.
4. Create an upload keystore and store it somewhere safe, not in git. Turn on Play App Signing.
5. Play listing: adults 18 and up, self-chosen accountability, privacy policy `https://remainfaithful.com/privacy`. On the Data safety form, say the app collects account info (name, email) and a push token. It does not collect browsing history, screen contents, or the names of other apps. Data is not sold. It is sent to the Remain Faithful server so a partner can be notified.
6. Ship a signed Android App Bundle to the internal testing track first. The CI debug APK is for sideloading a test build, not the Play upload. Move to production after the listing, content rating, and Data safety form are done.
