// Screen-by-screen guide to the Remain Faithful iPhone app.
// Every entry here was checked against the SwiftUI source in RemainFaithful/*.swift.
// Only describe controls that do what this text says. If the app changes, update this file.

export type GuideButton = {
  name: string
  does: string
  why?: string
}

export type GuideShot = {
  src: string
  alt: string
  caption: string
}

export type GuidePlate = {
  kicker: string
  title: string
  lines: { label: string; detail: string }[]
  caption: string
}

export type GuideScreen = {
  id: string
  title: string
  where: string
  intro: string
  shots?: GuideShot[]
  plate?: GuidePlate
  buttons: GuideButton[]
  note?: string
}

export type GuidePart = {
  id: string
  label: string
  title: string
  summary: string
  screens: GuideScreen[]
}

export const guideParts: GuidePart[] = [
  {
    id: 'getting-started',
    label: 'Getting started',
    title: 'Getting started',
    summary:
      'The first time you open the app, it walks you through a few short screens: your account, a partner, Screen Time, Deep Scan, and the apps you choose.',
    screens: [
      {
        id: 'welcome',
        title: 'Welcome',
        where: 'First screen after you install the app',
        intro: 'This is the first thing you see. It has two choices.',
        shots: [
          {
            src: '/screens/welcome.png',
            alt: 'Remain Faithful welcome screen with a gold shield, the tagline "Accountability between trusted friends," a Begin Your Journey button, and a Sign In link.',
            caption: 'Welcome screen',
          },
        ],
        buttons: [
          {
            name: 'Begin Your Journey',
            does: 'Starts setting up a new account.',
            why: 'Pick this if you have never used Remain Faithful before.',
          },
          {
            name: 'Sign In',
            does: 'Opens the sign-in screen.',
            why: 'Use this if you already have an account, like after you get a new phone.',
          },
        ],
      },
      {
        id: 'create-account',
        title: 'Create your account',
        where: 'After you tap Begin Your Journey',
        intro: 'Three boxes and one button.',
        shots: [
          {
            src: '/screens/create.png',
            alt: 'Create your account screen with fields for name, email address, and password, and a Get Started button.',
            caption: 'Create your account',
          },
        ],
        buttons: [
          {
            name: 'Your name',
            does: 'The name your partners and group see.',
            why: 'Use the name they know you by. Accountability only works when your partners know who they are walking with.',
          },
          {
            name: 'Email address',
            does: 'The email you sign in with.',
            why: 'If a partner or group already invited you, use the same email they sent the invite to. The app connects you to them as soon as your account is made.',
          },
          {
            name: 'Password',
            does: 'Must be at least 8 characters.',
          },
          {
            name: 'Get Started',
            does: 'Creates your account and moves to the next step. Tapping it means you agree to the Terms of Service and Privacy Policy.',
          },
        ],
      },
      {
        id: 'sign-in',
        title: 'Sign in',
        where: 'Tap Sign In on the welcome screen',
        intro: 'For people who already have an account.',
        buttons: [
          { name: 'Sign In', does: 'Signs you in with your email and password.' },
          {
            name: 'Forgot Password?',
            does: 'Sends a reset link to your email. Open the link, choose a new password, and the link works for one hour.',
          },
          {
            name: 'Sign in with Apple',
            does: 'Signs you in with your Apple ID instead of a password.',
          },
          { name: 'Create Account', does: 'Takes you back to new account setup.' },
        ],
      },
      {
        id: 'invite-partner',
        title: 'Invite your accountability partner',
        where: 'Second setup screen',
        intro: 'Here you invite the one person you want to walk with you.',
        buttons: [
          {
            name: "Partner's email address",
            does: 'Type the email of a trusted friend.',
          },
          {
            name: 'Send Invitation',
            does: 'Emails your partner an invite. If they do not have the app yet, they are connected to you once they create an account with that email.',
            why: 'It is easier to stay on track when someone you trust is paying attention.',
          },
          {
            name: 'Skip for now',
            does: 'Moves on without inviting anyone.',
            why: 'You can add a partner later in Settings, then Manage Partners.',
          },
        ],
      },
      {
        id: 'turn-on-monitoring',
        title: 'Allow Screen Time',
        where: 'After you invite a partner, or skip that step',
        intro:
          'This screen is for Alerts and Blocking. Deep Scan is the next screen, explained once.',
        buttons: [
          {
            name: 'Allow Screen Time',
            does: "Opens Apple's Screen Time permission box. You approve it on your own phone. If it is approved, this screen moves on.",
            why: 'The app needs this permission so Alerts and Blocking can use the apps you choose. If you skip it, monitoring stays off until you allow it in Settings.',
          },
          {
            name: 'Skip for now',
            does: 'Moves on without Screen Time permission. Monitoring stays off.',
          },
          {
            name: 'Continue',
            does: 'Shows when Screen Time is already allowed. Opens the Deep Scan screen.',
          },
        ],
      },
      {
        id: 'deep-scan-setup',
        title: 'Deep Scan',
        where: 'After Screen Time, before you choose apps',
        intro:
          'One short screen explains optional Deep Scan. Home does not keep a Deep Scan button after this.',
        buttons: [
          {
            name: 'Continue',
            does: 'Moves on to choosing apps. This button does not start Deep Scan. You start it later from Control Center.',
          },
        ],
        note: 'The screen says what Deep Scan does, how to open Control Center, and how to add Screen Recording in Settings if it is missing. It also says protected video such as Netflix is not visible, and that stopping or pausing Deep Scan notifies your partners.',
      },
    ],
  },
  {
    id: 'home',
    label: 'Home tab',
    title: 'Home tab',
    summary:
      'Home is your daily check-in. It shows whether Deep Scan is on, your clean streak, and any recent flags.',
    screens: [
      {
        id: 'home-screen',
        title: 'Home',
        where: 'First tab at the bottom of the screen',
        intro: 'Here is everything you can see and tap on Home.',
        plate: {
          kicker: 'Home',
          title: 'Daily check-in',
          caption: 'What Home shows',
          lines: [
            { label: 'Status', detail: 'Whether Alerts and Blocking are on' },
            { label: 'Clean streak', detail: 'Days in a row without a flag' },
            { label: 'Recent flags', detail: 'Your latest alerts' },
          ],
        },
        buttons: [
          {
            name: 'Monitoring status (card)',
            does: 'Shows whether Alerts and Blocking are on. Deep Scan is a smaller line. There is no Deep Scan button on this card. You start and stop Deep Scan from Control Center. Pausing or turning off monitoring notifies your partners. If setup is unfinished, the card has one button: Allow Screen Time, Choose apps, or Turn monitoring on.',
          },
          {
            name: 'Verse of the Day',
            does: 'A short Bible verse. A new one shows each day.',
          },
          {
            name: 'Clean streak',
            does: 'The big number is how many days in a row you have gone without a flag. "Best" is your longest streak. The row of circles is this week. A check mark is a clean day. An X is a day with a flag.',
            why: 'Seeing the days add up helps you keep going.',
          },
          {
            name: 'Recent Flags',
            does: 'Your latest alerts. Tap one to open it. If you have none, it says "No flags yet."',
          },
          {
            name: 'I Need Support Right Now',
            does: 'Opens a screen for a hard moment, with a verse, your partner\'s name, and a button to reach them.',
            why: 'When temptation hits, the hardest part is reaching out. This makes it one tap.',
          },
          {
            name: 'Monitored app used today (banner)',
            does: 'Can show up after you use an app you chose to monitor. It says Alerts and Blocking stay as you set them, and that Deep Scan is optional. The X hides it. It does not open Deep Scan steps.',
          },
          {
            name: 'Number on the Home icon',
            does: 'How many new alerts you have. It clears when you open Home.',
          },
          {
            name: 'Partner alert at the top',
            does: 'If a partner gets flagged while you have the app open, a banner slides in with their name. Tap it to read the alert. It goes away by itself after six seconds.',
          },
        ],
      },
      {
        id: 'deep-scan',
        title: 'Starting Deep Scan',
        where: 'Explained once during setup. Home has no Deep Scan button.',
        intro:
          'Deep Scan is optional. You turn it on yourself from Control Center when you want a closer look. Setup walks through this once.',
        buttons: [
          {
            name: 'How you start it',
            does: 'Swipe down from the top-right corner to open Control Center. If Screen Recording is missing, open Settings, then Control Center, and add Screen Recording. Touch and hold Screen Recording, choose Remain Faithful, then tap Start Broadcast.',
          },
        ],
        note: 'While Deep Scan runs, your phone checks what is on screen, on the phone itself. iOS shows its red recording indicator so you always know it is on. Streaming apps like Netflix, Disney+, and Hulu show up black to Deep Scan, so it cannot check them. Stopping or pausing Deep Scan notifies your partners. Screenshots and the text on your screen never leave your phone.',
      },
      {
        id: 'alert-detail',
        title: 'An alert',
        where: 'Tap any item in Recent Flags',
        intro: 'This is what you find when you open an alert.',
        plate: {
          kicker: 'Alert',
          title: 'Inside an alert',
          caption: 'What an alert includes',
          lines: [
            { label: 'Category and level', detail: 'The kind of alert, and how serious' },
            { label: 'When', detail: 'The time it happened' },
            { label: 'Summary', detail: 'A short line the app wrote' },
            { label: 'Conversation Starter', detail: 'Questions for the talk' },
            { label: 'Mark as Discussed', detail: 'Saved after you talk it through' },
          ],
        },
        buttons: [
          {
            name: 'Category and level',
            does: 'What kind of alert it was and how serious.',
          },
          {
            name: 'Detected',
            does: 'When it happened and a short summary the app wrote.',
          },
          {
            name: 'Conversation Starter',
            does: 'A few questions to use when you talk it through with your partner.',
            why: 'Starting the conversation is often the hardest part. These give you a first line.',
          },
          {
            name: 'Mark as Discussed',
            does: 'Tap it after you have talked it through. It is saved to your account and the button turns green.',
          },
        ],
      },
      {
        id: 'support-now',
        title: 'I Need Support Right Now',
        where: 'Purple button on Home',
        intro: 'This screen is for a moment of temptation.',
        buttons: [
          {
            name: 'Hold On To This',
            does: 'Shows 1 Corinthians 10:13.',
          },
          {
            name: 'Your Accountability Partner',
            does: 'Shows your primary partner\'s name. The star in Manage Partners marks who that is. The alert itself goes to every accepted partner and to your group, not only the starred person. If you have no partner and no group yet, it asks you to add one.',
          },
          {
            name: 'Alert Partners and Group',
            does: 'Sends an urgent notification to every accepted partner and to the other people in your groups. The button turns green and says Alert Sent only after the notice goes out. If it cannot send, the screen says so.',
            why: 'The people you chose hear from you right away, before things go further.',
          },
          { name: "I'm okay, close this", does: 'Closes the screen. Nothing is sent.' },
        ],
      },
    ],
  },
  {
    id: 'group',
    label: 'Group tab',
    title: 'Group tab',
    summary:
      'A group is a small circle of people, up to 12, who keep each other accountable. The Group tab shows how everyone is doing.',
    screens: [
      {
        id: 'group-screen',
        title: 'Your group',
        where: 'Middle tab at the bottom of the screen',
        intro: 'What you see once you are in a group.',
        shots: [
          {
            src: '/screens/group.png',
            alt: 'Group tab for Tuesday Men\'s Group with five members, each showing a status label and day streak, a Group Covenant row, and an Invite Member button.',
            caption: 'Group tab',
          },
        ],
        buttons: [
          {
            name: 'Group name',
            does: 'Your group\'s name and how many people are in it. The pencil saves a new name for everyone in the group.',
          },
          {
            name: 'Members',
            does: 'Each person\'s clean streak and a colored dot. Green "Strong" means no flags in the last 30 days. Orange "Watchful" means one or two. Red "Struggling" means three or more.',
            why: 'You can see at a glance who might need a call this week.',
          },
          {
            name: 'Tap a member',
            does: 'Opens their streak and their recent flags. Send Encouragement sends that person a short notice. It says Encouragement Sent only after the notice goes out. Close takes you back.',
          },
          {
            name: 'Group Covenant',
            does: 'Shows the wording saved for your group. Edit Covenant saves it and notifies the other members. They are not asked to sign it again. The notice does not include the wording.',
          },
          {
            name: 'Invite Member',
            does: 'Opens a sheet to add someone. Type an email and the app sends them an invite, even if they do not have the app yet. They join when they sign up with that email. Or tap Share Invite Link to send it by text or any other app.',
          },
        ],
      },
      {
        id: 'create-group',
        title: 'Create a Group',
        where: 'Group tab, when you are not in a group yet',
        intro:
          'If you are not in a group, the tab says "No Group Yet" and shows a Create a Group button. It takes four short steps.',
        buttons: [
          { name: 'Name your group', does: 'Pick a name your members will recognize.' },
          {
            name: 'Group Covenant',
            does: 'The agreement members sign. You can use ours or change the wording.',
          },
          {
            name: 'Invite Members',
            does: 'Type the emails of the people you want. Add another email adds a new box. Groups hold up to 12 people.',
          },
          { name: 'Review', does: 'Check everything, then create the group. Invites go out by email.' },
        ],
      },
    ],
  },
  {
    id: 'settings',
    label: 'Settings tab',
    title: 'Settings tab',
    summary:
      'Settings is where you manage partners, groups, the apps you want blocked, and your account.',
    screens: [
      {
        id: 'settings-account',
        title: 'Profile and accountability',
        where: 'Top of the Settings tab',
        intro: 'Your name and email are at the top, followed by the Accountability section.',
        buttons: [
          { name: 'Edit', does: 'Change your name or email.' },
          {
            name: 'Manage Partners',
            does: 'Lists your one-to-one partners. Tap the star to make someone your primary partner. The star changes only after it is saved. The I Need Support button reaches every accepted partner and your group, not only the starred person. Remove ends a partnership and lets them know. Add New Partner sends an invite by email.',
          },
          {
            name: 'Manage Groups',
            does: 'Lists the groups you belong to. Each one has Invite Member and Leave Group. Create New Group starts a new one.',
          },
          { name: 'View Covenant', does: 'Shows the wording saved for your group. If the group has not saved its own, you see a sample labeled as a sample.' },
        ],
      },
      {
        id: 'settings-restrictions',
        title: 'App Restrictions',
        where: 'Settings, then App Restrictions',
        intro:
          'This is where you choose which apps the app watches or blocks. It uses Apple\'s Screen Time, so it keeps working in the background.',
        buttons: [
          {
            name: 'Screen Time Permission',
            does: 'Shows if you have given permission. If not, tap to turn it on.',
          },
          {
            name: 'Choose Apps to Restrict',
            does: "Opens Apple's own list of apps and categories. Pick the ones that tempt you.",
          },
          {
            name: 'Monitor App Usage',
            does: 'If Alerts is on, your partner gets one alert that day after a chosen app or category is used for a minute.',
          },
          {
            name: 'Block Selected Apps',
            does: 'If Blocking is on and you open an app you asked to block, your partner gets an alert that you opened an app you asked to block, and Apple\'s Screen Time shield appears. The app name is never sent.',
            why: 'A block that is already in place does more than willpower in the moment.',
          },
        ],
      },
      {
        id: 'settings-protection',
        title: 'Monitoring and protection',
        where: 'Settings, Monitoring and Protection sections',
        intro: 'These keep your setup from being switched off quietly.',
        buttons: [
          {
            name: 'Notifications',
            does: 'Opens iPhone Settings for Remain Faithful, where you turn alerts on or off.',
          },
          {
            name: 'App Lockout (Deep Scan)',
            does: 'Optional. When on, the apps you chose stay blocked unless Deep Scan is running. If you stop Deep Scan, the block comes back. Pausing or turning off monitoring notifies your partners, whether or not this is on. Turning App Lockout off also sends a notice.',
          },
          {
            name: 'Partner PIN',
            does: 'Shows whether a partner has set a PIN for you. Your partner sets it from their own phone in Manage Partners. Once it is set, opening App Restrictions or turning off App Lockout asks for that PIN. A wrong PIN sends your partner a notice.',
            why: 'It keeps a weak moment from undoing a decision you made on a strong day.',
          },
        ],
      },
      {
        id: 'settings-support',
        title: 'Activity log, support, and account',
        where: 'Bottom half of Settings',
        intro: 'The rest of the Settings tab.',
        buttons: [
          { name: 'View My Activity Log', does: 'Your full history of flagged activity.' },
          { name: 'How It Works', does: 'A short three-step summary of how the app works, with a link to this guide.' },
          { name: 'Replay App Tour', does: 'Walks through the main buttons again, starting on Home. You can skip any step.' },
          { name: 'Contact Support', does: 'Opens an email to support@remainfaithful.com.' },
          { name: 'Rate the App', does: 'Asks iOS to show the App Store rating box.' },
          {
            name: 'Send Ideas or Report a Problem',
            does: 'Opens a form for an idea, an improvement, or a problem. Your account name and email are included so we can reply. App version and iPhone model are included only if you turn that on. If sending fails, you can email support@remainfaithful.com instead.',
          },
          {
            name: 'Leave All Groups',
            does: 'Takes you out of every group and lets the other members know. You confirm first. Your alert history stays until you delete your account.',
          },
          {
            name: 'Delete Account',
            does: 'Deletes your account and your data for good and lets your partners know you left. It cannot be undone.',
          },
          { name: 'Sign Out', does: 'Signs you out of this phone.' },
        ],
      },
    ],
  },
]
