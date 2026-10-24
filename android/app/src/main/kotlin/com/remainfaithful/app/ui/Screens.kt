package com.remainfaithful.app.ui

import android.content.Intent
import android.net.Uri
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.heightIn
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Group
import androidx.compose.material.icons.filled.Home
import androidx.compose.material.icons.filled.Settings
import androidx.compose.material3.Checkbox
import androidx.compose.material3.CheckboxDefaults
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.Icon
import androidx.compose.material3.NavigationBar
import androidx.compose.material3.NavigationBarItem
import androidx.compose.material3.NavigationBarItemDefaults
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Switch
import androidx.compose.material3.SwitchDefaults
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.font.FontFamily
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.remainfaithful.app.theme.RfButton
import com.remainfaithful.app.theme.RfColors
import com.remainfaithful.app.theme.RfMetrics
import com.remainfaithful.app.theme.SectionLabel
import com.remainfaithful.app.theme.rfCard
import com.remainfaithful.app.theme.rfFieldColors
import com.remainfaithful.core.AuthRules
import com.remainfaithful.core.CategoryLabels
import com.remainfaithful.core.Conversation
import com.remainfaithful.core.Copy
import com.remainfaithful.core.FeedbackRules
import com.remainfaithful.core.GroupPolicy
import com.remainfaithful.core.PinCheck
import com.remainfaithful.core.StreakMark

@Composable
fun RootScreen(model: AppModel) {
    Box(Modifier.fillMaxSize().background(RfColors.page)) {
        when (model.phase) {
            Phase.Loading -> CircularProgressIndicator(color = RfColors.gold, modifier = Modifier.align(Alignment.Center))
            Phase.Onboarding -> Onboarding(model)
            Phase.Login -> LoginScreen(model)
            Phase.App -> MainTabs(model)
        }
        model.banner?.let { note ->
            Text(
                note,
                color = RfColors.textPrimary,
                modifier = Modifier
                    .align(Alignment.BottomCenter)
                    .padding(16.dp)
                    .rfCard()
                    .clickable { model.banner = null },
            )
        }
        if (model.sheet != null) SheetHost(model)
        if (model.tourOn && model.phase == Phase.App) TourCard(model)
    }
}

@Composable
private fun Onboarding(model: AppModel) {
    Column(Modifier.fillMaxSize().verticalScroll(rememberScrollState()).padding(RfMetrics.gutter)) {
        if (model.step in 1..6) {
            Text("${model.step} of 6", color = RfColors.gold, fontSize = 13.sp)
            Spacer(Modifier.height(12.dp))
        }
        when (model.step) {
            0 -> Welcome(model)
            1 -> CreateAccount(model)
            2 -> InviteStep(model)
            3 -> DisclosureStep(model, Copy.usageTitle, Copy.usageBody, Copy.alertsRow, Copy.blockingRow, "Allow usage access") {
                model.askSystem("usage")
            }
            4 -> DeepScanStep(model)
            5 -> ChooseStep(model)
            else -> DoneStep(model)
        }
    }
}

@Composable
private fun Welcome(model: AppModel) {
    Spacer(Modifier.height(48.dp))
    Text(Copy.appName, color = RfColors.textPrimary, fontFamily = FontFamily.Serif, fontWeight = FontWeight.Bold, fontSize = 40.sp)
    Spacer(Modifier.height(8.dp))
    Text(Copy.welcomeSubtitle, color = RfColors.textSecondary, fontSize = 18.sp)
    Spacer(Modifier.height(28.dp))
    Text(Copy.scripture, color = RfColors.textPrimary, fontFamily = FontFamily.Serif, fontSize = 16.sp)
    Text(Copy.scriptureRef, color = RfColors.gold, modifier = Modifier.padding(top = 8.dp))
    Spacer(Modifier.height(36.dp))
    RfButton(Copy.begin, onClick = { model.step = 1 })
    TextButtonRow(Copy.alreadyAccount, Copy.signIn) { model.phase = Phase.Login }
}

@Composable
private fun CreateAccount(model: AppModel) {
    Text(Copy.createTitle, color = RfColors.textPrimary, fontFamily = FontFamily.Serif, fontWeight = FontWeight.Bold, fontSize = 34.sp)
    Text(Copy.createSubtitle, color = RfColors.textSecondary, modifier = Modifier.padding(top = 8.dp, bottom = 16.dp))
    Field("Name", model.formName) { model.formName = it }
    Field("Email", model.formEmail, KeyboardType.Email) { model.formEmail = it }
    Field("Password", model.formPassword, KeyboardType.Password, password = true) { model.formPassword = it }
    Text(AuthRules.passwordHint(model.formPassword), color = RfColors.textTertiary, fontSize = 13.sp)
    Row(verticalAlignment = Alignment.CenterVertically, modifier = Modifier.padding(vertical = 8.dp)) {
        Checkbox(
            checked = model.adult,
            onCheckedChange = { model.adult = it },
            colors = CheckboxDefaults.colors(checkedColor = RfColors.gold, checkmarkColor = RfColors.navy),
        )
        Text(Copy.ageGate, color = RfColors.textPrimary, modifier = Modifier.weight(1f))
    }
    val context = LocalContext.current
    Text(Copy.legalLead, color = RfColors.textTertiary, fontSize = 13.sp)
    Row {
        Text(Copy.terms, color = RfColors.gold, modifier = Modifier.clickable {
            context.startActivity(Intent(Intent.ACTION_VIEW, Uri.parse(Copy.termsUrl)))
        }.padding(end = 16.dp, top = 8.dp, bottom = 8.dp))
        Text(Copy.privacy, color = RfColors.gold, modifier = Modifier.clickable {
            context.startActivity(Intent(Intent.ACTION_VIEW, Uri.parse(Copy.privacyUrl)))
        }.padding(top = 8.dp, bottom = 8.dp))
    }
    RfButton(if (model.busy) Copy.creating else Copy.getStarted, onClick = { model.register() }, enabled = !model.busy)
}

@Composable
private fun InviteStep(model: AppModel) {
    var address by remember { mutableStateOf("") }
    Text(Copy.inviteTitle, color = RfColors.textPrimary, fontFamily = FontFamily.Serif, fontWeight = FontWeight.Bold, fontSize = 34.sp)
    Text("They accept in Settings. Or share a link they can open. They still accept from the email you send. ${Copy.groupCap}", color = RfColors.textSecondary, modifier = Modifier.padding(vertical = 12.dp))
    Field("Partner email", address, KeyboardType.Email) { address = it }
    RfButton("Send invite", onClick = { if (address.contains("@")) model.invitePartner(address) })
    val context = LocalContext.current
    RfButton("Share a link", onClick = {
        val send = Intent(Intent.ACTION_SEND).setType("text/plain").putExtra(Intent.EXTRA_TEXT, Copy.sharePartner)
        context.startActivity(Intent.createChooser(send, "Share"))
    }, secondary = true)
    RfButton(Copy.inviteSkip, onClick = { model.step = 3 }, secondary = true)
}

@Composable
private fun DisclosureStep(
    model: AppModel,
    title: String,
    body: String,
    row1: String,
    row2: String,
    action: String,
    onAllow: () -> Unit,
) {
    Text(title, color = RfColors.textPrimary, fontFamily = FontFamily.Serif, fontWeight = FontWeight.Bold, fontSize = 32.sp)
    Text(body, color = RfColors.textPrimary, modifier = Modifier.padding(vertical = 16.dp))
    Text("Alerts", color = RfColors.gold, fontWeight = FontWeight.Bold)
    Text(row1, color = RfColors.textSecondary, modifier = Modifier.padding(bottom = 12.dp))
    Text("Blocking", color = RfColors.gold, fontWeight = FontWeight.Bold)
    Text(row2, color = RfColors.textSecondary, modifier = Modifier.padding(bottom = 16.dp))
    RfButton(action, onClick = onAllow)
    RfButton(Copy.inviteSkip, onClick = { model.step = 4 }, secondary = true)
}

@Composable
private fun DeepScanStep(model: AppModel) {
    Text(Copy.deepScanTitle, color = RfColors.textPrimary, fontFamily = FontFamily.Serif, fontWeight = FontWeight.Bold, fontSize = 32.sp)
    Text(Copy.deepScanSubtitle, color = RfColors.gold, modifier = Modifier.padding(vertical = 8.dp))
    Text(Copy.deepScanSummary, color = RfColors.textPrimary)
    Text(Copy.deepScanStep, color = RfColors.textSecondary, modifier = Modifier.padding(vertical = 12.dp))
    Text(Copy.deepScanLimits, color = RfColors.textSecondary)
    Spacer(Modifier.height(16.dp))
    RfButton("Continue", onClick = { model.step = 5 })
    RfButton(Copy.setUpLater, onClick = { model.step = 5 }, secondary = true)
}

@Composable
private fun ChooseStep(model: AppModel) {
    Text("Choose apps", color = RfColors.textPrimary, fontFamily = FontFamily.Serif, fontWeight = FontWeight.Bold, fontSize = 32.sp)
    Text("Pick the apps that tempt you. The names stay on this phone.", color = RfColors.textSecondary, modifier = Modifier.padding(vertical = 8.dp))
    if (model.apps.isEmpty()) {
        RfButton("Show my apps", onClick = { model.loadApps() })
    }
    model.apps.take(40).forEach { app ->
        val on = app.packageName in model.chosen
        Row(
            Modifier.fillMaxWidth().heightIn(min = RfMetrics.minTap).clickable {
                model.saveChosen(if (on) model.chosen - app.packageName else model.chosen + app.packageName)
            },
            verticalAlignment = Alignment.CenterVertically,
        ) {
            Checkbox(checked = on, onCheckedChange = {
                model.saveChosen(if (it) model.chosen + app.packageName else model.chosen - app.packageName)
            }, colors = CheckboxDefaults.colors(checkedColor = RfColors.gold, checkmarkColor = RfColors.navy))
            Text(app.label, color = RfColors.textPrimary)
        }
    }
    RfButton("Allow notifications", onClick = { model.askSystem("notifications") }, secondary = true)
    RfButton("I'll do this later", onClick = { model.step = 6 }, secondary = true)
    RfButton("Continue", onClick = { model.step = 6 })
}

@Composable
private fun DoneStep(model: AppModel) {
    Text(Copy.youAreSet, color = RfColors.textPrimary, fontFamily = FontFamily.Serif, fontWeight = FontWeight.Bold, fontSize = 34.sp)
    Text(Copy.deepScanHomeNote, color = RfColors.textSecondary, modifier = Modifier.padding(vertical = 16.dp))
    RfButton(Copy.goHome, onClick = { model.finishOnboarding() })
}

@Composable
private fun LoginScreen(model: AppModel) {
    Column(Modifier.fillMaxSize().verticalScroll(rememberScrollState()).padding(RfMetrics.gutter)) {
        Spacer(Modifier.height(48.dp))
        Text(Copy.welcomeBack, color = RfColors.textPrimary, fontFamily = FontFamily.Serif, fontWeight = FontWeight.Bold, fontSize = 36.sp)
        Text(Copy.signInSubtitle, color = RfColors.textSecondary, modifier = Modifier.padding(bottom = 20.dp))
        Field("Email", model.formEmail, KeyboardType.Email) { model.formEmail = it }
        Field("Password", model.formPassword, KeyboardType.Password, password = true) { model.formPassword = it }
        Text(Copy.forgot, color = RfColors.gold, modifier = Modifier.clickable { model.sheet = Sheet.Forgot }.padding(vertical = 8.dp))
        RfButton(Copy.signIn, onClick = { model.login() }, enabled = !model.busy)
        TextButtonRow(Copy.newHere, Copy.createAccount) {
            model.phase = Phase.Onboarding
            model.step = 1
        }
    }
}

@Composable
private fun MainTabs(model: AppModel) {
    Column(Modifier.fillMaxSize()) {
        Box(Modifier.weight(1f)) {
            when (model.tab) {
                0 -> HomeScreen(model)
                1 -> GroupScreen(model)
                else -> SettingsScreen(model)
            }
        }
        NavigationBar(containerColor = RfColors.navy) {
            val colors = NavigationBarItemDefaults.colors(
                selectedIconColor = RfColors.navy,
                selectedTextColor = RfColors.gold,
                indicatorColor = RfColors.gold,
                unselectedIconColor = RfColors.textTertiary,
                unselectedTextColor = RfColors.textTertiary,
            )
            NavigationBarItem(selected = model.tab == 0, onClick = { model.tab = 0 }, icon = { Icon(Icons.Filled.Home, null) }, label = { Text("Home") }, colors = colors)
            NavigationBarItem(selected = model.tab == 1, onClick = { model.tab = 1 }, icon = { Icon(Icons.Filled.Group, null) }, label = { Text("Group") }, colors = colors)
            NavigationBarItem(selected = model.tab == 2, onClick = { model.tab = 2 }, icon = { Icon(Icons.Filled.Settings, null) }, label = { Text("Settings") }, colors = colors)
        }
    }
}

@Composable
private fun HomeScreen(model: AppModel) {
    LazyColumn(Modifier.fillMaxSize().padding(horizontal = RfMetrics.gutter), verticalArrangement = Arrangement.spacedBy(14.dp)) {
        item { Spacer(Modifier.height(18.dp)) }
        item {
            Text(model.greeting(), color = RfColors.textPrimary, fontFamily = FontFamily.Serif, fontWeight = FontWeight.Bold, fontSize = 32.sp)
            Text(Copy.glad, color = RfColors.textSecondary)
            if (!model.verified) {
                Text("Confirm your email.", color = RfColors.gold, modifier = Modifier.clickable { model.resend() }.padding(top = 8.dp))
            }
        }
        item {
            Column(Modifier.rfCard()) {
                Text(model.status.title, color = RfColors.textPrimary, fontWeight = FontWeight.Bold, fontSize = 18.sp)
                Text(model.status.detail, color = RfColors.textSecondary, modifier = Modifier.padding(top = 6.dp))
                Text(model.status.deepScanLine, color = RfColors.gold, fontSize = 13.sp, modifier = Modifier.padding(top = 8.dp))
                model.status.actionTitle?.let { title ->
                    Spacer(Modifier.height(10.dp))
                    RfButton(title, onClick = { model.openMonitoring() })
                }
            }
        }
        item {
            Column(Modifier.rfCard()) {
                Text("$model.streak".let { model.streak.toString() }, color = RfColors.textPrimary, fontSize = 42.sp, fontWeight = FontWeight.Bold, fontFamily = FontFamily.Serif, style = androidx.compose.ui.text.TextStyle(fontFeatureSettings = "tnum"))
                Text("days you've stayed with it", color = RfColors.textSecondary)
                Text("Best: ${model.best}d", color = RfColors.gold, fontSize = 13.sp)
                Row(Modifier.padding(top = 10.dp), horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                    model.week.forEach { mark ->
                        val color = when (mark) {
                            StreakMark.Clean -> RfColors.success
                            StreakMark.Flagged -> RfColors.flag
                            StreakMark.BeforeSignup -> RfColors.hairline
                        }
                        Box(Modifier.size(18.dp).clip(CircleShape).background(color))
                    }
                }
            }
        }
        item {
            val verse = model.verse()
            Column(Modifier.rfCard()) {
                Text(Copy.verseLabel, color = RfColors.gold, fontSize = 11.sp, fontWeight = FontWeight.Bold, letterSpacing = 1.2.sp)
                Text(verse.text, color = RfColors.textPrimary, fontFamily = FontFamily.Serif, modifier = Modifier.padding(top = 8.dp))
                Text(verse.ref, color = RfColors.gold, modifier = Modifier.padding(top = 8.dp))
            }
        }
        item { SectionLabel("ALERTS") }
        if (model.alerts.isEmpty()) item { Text(Copy.partnerEmpty, color = RfColors.textSecondary) }
        items(model.alerts, key = { it.key }) { event -> EventRow(event) { model.openAlert(event) } }
        item { SectionLabel("YOUR ACTIVITY") }
        if (model.mine.isEmpty()) item { Text(Copy.noAlerts, color = RfColors.textSecondary) }
        items(model.mine, key = { it.key }) { event -> EventRow(event) { model.openAlert(event) } }
        item {
            Text(Copy.pauseNotice, color = RfColors.textTertiary, fontSize = 13.sp)
            Spacer(Modifier.height(8.dp))
            RfButton(Copy.panic, onClick = { model.sheet = Sheet.Panic }, destructive = true)
            Spacer(Modifier.height(24.dp))
        }
    }
}

@Composable
private fun EventRow(event: UiEvent, onClick: () -> Unit) {
    Column(Modifier.fillMaxWidth().rfCard().clickable(onClick = onClick)) {
        Text(CategoryLabels.label(event.category), color = RfColors.gold, fontWeight = FontWeight.Bold)
        Text(event.summary, color = RfColors.textPrimary, modifier = Modifier.padding(top = 4.dp))
        Text(event.timeLabel, color = RfColors.textTertiary, fontSize = 12.sp, modifier = Modifier.padding(top = 4.dp))
    }
}

@Composable
private fun GroupScreen(model: AppModel) {
    var groupName by remember { mutableStateOf("") }
    var invite by remember { mutableStateOf("") }
    var covenant by remember(model.group?.id) { mutableStateOf(model.group?.covenant.orEmpty()) }
    val context = LocalContext.current
    LazyColumn(Modifier.fillMaxSize().padding(horizontal = RfMetrics.gutter), verticalArrangement = Arrangement.spacedBy(12.dp)) {
        item { Spacer(Modifier.height(18.dp)) }
        item { Text("Group", color = RfColors.textPrimary, fontFamily = FontFamily.Serif, fontWeight = FontWeight.Bold, fontSize = 32.sp) }
        if (model.groupInvites.isNotEmpty()) {
            item { SectionLabel("INVITES") }
            items(model.groupInvites, key = { it.id }) { inviteRow ->
                Column(Modifier.rfCard()) {
                    Text("${inviteRow.inviterName} invited you to ${inviteRow.groupName}.", color = RfColors.textPrimary)
                    Text("You join after you accept.", color = RfColors.textSecondary, fontSize = 13.sp)
                    Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                        RfButton("Accept", onClick = { model.acceptGroup(inviteRow.id) }, modifier = Modifier.weight(1f))
                        RfButton("Decline", onClick = { model.declineGroup(inviteRow.id) }, secondary = true, modifier = Modifier.weight(1f))
                    }
                }
            }
        }
        val current = model.group
        if (current == null) {
            item {
                Text("Start a group here. A name is enough. ${Copy.groupCap}", color = RfColors.textSecondary)
                Field("Group name", groupName) { groupName = it }
                RfButton("Create group", onClick = { if (groupName.isNotBlank()) model.createGroup(groupName) }, enabled = groupName.isNotBlank())
            }
        } else {
            item {
                Text(current.name, color = RfColors.textPrimary, fontWeight = FontWeight.Bold, fontSize = 22.sp)
                Text("${current.members.size} of ${GroupPolicy.maxMembers}", color = RfColors.gold, fontSize = 13.sp)
            }
            items(current.members, key = { it.userId }) { member ->
                Column(Modifier.rfCard()) {
                    Text(member.name, color = RfColors.textPrimary, fontWeight = FontWeight.SemiBold)
                    Text("${member.streakDays} day streak. ${member.flagsLast30} flags in 30 days.", color = RfColors.textSecondary, fontSize = 13.sp)
                    if (member.userId != model.userId) {
                        RfButton("Encourage", onClick = { model.encourage(member) }, secondary = true)
                    }
                }
            }
            item {
                SectionLabel("COVENANT")
                Field("Covenant", covenant) { covenant = it }
                RfButton("Save covenant", onClick = { model.updateCovenant(covenant) }, secondary = true)
            }
            item {
                val full = GroupPolicy.fullMessage(current.members.size)
                if (full != null) Text(full, color = RfColors.gold)
                Field("Member email", invite, KeyboardType.Email) { invite = it }
                RfButton("Invite by email", onClick = { model.inviteMember(invite) }, enabled = full == null && invite.contains("@"))
                RfButton("Share link", onClick = {
                    val send = Intent(Intent.ACTION_SEND).setType("text/plain").putExtra(Intent.EXTRA_TEXT, model.inviteLink())
                    context.startActivity(Intent.createChooser(send, "Share"))
                }, secondary = true)
                RfButton("Leave group", onClick = { model.leaveGroup() }, destructive = true)
            }
        }
        item { Spacer(Modifier.height(24.dp)) }
    }
}

@Composable
private fun SettingsScreen(model: AppModel) {
    LazyColumn(Modifier.fillMaxSize().padding(horizontal = RfMetrics.gutter), verticalArrangement = Arrangement.spacedBy(10.dp)) {
        item { Spacer(Modifier.height(18.dp)) }
        item {
            Text("Settings", color = RfColors.textPrimary, fontFamily = FontFamily.Serif, fontWeight = FontWeight.Bold, fontSize = 32.sp)
            Text(model.name, color = RfColors.textPrimary, modifier = Modifier.padding(top = 8.dp))
            Text(model.email, color = RfColors.textSecondary)
            if (!model.verified) Text("Confirm your email.", color = RfColors.gold, modifier = Modifier.clickable { model.resend() })
        }
        item { SectionLabel("ACCOUNTABILITY") }
        item { SettingsRow("Manage Partners") { model.sheet = Sheet.Partners } }
        item { SettingsRow("Edit Profile") { model.sheet = Sheet.Profile } }
        item { SectionLabel("MONITORING") }
        item { SettingsRow("Notifications") { model.askSystem("notifications") } }
        item { SectionLabel("BLOCKING") }
        item { SettingsRow("Blocking") { model.openMonitoring() } }
        item { SectionLabel("PROTECTION") }
        item {
            Row(Modifier.rfCard(), verticalAlignment = Alignment.CenterVertically) {
                Column(Modifier.weight(1f)) {
                    Text(Copy.lockoutTitle, color = RfColors.textPrimary)
                }
                Switch(
                    checked = model.lockoutOn,
                    onCheckedChange = { model.setLockout(it) },
                    colors = SwitchDefaults.colors(checkedTrackColor = RfColors.gold, checkedThumbColor = RfColors.navy),
                )
            }
        }
        item {
            Text("Partner PIN", color = RfColors.textPrimary, fontWeight = FontWeight.SemiBold)
            Text(if (model.pinSet) Copy.pinActive else Copy.pinMissing, color = RfColors.textSecondary)
        }
        item { SectionLabel("PRIVACY") }
        item { SettingsRow("View My Activity Log") { model.sheet = Sheet.Activity } }
        item { SectionLabel("APP") }
        item { SettingsRow(Copy.how) { model.sheet = Sheet.How } }
        item { SettingsRow(Copy.replay) { model.replayTour() } }
        item {
            val context = LocalContext.current
            SettingsRow(Copy.rate) {
                val market = Intent(Intent.ACTION_VIEW, Uri.parse("market://details?id=${context.packageName}"))
                runCatching { context.startActivity(market) }.onFailure {
                    context.startActivity(Intent(Intent.ACTION_VIEW, Uri.parse("https://play.google.com/store/apps/details?id=${context.packageName}")))
                }
            }
        }
        item { SettingsRow(Copy.ideas) { model.sheet = Sheet.Feedback } }
        item { SectionLabel("DANGER ZONE") }
        item { SettingsRow("Leave All Groups") { model.leaveAll() } }
        item { SettingsRow("Delete Account") { model.sheet = Sheet.Delete } }
        item {
            var confirm by remember { mutableStateOf(false) }
            if (!confirm) {
                RfButton("Sign Out", onClick = { confirm = true }, destructive = true)
            } else {
                Text(Copy.signOutBody, color = RfColors.textSecondary)
                RfButton("Sign Out", onClick = { model.signOut() }, destructive = true)
                RfButton("Cancel", onClick = { confirm = false }, secondary = true)
            }
        }
        item { Spacer(Modifier.height(24.dp)) }
    }
}

@Composable
private fun SettingsRow(label: String, onClick: () -> Unit) {
    Text(
        label,
        color = RfColors.textPrimary,
        modifier = Modifier.fillMaxWidth().rfCard().clickable(onClick = onClick).heightIn(min = RfMetrics.minTap),
    )
}

@Composable
private fun SheetHost(model: AppModel) {
    Box(Modifier.fillMaxSize().background(RfColors.scrim).clickable { model.sheet = null }) {
        Column(
            Modifier
                .align(Alignment.BottomCenter)
                .fillMaxWidth()
                .clip(RoundedCornerShape(topStart = 20.dp, topEnd = 20.dp))
                .background(RfColors.navyMid)
                .clickable(enabled = false) {}
                .verticalScroll(rememberScrollState())
                .padding(20.dp),
        ) {
            when (model.sheet) {
                Sheet.Panic -> PanicSheet(model)
                Sheet.Feedback -> FeedbackSheet(model)
                Sheet.Pin -> PinSheet(model)
                Sheet.Activity -> ActivitySheet(model)
                Sheet.How -> HowSheet { model.sheet = null }
                Sheet.Delete -> DeleteSheet(model)
                Sheet.Partners -> PartnersSheet(model)
                Sheet.Monitoring -> MonitoringSheet(model)
                Sheet.DeepScan -> DeepScanManage(model)
                Sheet.Alert -> AlertSheet(model)
                Sheet.Profile -> ProfileSheet(model)
                Sheet.Forgot -> ForgotSheet(model)
                null -> Unit
            }
        }
    }
}

@Composable
private fun PanicSheet(model: AppModel) {
    Text(Copy.panic, color = RfColors.textPrimary, fontFamily = FontFamily.Serif, fontWeight = FontWeight.Bold, fontSize = 26.sp)
    Text("This sends an urgent notice to your partners and your group.", color = RfColors.textSecondary, modifier = Modifier.padding(vertical = 12.dp))
    RfButton("Send", onClick = { model.panic() }, destructive = true)
    RfButton("Close", onClick = { model.sheet = null }, secondary = true)
}

@Composable
private fun FeedbackSheet(model: AppModel) {
    var type by remember { mutableStateOf(FeedbackRules.types.first()) }
    var summary by remember { mutableStateOf("") }
    var details by remember { mutableStateOf("") }
    var device by remember { mutableStateOf(false) }
    Text(Copy.ideas, color = RfColors.textPrimary, fontFamily = FontFamily.Serif, fontWeight = FontWeight.Bold, fontSize = 24.sp)
    Row(horizontalArrangement = Arrangement.spacedBy(8.dp), modifier = Modifier.padding(vertical = 8.dp)) {
        FeedbackRules.types.forEach { item ->
            Text(
                item,
                color = if (item == type) RfColors.navy else RfColors.textPrimary,
                modifier = Modifier
                    .clip(RoundedCornerShape(12.dp))
                    .background(if (item == type) RfColors.gold else RfColors.card)
                    .clickable { type = item }
                    .padding(horizontal = 12.dp, vertical = 8.dp),
            )
        }
    }
    Field("Summary", summary) { summary = it }
    Field("Details", details) { details = it }
    Row(verticalAlignment = Alignment.CenterVertically) {
        Checkbox(checked = device, onCheckedChange = { device = it }, colors = CheckboxDefaults.colors(checkedColor = RfColors.gold, checkmarkColor = RfColors.navy))
        Text("Include app version and device", color = RfColors.textSecondary)
    }
    RfButton("Send", onClick = { model.feedback(type, summary, details, device) })
}

@Composable
private fun PinSheet(model: AppModel) {
    var pin by remember { mutableStateOf("") }
    var note by remember { mutableStateOf<String?>(null) }
    Text(Copy.pinTitle, color = RfColors.textPrimary, fontFamily = FontFamily.Serif, fontWeight = FontWeight.Bold, fontSize = 24.sp)
    Text(Copy.pinBody, color = RfColors.textSecondary, modifier = Modifier.padding(vertical = 8.dp))
    Field("4 digit PIN", pin, KeyboardType.NumberPassword, password = true) { if (it.length <= 4) pin = it.filter { c -> c.isDigit() } }
    note?.let { Text(it, color = RfColors.flag) }
    RfButton("Continue", onClick = {
        model.submitPin(pin) { check: PinCheck -> note = check.message.ifBlank { null } }
    }, enabled = pin.length == 4)
}

@Composable
private fun ActivitySheet(model: AppModel) {
    Text("Activity", color = RfColors.textPrimary, fontFamily = FontFamily.Serif, fontWeight = FontWeight.Bold, fontSize = 24.sp)
    model.mine.forEach { event ->
        Text("${CategoryLabels.label(event.category)}: ${event.summary}", color = RfColors.textPrimary, modifier = Modifier.padding(vertical = 6.dp))
    }
    if (model.mine.isEmpty()) Text(Copy.noAlerts, color = RfColors.textSecondary)
}

@Composable
private fun HowSheet(onClose: () -> Unit) {
    Text(Copy.how, color = RfColors.textPrimary, fontFamily = FontFamily.Serif, fontWeight = FontWeight.Bold, fontSize = 24.sp)
    Text("Three steps to freedom", color = RfColors.gold, modifier = Modifier.padding(bottom = 8.dp))
    Copy.howSteps.forEach { (title, body) ->
        Text(title, color = RfColors.gold, fontWeight = FontWeight.Bold, modifier = Modifier.padding(top = 10.dp))
        Text(body, color = RfColors.textPrimary)
    }
    Spacer(Modifier.height(12.dp))
    RfButton("Close", onClick = onClose, secondary = true)
}

@Composable
private fun DeleteSheet(model: AppModel) {
    var confirm by remember { mutableStateOf("") }
    Text("Delete Account", color = RfColors.textPrimary, fontFamily = FontFamily.Serif, fontWeight = FontWeight.Bold, fontSize = 24.sp)
    Text(Copy.deleteBody, color = RfColors.textSecondary, modifier = Modifier.padding(vertical = 8.dp))
    Field("Type DELETE to confirm", confirm) { confirm = it }
    RfButton("Delete Account", onClick = { model.deleteAccount(confirm) }, destructive = true, enabled = confirm == "DELETE")
    RfButton("Cancel", onClick = { model.sheet = null }, secondary = true)
}

@Composable
private fun PartnersSheet(model: AppModel) {
    var email by remember { mutableStateOf("") }
    var pinFor by remember { mutableStateOf<Long?>(null) }
    var pin by remember { mutableStateOf("") }
    Text("Partners", color = RfColors.textPrimary, fontFamily = FontFamily.Serif, fontWeight = FontWeight.Bold, fontSize = 24.sp)
    model.partnerInvites.forEach { invite ->
        Column(Modifier.padding(vertical = 6.dp)) {
            Text("${invite.inviterName} invited you.", color = RfColors.textPrimary)
            Row {
                RfButton("Accept", onClick = { model.acceptPartner(invite.id) }, modifier = Modifier.weight(1f))
                Spacer(Modifier.width(8.dp))
                RfButton("Decline", onClick = { model.declinePartner(invite.id) }, secondary = true, modifier = Modifier.weight(1f))
            }
        }
    }
    if (model.rels.isEmpty()) Text("No partners yet. Invite someone below. They must accept before alerts start.", color = RfColors.textSecondary)
    model.rels.forEach { rel ->
        Column(Modifier.padding(vertical = 8.dp)) {
            Text(rel.partnerName, color = RfColors.textPrimary, fontWeight = FontWeight.SemiBold)
            Text(rel.status, color = RfColors.textTertiary, fontSize = 12.sp)
            if (rel.isPrimary) Text("Primary", color = RfColors.gold, fontSize = 12.sp)
            if (!rel.isPrimary && rel.userId == model.userId) {
                RfButton("Make primary", onClick = { model.makePrimary(rel) }, secondary = true)
            }
            if (rel.partnerId == model.userId && rel.status == "accepted") {
                if (pinFor == rel.id) {
                    Field("4 digit PIN", pin, KeyboardType.NumberPassword, password = true) { if (it.length <= 4) pin = it.filter(Char::isDigit) }
                    RfButton("Save PIN", onClick = { model.setPartnerPin(rel.id, pin); pinFor = null }, enabled = pin.length == 4)
                } else {
                    RfButton(if (rel.pinSet) "Change Partner PIN" else "Set Partner PIN", onClick = { pinFor = rel.id }, secondary = true)
                }
                if (rel.pinSet) RfButton("Remove PIN", onClick = { model.removePartnerPin(rel.id) }, secondary = true)
            }
            RfButton("Remove", onClick = { model.removePartner(rel) }, destructive = true)
        }
    }
    Field("Partner email", email, KeyboardType.Email) { email = it }
    RfButton("Invite", onClick = { model.invitePartner(email) }, enabled = email.contains("@"))
}

@Composable
private fun MonitoringSheet(model: AppModel) {
    Text("Blocking", color = RfColors.textPrimary, fontFamily = FontFamily.Serif, fontWeight = FontWeight.Bold, fontSize = 24.sp)
    Text(Copy.pauseNotice, color = RfColors.textTertiary, fontSize = 13.sp, modifier = Modifier.padding(vertical = 8.dp))
    if (!model.usageGranted) {
        Text(Copy.usageBody, color = RfColors.textPrimary)
        RfButton(Copy.usageTitle, onClick = { model.askSystem("usage") })
    }
    if (!model.overlayGranted) {
        Text(Copy.overlayBody, color = RfColors.textSecondary, modifier = Modifier.padding(vertical = 8.dp))
        RfButton(Copy.overlayTitle, onClick = { model.askSystem("overlay") }, secondary = true)
    }
    Row(verticalAlignment = Alignment.CenterVertically) {
        Text("Alerts", color = RfColors.textPrimary, modifier = Modifier.weight(1f))
        Switch(checked = model.alertsOn, onCheckedChange = { model.setAlerts(it) }, colors = SwitchDefaults.colors(checkedTrackColor = RfColors.gold, checkedThumbColor = RfColors.navy))
    }
    Row(verticalAlignment = Alignment.CenterVertically) {
        Text("Blocking", color = RfColors.textPrimary, modifier = Modifier.weight(1f))
        Switch(checked = model.blockingOn, onCheckedChange = { model.setBlocking(it) }, colors = SwitchDefaults.colors(checkedTrackColor = RfColors.gold, checkedThumbColor = RfColors.navy))
    }
    if (model.apps.isEmpty()) RfButton("Show my apps", onClick = { model.loadApps() }, secondary = true)
    model.apps.forEach { app ->
        val on = app.packageName in model.chosen
        Row(verticalAlignment = Alignment.CenterVertically, modifier = Modifier.fillMaxWidth()) {
            Checkbox(checked = on, onCheckedChange = {
                model.saveChosen(if (it) model.chosen + app.packageName else model.chosen - app.packageName)
            }, colors = CheckboxDefaults.colors(checkedColor = RfColors.gold, checkmarkColor = RfColors.navy))
            Text(app.label, color = RfColors.textPrimary)
        }
    }
    RfButton("Deep Scan", onClick = { model.sheet = Sheet.DeepScan }, secondary = true)
    RfButton("Done", onClick = { model.sheet = null })
}

@Composable
private fun DeepScanManage(model: AppModel) {
    Text(Copy.deepScanTitle, color = RfColors.textPrimary, fontFamily = FontFamily.Serif, fontWeight = FontWeight.Bold, fontSize = 24.sp)
    Text(Copy.deepScanSummary, color = RfColors.textSecondary, modifier = Modifier.padding(vertical = 8.dp))
    Text(Copy.deepScanLimits, color = RfColors.textSecondary)
    RfButton(if (model.deepScanOn) "Deep Scan is running" else "Start Deep Scan", onClick = { model.askSystem("capture") }, enabled = !model.deepScanOn)
}

@Composable
private fun AlertSheet(model: AppModel) {
    val event = model.selectedAlert
    if (event == null) {
        model.sheet = null
        return
    }
    Text(CategoryLabels.label(event.category), color = RfColors.gold, fontWeight = FontWeight.Bold)
    Text(event.summary, color = RfColors.textPrimary, fontFamily = FontFamily.Serif, fontSize = 20.sp, modifier = Modifier.padding(vertical = 8.dp))
    Text(event.timeLabel, color = RfColors.textTertiary)
    Conversation.forCategory(CategoryLabels.starterKey(event.category)).forEach { line ->
        Text(line, color = RfColors.textPrimary, modifier = Modifier.padding(top = 10.dp))
    }
    if (event.alertId != null && !event.discussed) {
        RfButton("Mark as discussed", onClick = { model.discuss(event) })
    }
    RfButton("Close", onClick = { model.sheet = null }, secondary = true)
}

@Composable
private fun ProfileSheet(model: AppModel) {
    var name by remember { mutableStateOf(model.name) }
    var email by remember { mutableStateOf(model.email) }
    Text("Edit Profile", color = RfColors.textPrimary, fontFamily = FontFamily.Serif, fontWeight = FontWeight.Bold, fontSize = 24.sp)
    Field("Name", name) { name = it }
    Field("Email", email, KeyboardType.Email) { email = it }
    RfButton("Save", onClick = { model.updateProfile(name, email) })
}

@Composable
private fun ForgotSheet(model: AppModel) {
    var email by remember { mutableStateOf(model.formEmail) }
    Text(Copy.forgot, color = RfColors.textPrimary, fontFamily = FontFamily.Serif, fontWeight = FontWeight.Bold, fontSize = 24.sp)
    Field("Email", email, KeyboardType.Email) { email = it }
    RfButton("Send reset email", onClick = { model.forgot(email) }, enabled = email.contains("@"))
}

@Composable
private fun TourCard(model: AppModel) {
    val steps = model.tourSteps
    val step = steps.getOrNull(model.tourIndex) ?: return
    Box(Modifier.fillMaxSize().background(RfColors.scrim)) {
        Column(Modifier.align(Alignment.BottomCenter).padding(16.dp).rfCard()) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                Text(step.title, color = RfColors.gold, fontFamily = FontFamily.Serif, fontWeight = FontWeight.Bold, fontSize = 20.sp, modifier = Modifier.weight(1f))
                Text(Copy.skipTour, color = RfColors.textPrimary, modifier = Modifier.clickable { model.skipTour() }.padding(8.dp))
            }
            Text(step.body, color = RfColors.textPrimary, modifier = Modifier.padding(vertical = 8.dp))
            Text("${model.tourIndex + 1} of ${steps.size}", color = RfColors.gold, fontSize = 13.sp)
            Row {
                if (model.tourIndex > 0) Text("Back", color = RfColors.textPrimary, modifier = Modifier.clickable { model.backTour() }.padding(12.dp))
                Spacer(Modifier.weight(1f))
                RfButton(if (model.tourIndex + 1 >= steps.size) "Done" else "Next", onClick = {
                    model.tab = steps.getOrNull(model.tourIndex + 1)?.tab ?: model.tab
                    model.nextTour()
                }, modifier = Modifier.width(140.dp))
            }
        }
    }
}

@Composable
private fun Field(
    label: String,
    value: String,
    keyboard: KeyboardType = KeyboardType.Text,
    password: Boolean = false,
    onChange: (String) -> Unit,
) {
    OutlinedTextField(
        value = value,
        onValueChange = onChange,
        label = { Text(label) },
        modifier = Modifier.fillMaxWidth().padding(vertical = 6.dp),
        singleLine = label != "Details" && label != "Covenant",
        visualTransformation = if (password) PasswordVisualTransformation() else androidx.compose.ui.text.input.VisualTransformation.None,
        keyboardOptions = KeyboardOptions(keyboardType = keyboard),
        colors = rfFieldColors(),
        shape = RoundedCornerShape(RfMetrics.control),
    )
}

@Composable
private fun TextButtonRow(lead: String, action: String, onClick: () -> Unit) {
    Row(Modifier.fillMaxWidth().padding(top = 12.dp), horizontalArrangement = Arrangement.Center) {
        Text(lead, color = RfColors.textSecondary)
        Text(action, color = RfColors.gold, fontWeight = FontWeight.SemiBold, modifier = Modifier.padding(start = 6.dp).clickable(onClick = onClick))
    }
}
