package com.remainfaithful.core

object Copy {
    const val appName = "Remain Faithful"
    const val welcomeSubtitle = "Accountability between\ntrusted friends"
    const val scripture = "But you must remain faithful to the things you have been taught."
    const val scriptureRef = "2 Timothy 3:14"
    const val begin = "Begin Your Journey"
    const val alreadyAccount = "Already have an account?"
    const val signIn = "Sign In"
    const val google = "Sign in with Google"
    const val googleContinue = "Continue with Google"
    const val googleUnavailable = "Google Sign In is not available"
    const val googleFailed = "Google Sign In failed. Try again."
    const val googleNone = "No Google account is available on this phone."
    const val createTitle = "Create your\naccount"
    const val createSubtitle = "Your journey begins here"
    const val ageGate = "I am 18 or older, and I am setting this up for myself."
    const val legalLead = "By continuing you agree to our"
    const val terms = "Terms"
    const val privacy = "Privacy"
    const val getStarted = "Get Started"
    const val creating = "Creating account..."
    const val inviteTitle = "Invite your\naccountability partner"
    const val inviteSkip = "Skip for now"
    const val groupCap = "Groups hold up to 12 people."
    const val sharePartner = "I would like you as my accountability partner on Remain Faithful. https://www.remainfaithful.com"
    const val usageTitle = "Allow usage access"
    const val usageBody = "Remain Faithful reads which apps are open on this phone so it can time the apps you choose. The app name stays on this phone. Your partner receives an alert category, a short summary, and a heartbeat. Photos, messages, web addresses, and the screen stay here."
    const val alertsRow = "Your partner is told after about one minute of use in an app you choose, once a day."
    const val blockingRow = "Chosen apps can be covered by a shield you turn on. You can dismiss it when App Lockout is off."
    const val deepScanTitle = "Deep Scan"
    const val deepScanSubtitle = "Optional. It does not replace Alerts or Blocking."
    const val deepScanSummary = "Deep Scan looks at this phone only while you start it in Remain Faithful. Alerts and Blocking stay as you set them."
    const val deepScanStep = "Start Deep Scan here. Android asks you to allow a screen capture for this session. Frames stay on this phone and are dropped. Nothing from the screen is uploaded."
    const val deepScanLimits = "Protected video such as Netflix is not visible. Stopping or pausing Deep Scan notifies your partners."
    const val setUpLater = "Set up later"
    const val youAreSet = "You are set up"
    const val deepScanHomeNote = "Deep Scan stays off until you start it in the app. Protected video such as Netflix is not visible."
    const val goHome = "Go to Home"
    const val welcomeBack = "Welcome Back"
    const val signInSubtitle = "Sign in to continue your journey"
    const val forgot = "Forgot Password?"
    const val newHere = "New here?"
    const val createAccount = "Create Account"
    const val networkError = "Network error. Check your connection."
    const val pauseNotice = "Pausing or turning off monitoring notifies your partners"
    const val panic = "I Need Support Right Now"
    const val noAlerts = "No alerts yet."
    const val partnerEmpty = "No alerts yet. This does not mean their Blocking is on."
    const val appUsed = "An app you chose was used"
    const val appUsedDetail = "Alerts and Blocking stay as you set them. Deep Scan is optional."
    const val skipTour = "Skip tour"
    const val signOutBody = "Signing out stops monitoring on this phone. Your partners will be told. Blocking on this phone is separate from the account and turns off until you sign back in."
    const val howItWorksPrivacy = "Only an alert and a heartbeat leave this phone."
    const val deleteBody = "This will permanently delete your account, all your data, and notify your accountability partners that you've left. This cannot be undone."
    const val pinTitle = "Partner PIN Required"
    const val pinBody = "Your partner set this PIN. They are the one who knows it."
    const val lockoutTitle = "Keep apps blocked unless Deep Scan is running"
    const val lockoutTitleNoScan = "Keep chosen apps covered by the shield"
    const val shieldLockoutPin = "Ask your partner for the Partner PIN."
    const val pinActive = "Active. Your partner controls gated actions."
    const val pinMissing = "Not set. Ask your partner to enable it in Manage Partners."
    const val overlayTitle = "Allow the shield"
    const val overlayBody = "Remain Faithful can draw a shield over an app you chose to block. The shield does not read the screen. It does not see web addresses or the app name, and it does not send them. You can turn this off in system settings."
    const val notifyTitle = "Allow notifications"
    const val notifyBody = "Remain Faithful uses notifications for partner alerts, the urgent prayer request, and the monitoring reminder while Alerts or Blocking are on."
    const val ideas = "Send Ideas or Report a Problem"
    const val contact = "Contact Support"
    const val rate = "Rate the App"
    const val how = "How It Works"
    const val replay = "Replay App Tour"
    const val morning = "Good morning"
    const val afternoon = "Good afternoon"
    const val evening = "Good evening"
    const val glad = "Glad you are here."
    const val verseLabel = "VERSE OF THE DAY"
    const val shieldLockout = "Open Remain Faithful and start Deep Scan to unlock."
    const val shieldOpen = "Your partner is standing with you."
    const val needHelp = "I need help"
    const val termsUrl = "https://remainfaithful.com/terms"
    const val privacyUrl = "https://remainfaithful.com/privacy"
    const val supportEmail = "support@remainfaithful.com"

    fun howSteps(includeDeepScan: Boolean = true): List<Pair<String, String>> = listOf(
        "Choose Your Partners" to "Invite a trusted friend or join a small group of men who will walk with you. Accountability works because you choose people who care.",
        "Turn on Alerts and Blocking" to if (includeDeepScan) {
            "Choose the apps you want help with. Alerts tell your partner after about a minute of use, once a day. Blocking covers those apps with a shield you asked for. Deep Scan is optional and cannot see protected video such as Netflix. Only an alert and a heartbeat leave this phone."
        } else {
            "Choose the apps you want help with. Alerts tell your partner after about a minute of use, once a day. Blocking covers those apps with a shield you asked for. Only an alert and a heartbeat leave this phone."
        },
        "Stay Accountable" to "Review your activity together and spur one another toward freedom. Iron sharpens iron, and so one man sharpens another.",
    )

    fun forgotSent(email: String): String =
        "If an account exists for $email, you'll get an email with a link to remainfaithful.com/reset-password. Open that page on this phone. The link works for one hour."

    fun allStrings(): List<String> {
        val base = listOf(
            appName, welcomeSubtitle, scripture, scriptureRef, begin, alreadyAccount, signIn,
            google, googleContinue, googleUnavailable, googleFailed, googleNone,
            createTitle, createSubtitle, ageGate, legalLead, terms, privacy, getStarted, creating,
            inviteTitle, inviteSkip, groupCap, sharePartner, usageTitle, usageBody, alertsRow,
            blockingRow, deepScanTitle, deepScanSubtitle, deepScanSummary, deepScanStep,
            deepScanLimits, setUpLater, youAreSet, deepScanHomeNote, goHome, welcomeBack,
            signInSubtitle, forgot, newHere, createAccount, networkError, pauseNotice, panic,
            noAlerts, partnerEmpty, appUsed, appUsedDetail, skipTour, signOutBody,
            howItWorksPrivacy, deleteBody, pinTitle, pinBody, lockoutTitle, lockoutTitleNoScan,
            shieldLockoutPin, pinActive, pinMissing,
            overlayTitle, overlayBody, notifyTitle, notifyBody, ideas, contact, rate, how, replay,
            morning, afternoon, evening, glad, verseLabel, shieldLockout, shieldOpen, needHelp,
        )
        val verses = Verses.all.flatMap { listOf(it.text, it.ref) }
        val starters = Conversation.starters.values.flatten()
        val tour = listOf(true, false).flatMap { deep ->
            listOf(true, false).flatMap { inGroup ->
                TourContent.steps(inGroup, deep).flatMap { listOf(it.title, it.body) }
            }
        }
        val howLines = howSteps(true).flatMap { listOf(it.first, it.second) } +
            howSteps(false).flatMap { listOf(it.first, it.second) }
        return base + verses + starters + tour + howLines
    }
}

data class Verse(val text: String, val ref: String)

object Verses {
    val all = listOf(
        Verse("Watch, stand fast in the faith, be brave, be strong.", "1 Corinthians 16:13"),
        Verse("Flee from sexual immorality. Every other sin a person commits is outside the body, but the sexually immoral person sins against his own body.", "1 Corinthians 6:18"),
        Verse("Submit yourselves therefore to God. Resist the devil, and he will flee from you.", "James 4:7"),
        Verse("No temptation has overtaken you that is not common to man. God is faithful, and he will not let you be tempted beyond your ability.", "1 Corinthians 10:13"),
        Verse("I have stored up your word in my heart, that I might not sin against you.", "Psalm 119:11"),
        Verse("Blessed is the man who remains steadfast under trial, for when he has stood the test he will receive the crown of life.", "James 1:12"),
        Verse("Iron sharpens iron, and one man sharpens another.", "Proverbs 27:17"),
        Verse("Create in me a clean heart, O God, and renew a right spirit within me.", "Psalm 51:10"),
        Verse("Walk by the Spirit, and you will not gratify the desires of the flesh.", "Galatians 5:16"),
        Verse("Put on the whole armor of God, that you may be able to stand against the schemes of the devil.", "Ephesians 6:11"),
        Verse("I can do all things through him who strengthens me.", "Philippians 4:13"),
        Verse("For everyone who has been born of God overcomes the world. And this is the victory that has overcome the world. Our faith.", "1 John 5:4"),
        Verse("Do you not know that your body is a temple of the Holy Spirit within you, whom you have from God? You are not your own.", "1 Corinthians 6:19"),
        Verse("Set your minds on things that are above, not on things that are on earth.", "Colossians 3:2"),
        Verse("Finally, brothers, whatever is true, whatever is honorable, whatever is just, whatever is pure. Think about these things.", "Philippians 4:8"),
        Verse("If we confess our sins, he is faithful and just to forgive us our sins and to cleanse us from all unrighteousness.", "1 John 1:9"),
        Verse("How can a young man keep his way pure? By guarding it according to your word.", "Psalm 119:9"),
        Verse("For the grace of God has appeared, bringing salvation for all people, training us to renounce ungodliness and worldly passions.", "Titus 2:11-12"),
        Verse("Confess your sins to one another and pray for one another, that you may be healed.", "James 5:16"),
        Verse("He who conceals his transgressions will not prosper, but he who confesses and forsakes them will obtain mercy.", "Proverbs 28:13"),
        Verse("Be sober-minded; be watchful. Your adversary the devil prowls around like a roaring lion, seeking someone to devour.", "1 Peter 5:8"),
        Verse("Therefore, since we are surrounded by so great a cloud of witnesses, let us lay aside every weight, and sin which clings so closely.", "Hebrews 12:1"),
        Verse("The Lord will fight for you; you need only to be still.", "Exodus 14:14"),
        Verse("For we do not have a high priest who is unable to sympathize with our weaknesses, but one who in every respect has been tempted as we are, yet without sin.", "Hebrews 4:15"),
        Verse("But I say, walk by the Spirit, and you will not gratify the desires of the flesh.", "Galatians 5:16"),
        Verse("Or do you not know that the unrighteous will not inherit the kingdom of God? Do not be deceived.", "1 Corinthians 6:9"),
        Verse("The name of the Lord is a strong tower; the righteous man runs into it and is safe.", "Proverbs 18:10"),
        Verse("So flee youthful passions and pursue righteousness, faith, love, and peace, along with those who call on the Lord from a pure heart.", "2 Timothy 2:22"),
        Verse("For the weapons of our warfare are not of the flesh but have divine power to destroy strongholds.", "2 Corinthians 10:4"),
        Verse("And we all, with unveiled face, beholding the glory of the Lord, are being transformed into the same image from one degree of glory to another.", "2 Corinthians 3:18"),
    )

    fun forDayOfYear(dayOfYear: Int): Verse {
        val index = ((dayOfYear - 1).coerceAtLeast(0)) % all.size
        return all[index]
    }
}

object Conversation {
    val starters: Map<String, List<String>> = mapOf(
        "adult_content" to listOf(
            "What was going on in your heart before this happened? Were you stressed, lonely, or bored?",
            "Was this something you sought out, or did you stumble onto it? What happened in that moment?",
            "What does your relationship with purity look like right now? What would help most this week?",
            "How can I be praying for you specifically in this area?",
        ),
        "gambling" to listOf(
            "Are there financial pressures or anxieties you have been carrying lately that we should talk about?",
            "What draws you to this: the thrill, a sense of control, or something else entirely?",
            "Have you spoken with a pastor or counselor about this pattern? Would that be helpful?",
            "What would one concrete boundary look like here, and how can I help hold you to it?",
        ),
        "violence" to listOf(
            "What were you looking at, and how are you doing right now?",
            "Was this something you went looking for, or did it show up in front of you?",
            "What would help you step away from this the next time it appears?",
            "How can I check in with you later today?",
        ),
        "self_harm" to listOf(
            "Are you safe right now? I want to hear that from you.",
            "You don't have to explain everything. What do you need in the next hour?",
            "Who else can be with you: me, a pastor, or a counselor?",
            "Can we stay in touch until this feels smaller?",
        ),
        "app_usage" to listOf(
            "You opened an app you asked to be watched. What was going on just before?",
            "Was this a slip, or a choice you had already talked about?",
            "What would make the next hour easier?",
            "Do you want Blocking on for this app, or is an Alert enough?",
        ),
        "clean" to listOf(
            "What's been helping this week?",
            "Is there a pressure point we should name before it grows?",
            "How can I pray for you specifically?",
        ),
        "protection" to listOf(
            "Monitoring changed on the phone. What happened?",
            "Was this intentional, or did something else turn it off?",
            "What do you want back on: Alerts, Blocking, or both?",
            "When can we confirm it's on again?",
        ),
    )

    fun forCategory(category: String): List<String> =
        starters[category] ?: listOf(
            "What does this alert mean on your end?",
            "What was happening when it arrived?",
            "What would a helpful next step look like today?",
        )
}

object CategoryLabels {
    fun label(api: String): String = when (api.lowercase().replace(' ', '_')) {
        "adult_content" -> "Adult Content"
        "gambling" -> "Gambling"
        "violence" -> "Violence"
        "self_harm" -> "Self-Harm"
        "app_usage" -> "App Use"
        "clean" -> "All Clear"
        "monitoring_disabled", "shielding_disabled", "signed_out",
        "family_controls_revoked", "deep_scan_stopped", "lockout_disabled",
        "selection_cleared", "heartbeat_silence", "wrong_pin_attempt",
        "pin_removed", "pin_changed", "protection", "tamper",
        -> "Monitoring Change"
        else -> "Alert"
    }

    fun starterKey(api: String): String = when (api.lowercase().replace(' ', '_')) {
        "adult_content", "gambling", "violence", "self_harm", "app_usage", "clean" ->
            api.lowercase().replace(' ', '_')
        "monitoring_disabled", "shielding_disabled", "signed_out",
        "family_controls_revoked", "deep_scan_stopped", "lockout_disabled",
        "selection_cleared", "heartbeat_silence", "wrong_pin_attempt",
        "pin_removed", "pin_changed", "protection", "tamper",
        -> "protection"
        else -> "unknown"
    }
}
