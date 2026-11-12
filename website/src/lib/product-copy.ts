/**
 * Paste-ready product sentences.
 * ALERT_TRIGGER must match the sentence in docs/alert-trigger-copy.md.
 */

export const ALERT_TRIGGER =
  'If Alerts is on, your partner is told when a chosen app or category is used for 1 minute, then again at 5, 15, and 30 minutes that day. If Blocking is on, opening a chosen app shows Apple\'s Screen Time shield, Apple\'s adult web filter is on, and your partner gets an alert that a blocked app was opened. That alert does not include the app name.'

export const PRIVACY_PROMISE =
  'Only event metadata leaves the phone: category, severity, timestamp, and a short system summary, plus a heartbeat while protection is on. Screenshots and screen content stay on the device.'

export const PRIVACY_POINTER =
  'Screenshots and screen content stay on the device. What a partner can receive is on the privacy page.'

export const FAQ_PARTNERS_SEE =
  'No. Screenshots and screen content stay on your phone. The privacy section on this page says what a partner can receive.'

export const FAQ_SLIP =
  `${ALERT_TRIGGER} The goal is a conversation. Alerts include prompts that help your partner respond with grace.`

export const FAQ_DIFFERENT =
  'It is free, it is for adults who choose their own partners, and screenshots stay on the phone. The privacy section on this page says what a partner can receive. The iOS app, backend, and website source are public on GitHub.'

export const FAQ_STREAMING =
  'Deep Scan is optional and runs on your phone. DRM video apps such as Netflix do not share frames, so there is nothing from those apps to classify. Screenshots and screen content stay on your phone. Only an event category, severity, time, a short summary, and a heartbeat leave the device.'
