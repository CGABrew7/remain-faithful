(function () {
  const field = {};
  const discussed = {};
  const state = {
    screen: "loading",
    tab: "home",
    boot: null,
    home: null,
    settings: null,
    groups: [],
    group: null,
    onStep: 0,
    tour: null,
    tourIndex: 0,
    modal: null,
    toast: "",
    error: "",
    hint: "At least 8 characters.",
    busy: false,
    anim: true,
    deviceCaption: ""
  };

  const shareText = "I would like you as my accountability partner on Remain Faithful. https://www.remainfaithful.com";

  function $(id) { return document.getElementById(id); }

  function esc(value) {
    return String(value == null ? "" : value).replace(/[&<>"']/g, function (c) {
      return { "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;" }[c];
    });
  }

  function errText(err) {
    let msg = "Something went wrong. Please try again.";
    if (typeof err === "string" && err) msg = err;
    else if (err && err.message) msg = String(err.message);
    else if (err) msg = String(err);
    if (msg.indexOf("Error: ") === 0) msg = msg.slice(7);
    return msg;
  }

  function call(name) {
    const args = Array.prototype.slice.call(arguments, 1);
    const app = window.go && window.go.main && window.go.main.App;
    const fn = app && app[name];
    if (!fn) return Promise.reject("Remain Faithful is still starting.");
    return fn.apply(app, args);
  }

  function machine() {
    const platform = state.boot && state.boot.platform;
    if (platform === "mac") return "this Mac";
    if (platform === "windows") return "this PC";
    return "this computer";
  }

  function consentCopy() {
    return "Remain Faithful can see which app is in front and close apps you chose to block. It does not read pages, addresses, or what is on screen. What leaves " + machine() + " is the time, a category, a severity, a short summary, and a heartbeat. The system may ask before a chosen app can be closed.";
  }

  function toast(msg) {
    state.toast = msg;
    const node = $("toast");
    if (node) node.textContent = msg;
    else render();
    clearTimeout(toast.timer);
    toast.timer = setTimeout(function () {
      state.toast = "";
      const live = $("toast");
      if (live) live.remove();
    }, 4200);
  }

  function openURL(url) {
    if (window.runtime && window.runtime.BrowserOpenURL) window.runtime.BrowserOpenURL(url);
  }

  function copyText(text) {
    const done = function () { toast("Copied."); };
    if (window.runtime && window.runtime.ClipboardSetText) {
      Promise.resolve(window.runtime.ClipboardSetText(text)).then(done).catch(function () { toast(text); });
      return;
    }
    toast(text);
  }

  function val(name) { return field[name] == null ? "" : field[name]; }

  function btn(label, act, extra) {
    const cls = extra || "btn btn-primary btn-block";
    return '<button class="' + cls + '" data-act="' + act + '">' + esc(label) + "</button>";
  }

  function fieldInput(name, label, type) {
    return '<label class="stack"><span class="fine">' + esc(label) + '</span><input class="field" data-field="' + name + '" type="' + (type || "text") + '" value="' + esc(val(name)) + '" /></label>';
  }

  function render() {
    const root = $("root");
    if (!root) return;
    const play = state.anim ? " play" : "";
    state.anim = false;
    let body;
    if (state.screen === "loading") {
      body = '<div class="stage"><p class="muted">Remain Faithful</p></div>';
    } else if (state.screen === "main") {
      body = '<div class="shell' + play + '">' + shell() + "</div>";
    } else {
      const inner = state.screen === "register" ? register()
        : state.screen === "login" ? login()
        : state.screen === "forgot" ? forgot()
        : state.screen === "onboard" ? onboard()
        : welcome();
      body = '<div class="stage' + play + '"><div class="narrow">' + inner + "</div></div>";
    }
    const toastHTML = state.toast ? '<div class="toast" id="toast">' + esc(state.toast) + "</div>" : "";
    root.innerHTML = body + modalHTML() + toastHTML;
  }

  function welcome() {
    return '<div class="stack center chunk">' +
      '<div class="halo" aria-hidden="true">RF</div>' +
      "<h1>Remain Faithful</h1>" +
      '<p class="lede">Accountability between trusted friends</p>' +
      '<p class="verse">"But you must remain faithful to the things you have been taught."</p>' +
      '<p class="cite">2 Timothy 3:14</p>' +
      btn("Begin Your Journey", "show-register") +
      '<div class="linkrow"><span class="dim">Already have an account?</span><button class="link" data-act="show-login">Sign In</button></div>' +
      "</div>";
  }

  function register() {
    return '<div class="stack chunk">' +
      "<h1>Create your account</h1>" +
      '<p class="lede">Your journey begins here</p>' +
      fieldInput("name", "Your name") +
      fieldInput("email", "Email address", "email") +
      fieldInput("password", "Password", "password") +
      '<p class="fine" id="password-hint">' + esc(state.hint) + "</p>" +
      '<label class="check"><input type="checkbox" data-check="adult" ' + (field.adult ? "checked" : "") + ' /><span>I am 18 or older, and I am setting this up for myself.</span></label>' +
      '<p class="fine">By continuing you agree to the <button class="link" data-act="open-terms">Terms</button> and <button class="link" data-act="open-privacy">Privacy Policy</button>.</p>' +
      (state.error ? '<p class="warn">' + esc(state.error) + "</p>" : "") +
      btn(state.busy ? "Creating account..." : "Get Started", "do-register") +
      '<button class="link" data-act="show-welcome">Back</button>' +
      "</div>";
  }

  function login() {
    return '<div class="stack chunk">' +
      "<h1>Welcome Back</h1>" +
      '<p class="lede">Sign in to continue your journey</p>' +
      fieldInput("email", "Email address", "email") +
      fieldInput("password", "Password", "password") +
      (state.error ? '<p class="warn">' + esc(state.error) + "</p>" : "") +
      btn(state.busy ? "Signing in..." : "Sign In", "do-login") +
      '<button class="link" data-act="show-forgot">Forgot Password?</button>' +
      '<button class="link" data-act="show-welcome">Back</button>' +
      "</div>";
  }

  function forgot() {
    return '<div class="stack chunk">' +
      "<h1>Forgot Password</h1>" +
      '<p class="lede">We\'ll send a reset link to your email</p>' +
      fieldInput("email", "Email address", "email") +
      (state.error ? '<p class="warn">' + esc(state.error) + "</p>" : "") +
      (state.info ? '<p class="muted">' + esc(state.info) + "</p>" : "") +
      btn("Send reset link", "do-forgot") +
      '<button class="link" data-act="show-login">Back to Sign In</button>' +
      "</div>";
  }

  function onboard() {
    const dots = [0, 1, 2, 3].map(function (i) {
      return '<span class="dot' + (i === state.onStep ? " on" : "") + '"></span>';
    }).join("");
    let body = "";
    if (state.onStep === 0) {
      body = "<h1>Invite your accountability partner</h1>" +
        '<p class="muted">If someone already invited you, accept it in Settings after you sign in. You don\'t have to invite them back.</p>' +
        fieldInput("partnerEmail", "Partner's email address", "email") +
        (state.error ? '<p class="warn">' + esc(state.error) + "</p>" : "") +
        (state.info ? '<p class="muted">' + esc(state.info) + "</p>" : "") +
        btn("Send Invitation", "onboard-invite") +
        '<button class="link" data-act="share-partner">Share invite link</button>' +
        '<p class="fine">Or share a link they can open. They still accept from the email you send. Groups hold up to 12 people.</p>' +
        '<button class="link" data-act="onboard-next">Skip for now</button>';
    } else if (state.onStep === 1) {
      body = "<h1>Allow monitoring</h1><p class=\"muted\">" + esc(consentCopy()) + "</p>" +
        btn("Allow monitoring", "onboard-consent") +
        '<button class="link" data-act="onboard-next">Not now</button>';
    } else if (state.onStep === 2) {
      body = "<h1>Choose apps</h1>" +
        '<p class="muted">Nothing is being watched yet. Choose the apps you want help with. You can take your time.</p>' +
        appList() +
        btn("Save apps", "onboard-apps") +
        '<button class="link" data-act="onboard-next">Continue</button>';
    } else {
      body = "<h1>You are set</h1>" +
        '<p class="muted">Alerts tell your partner after one minute of use, once a day. Blocking closes a chosen app when it opens. You can change both in Settings.</p>' +
        btn("Continue", "onboard-finish");
    }
    return '<div class="stack chunk"><div class="dots">' + dots + "</div>" + body + "</div>";
  }

  function appList() {
    const apps = state.catalog || [];
    if (!apps.length) return '<p class="fine">No apps are listed for this computer yet.</p>';
    const browsers = apps.filter(function (a) { return a.kind === "browser"; });
    const rest = apps.filter(function (a) { return a.kind !== "browser"; });
    return groupApps("Browsers", browsers) + groupApps("Apps", rest);
  }

  function groupApps(title, rows) {
    if (!rows.length) return "";
    return '<div class="card stack"><h3>' + esc(title) + "</h3>" + rows.map(function (app) {
      const on = field["app:" + app.key] ? "checked" : "";
      return '<label class="apppick"><input type="checkbox" data-app="' + esc(app.key) + '" ' + on + ' /><span>' + esc(app.display) + "</span></label>";
    }).join("") + "</div>";
  }

  function shell() {
    const home = state.home || {};
    const pending = (home.pendingPartners || 0) + (home.pendingGroups || 0);
    return '<nav class="rail">' +
      '<div class="brand"><div class="mark" aria-hidden="true">RF</div><div><strong>Remain Faithful</strong><span>Trusted friends</span></div></div>' +
      nav("home", "Home", 0) +
      nav("group", "Group", home.pendingGroups || 0) +
      nav("settings", "Settings", pending && state.tab !== "settings" ? home.pendingPartners || 0 : home.pendingPartners || 0) +
      "</nav><main class=\"main\"><div class=\"column\">" + (state.tab === "group" ? groupTab() : state.tab === "settings" ? settingsTab() : homeTab()) + "</div></main>";
  }

  function nav(tab, label, count) {
    const current = state.tab === tab ? ' aria-current="page"' : "";
    const badge = count ? '<span class="badge num">' + esc(count) + "</span>" : "";
    return '<button class="navbtn" data-act="tab" data-tab="' + tab + '"' + current + "><span>" + label + "</span>" + badge + "</button>";
  }

  function homeTab() {
    const h = state.home;
    if (!h) return '<p class="muted">Loading...</p>';
    const week = (h.week || []).map(function (mark) {
      return '<span class="day ' + esc(mark) + '" title="' + esc(mark) + '"></span>';
    }).join("");
    let action = "";
    if (h.statusAction) action = btn(h.statusAction, "status-action");
    const mine = (h.myEvents || []).slice(0, 8).map(rowHTML).join("") || '<p class="fine">No flags yet.</p>';
    const partners = (h.partnerAlerts || []).slice(0, 12).map(function (row) {
      const done = discussed[row.id];
      return '<div class="item"><div class="spread"><strong>' + esc(row.label) + '</strong><span class="fine num">' + esc(row.when) + "</span></div><p>" + esc(row.summary) + "</p>" +
        '<button class="btn btn-quiet" data-act="discussed" data-id="' + row.id + '"' + (done ? " disabled" : "") + ">" + (done ? "Marked as Discussed" : "Mark as Discussed") + "</button></div>";
    }).join("") || '<p class="fine">No partner alerts yet. Partners see new alerts here.</p>';
    let banners = "";
    if (!h.emailVerified) banners += '<div class="banner spread"><span>Confirm your email to finish account setup.</span><button class="btn btn-secondary" data-act="resend">Resend</button></div>';
    if (h.pendingPartners) banners += '<div class="banner spread"><span class="num">' + h.pendingPartners + ' partner invite' + (h.pendingPartners === 1 ? "" : "s") + '</span><button class="btn btn-secondary" data-act="tab" data-tab="settings">Review</button></div>';
    if (h.pendingGroups) banners += '<div class="banner spread"><span class="num">' + h.pendingGroups + ' group invite' + (h.pendingGroups === 1 ? "" : "s") + '</span><button class="btn btn-secondary" data-act="tab" data-tab="group">Review</button></div>';
    return '<section class="chunk stack"><p class="kicker">Home</p><h1>Hello, ' + esc(h.name || "friend") + "</h1>" + banners + "</section>" +
      '<section class="chunk card stack"><h2>' + esc(h.statusTitle) + "</h2><p class=\"muted\">" + esc(h.statusDetail) + "</p>" + action + "</section>" +
      '<section class="chunk card spread"><div><p class="kicker">Clean streak</p><p class="streak num">' + esc(h.streak || 0) + '</p><p class="fine">Best <span class="num">' + esc(h.best || 0) + '</span></p></div><div class="week" aria-label="This week">' + week + "</div></section>" +
      '<section class="chunk stack"><h2>Recent flags</h2><div class="list">' + mine + "</div></section>" +
      '<section class="chunk stack"><h2>Partner alerts</h2><div class="list">' + partners + "</div></section>" +
      '<section class="chunk">' + btn("I Need Support Right Now", "open-panic", "btn btn-primary btn-block") + "</section>";
  }

  function rowHTML(row) {
    return '<div class="item"><div class="spread"><strong>' + esc(row.label) + '</strong><span class="fine num">' + esc(row.when) + "</span></div><p>" + esc(row.summary) + "</p></div>";
  }

  function groupTab() {
    if (state.group) return groupDetail();
    const cards = (state.groups || []).map(function (g) {
      return '<button class="item" data-act="open-group" data-id="' + g.id + '"><strong>' + esc(g.name) + "</strong></button>";
    }).join("") || '<p class="fine">Start a group here. A name is enough. Groups hold up to 12 people.</p>';
    const pending = (state.pendingGroups || []).map(function (row) {
      return '<div class="item"><strong>' + esc(row.title) + '</strong><p class="fine">' + esc(row.detail) + '</p><div class="split"><button class="btn btn-primary" data-act="accept-group" data-id="' + row.id + '">Accept</button><button class="btn btn-quiet" data-act="decline-group" data-id="' + row.id + '">Decline</button></div></div>';
    }).join("");
    return '<section class="chunk stack"><p class="kicker">Group</p><h1>Your group</h1>' +
      (pending ? '<div class="list">' + pending + "</div>" : "") +
      '<div class="list">' + cards + "</div></section>" +
      '<section class="chunk card stack"><h2>Create a group</h2>' + fieldInput("groupName", "Group name") +
      (state.error ? '<p class="warn">' + esc(state.error) + "</p>" : "") +
      btn("Create group", "create-group") + '<p class="fine">Groups hold up to 12 people.</p></section>';
  }

  function groupDetail() {
    const g = state.group;
    const members = (g.members || []).map(function (m) {
      return '<div class="item"><div class="spread"><strong>' + esc(m.name) + (m.self ? " (you)" : "") + '</strong><span class="fine">' + esc(m.role) + '</span></div><p class="fine num">' + m.streak + " day streak. " + m.flags + " flags</p>" +
        (m.self ? "" : '<button class="btn btn-quiet" data-act="encourage" data-id="' + g.id + '" data-user="' + m.userId + '">Send Encouragement</button>') +
        "</div>";
    }).join("");
    const invite = g.canInvite
      ? fieldInput("inviteEmail", "Email address", "email") + btn("Send invite", "invite-group")
      : '<p class="muted">' + esc(g.fullMessage || "Groups hold up to 12 people.") + "</p>";
    return '<section class="chunk stack"><button class="link" data-act="close-group">All groups</button><h1>' + esc(g.name) + '</h1><p class="fine num">' + esc(g.memberCount) + " of 12</p></section>" +
      '<section class="chunk card stack"><h2>Invite someone</h2><p class="muted">Type an email. They join when they accept, or when they sign up with that email.</p>' + invite +
      (state.error ? '<p class="warn">' + esc(state.error) + "</p>" : "") +
      (state.info ? '<p class="muted">' + esc(state.info) + "</p>" : "") +
      '<button class="btn btn-secondary btn-block" data-act="copy-link" data-id="' + g.id + '">Copy invite link</button></section>' +
      '<section class="chunk stack"><h2>Members</h2><div class="list">' + members + "</div></section>" +
      '<section class="chunk"><button class="btn btn-quiet btn-block" data-act="leave-group" data-id="' + g.id + '">Leave group</button></section>';
  }

  function settingsTab() {
    const s = state.settings;
    if (!s) return '<p class="muted">Loading...</p>';
    const partners = (s.partners || []).map(partnerHTML).join("") || '<p class="fine">Add a partner to start alerts.</p>';
    const pendingP = (s.pendingPartners || []).map(function (row) {
      return '<div class="item"><strong>' + esc(row.title) + '</strong><p class="fine">' + esc(row.detail) + '</p><div class="split"><button class="btn btn-primary" data-act="accept-partner" data-id="' + row.id + '">Accept</button><button class="btn btn-quiet" data-act="decline-partner" data-id="' + row.id + '">Decline</button></div></div>';
    }).join("");
    return '<section class="chunk stack"><p class="kicker">Settings</p><h1>' + esc(s.name || "Settings") + '</h1><p class="fine">' + esc(s.email) + (s.emailVerified ? "" : " (not confirmed)") + "</p></section>" +
      '<section class="chunk card stack"><h2>Profile</h2>' + fieldInput("profileName", "Name") + fieldInput("profileEmail", "Email", "email") +
      btn("Save profile", "save-profile", "btn btn-secondary btn-block") +
      (s.emailVerified ? "" : '<button class="btn btn-quiet btn-block" data-act="resend">Resend confirmation</button>') +
      "</section>" +
      '<section class="chunk card stack"><h2>' + esc(s.statusTitle) + "</h2><p class=\"muted\">" + esc(s.statusDetail) + "</p><p class=\"fine\">" + esc(s.pauseNotice) + "</p>" +
      toggle("Alerts", "Tell your partner after 1 minute of use, once a day", s.alertsOn, "toggle-alerts") +
      toggle("Blocking", "Closes a chosen app when it opens", s.blockingOn, "toggle-blocking") +
      toggle("App Lockout", "Blocking stays on until this is turned off with your partner's PIN. Blocking stays on after lockout is turned off.", s.lockoutOn, "toggle-lockout") +
      '<p class="fine">Partner PIN is ' + (s.pinSet ? "active. Your partner controls gated actions." : "not set. Ask your partner to enable it under Partners.") + "</p>" +
      (s.consent ? '<button class="btn btn-quiet btn-block" data-act="revoke">Turn monitoring access off</button>' : '<button class="btn btn-secondary btn-block" data-act="open-consent">Allow monitoring</button>') +
      "</section>" +
      devicesHTML(s) +
      '<section class="chunk stack"><h2>Choose apps</h2>' + appList() + btn("Save chosen apps", "save-apps", "btn btn-secondary btn-block") + "</section>" +
      '<section class="chunk stack"><h2>Your partners</h2>' + fieldInput("partnerEmail", "Partner's email address", "email") +
      btn("Send Invitation", "invite-partner", "btn btn-secondary btn-block") +
      '<button class="link" data-act="share-partner">Share invite link</button>' +
      (state.info ? '<p class="muted">' + esc(state.info) + "</p>" : "") +
      (pendingP ? '<div class="list">' + pendingP + "</div>" : "") +
      '<div class="list">' + partners + "</div></section>" +
      '<section class="chunk stack">' +
      btn("Send Ideas or Report a Problem", "open-feedback", "btn btn-secondary btn-block") +
      btn("How It Works", "replay-tour", "btn btn-quiet btn-block") +
      btn("Sign Out", "ask-signout", "btn btn-quiet btn-block") +
      btn("Delete Account", "ask-delete", "btn btn-danger btn-block") +
      '<p class="fine">' + esc(s.deviceCaption || s.version) + "</p></section>";
  }

  function devicesHTML(s) {
    if (s.devicesUnavailable) {
      return '<section class="chunk stack"><h2>My devices</h2><p class="muted">Device list is not available yet.</p></section>';
    }
    const rows = (s.devices || []).map(function (d) {
      const mark = d.monitoring ? '<p class="ok">Monitoring on</p>' : '<p class="fine">Monitoring off</p>';
      const here = d.current ? '<span class="fine">This computer</span>' : "<span></span>";
      return '<div class="item"><div class="spread"><strong>' + esc(d.label) + "</strong>" + here + "</div>" +
        mark +
        '<p class="fine num">' + esc(d.lastCheckIn || "No check-in yet") + "</p>" +
        '<button class="btn btn-quiet" data-act="remove-device" data-id="' + esc(d.deviceId) + '" data-label="' + esc(d.label) + '">Remove</button></div>';
    }).join("");
    const body = rows || '<p class="fine">No devices are registered yet.</p>';
    return '<section class="chunk stack"><h2>My devices</h2><p class="muted">Each phone or computer on this account is listed on its own. Removing one tells your partners, the same as turning monitoring off.</p><div class="list">' + body + "</div></section>";
  }

  function toggle(title, detail, on, act) {
    return '<button class="switch" role="switch" aria-checked="' + (on ? "true" : "false") + '" data-act="' + act + '"><span><strong>' + title + '</strong><br /><span class="fine">' + esc(detail) + '</span></span><span class="knob" aria-hidden="true"></span></button>';
  }

  function partnerHTML(p) {
    let actions = "";
    if (p.status === "accepted" && p.canStar && !p.primary) actions += '<button class="btn btn-quiet" data-act="primary" data-id="' + p.id + '">Make primary</button>';
    if (p.canSetPin) {
      actions += '<button class="btn btn-quiet" data-act="set-pin" data-id="' + p.id + '" data-name="' + esc(p.name) + '" data-replacing="' + (p.pinSet ? "1" : "0") + '">' + (p.pinSet ? "Change PIN" : "Set PIN") + "</button>";
      if (p.pinSet) actions += '<button class="btn btn-quiet" data-act="remove-pin" data-id="' + p.id + '">Remove PIN</button>';
    }
    actions += '<button class="btn btn-quiet" data-act="remove-partner" data-id="' + p.id + '">Remove</button>';
    return '<div class="item"><div class="spread"><strong>' + esc(p.name || "Partner") + '</strong><span class="fine">' + esc(p.status) + (p.primary ? " · primary" : "") + "</span></div><p class=\"fine\">" + esc(p.email) + '</p><div class="split">' + actions + "</div></div>";
  }

  function modalHTML() {
    const m = state.modal;
    if (state.tour) return tourHTML();
    if (!m) return "";
    let inner = "";
    if (m.type === "pin") {
      inner = "<h2>Partner PIN</h2><p class=\"muted\">Enter the Partner PIN your partner set.</p>" +
        '<input class="field num" data-field="pin" inputmode="numeric" maxlength="4" value="' + esc(val("pin")) + '" />' +
        (m.error ? '<p class="warn">' + esc(m.error) + "</p>" : "") +
        '<div class="split">' + btn("Continue", "submit-pin", "btn btn-primary") + '<button class="btn btn-quiet" data-act="close-modal">Cancel</button></div>';
    } else if (m.type === "consent") {
      inner = "<h2>Allow monitoring</h2><p class=\"muted\">" + esc(consentCopy()) + "</p>" + btn("Allow monitoring", "grant-consent") + '<button class="btn btn-quiet btn-block" data-act="close-modal">Not now</button>';
    } else if (m.type === "shield") {
      inner = '<div class="halo" aria-hidden="true">RF</div><h2>Blocking is on. This chosen app was closed.</h2><p class="lede center">' + esc(m.name || "A chosen app") + "</p>" + btn("Close", "close-modal");
    } else if (m.type === "panic") {
      inner = '<div class="halo" aria-hidden="true">RF</div><h2 class="center">You Are Not Alone</h2>' +
        '<p class="muted center">Your accountability partners want to hear from you right now. Reaching out is a sign of strength, not weakness.</p>' +
        '<div class="card stack"><p class="kicker">Hold on to this</p><p class="verse">"No temptation has overtaken you except what is common to mankind. And God is faithful; he will not let you be tempted beyond what you can bear. But when you are tempted, he will also provide a way out so that you can endure it."</p><p class="cite">1 Corinthians 10:13</p></div>' +
        '<p class="kicker">Your accountability partner</p><p>' + esc(m.partner || "Add an accountability partner to use this feature") + "</p>" +
        (m.sent ? btn("Alert Sent", "noop") : btn(m.busy ? "Sending..." : "Alert Partners and Group", "send-panic")) +
        (m.error ? '<p class="warn">' + esc(m.error) + "</p>" : "") +
        btn("I'm okay, close this", "close-modal", "btn btn-secondary btn-block");
    } else if (m.type === "feedback") {
      inner = "<h2>Send Ideas or Report a Problem</h2><p class=\"fine\">We'll reply to your account email.</p>" +
        '<div class="split"><button class="btn ' + (field.kind === "Idea" ? "btn-primary" : "btn-quiet") + '" data-act="kind" data-kind="Idea">Idea</button><button class="btn ' + (field.kind === "Improvement" ? "btn-primary" : "btn-quiet") + '" data-act="kind" data-kind="Improvement">Improvement</button><button class="btn ' + (field.kind === "Problem" ? "btn-primary" : "btn-quiet") + '" data-act="kind" data-kind="Problem">Problem</button></div>' +
        '<label class="stack"><span class="fine">In a few words</span><input class="field" data-field="summary" maxlength="120" value="' + esc(val("summary")) + '" /></label><p class="fine num" id="summary-count">' + val("summary").length + ' of 120</p>' +
        '<label class="stack"><span class="fine">Details</span><textarea class="area" data-field="details" maxlength="4000">' + esc(val("details")) + '</textarea></label><p class="fine num" id="details-count">' + val("details").length + ' of 4000</p>' +
        '<label class="check"><input type="checkbox" data-check="device" ' + (field.device ? "checked" : "") + ' /><span>Include app version and computer model<br /><span class="fine">' + esc(state.deviceCaption) + "</span></span></label>" +
        (m.error ? '<p class="warn">' + esc(m.error) + "</p>" : "") +
        (m.failed ? '<button class="link" data-act="mailto-support">Email support@remainfaithful.com</button>' : "") +
        btn("Send", "send-feedback") + '<button class="btn btn-quiet btn-block" data-act="close-modal">Close</button>';
    } else if (m.type === "delete") {
      inner = "<h2>Delete Account</h2><p class=\"muted\">This will permanently delete your account, all your data, and notify your accountability partners that you've left. This cannot be undone.</p>" +
        fieldInput("confirm", "Type DELETE to confirm") +
        (m.error ? '<p class="warn">' + esc(m.error) + "</p>" : "") +
        btn("Delete Account", "do-delete", "btn btn-danger btn-block") +
        '<button class="btn btn-quiet btn-block" data-act="close-modal">Cancel</button>';
    } else if (m.type === "signout") {
      inner = "<h2>Sign Out</h2><p class=\"muted\">Signing out stops monitoring on this computer. Your partners will be told. Blocking on this computer turns off until you sign back in.</p>" +
        btn("Sign Out", "do-signout", "btn btn-danger btn-block") +
        '<button class="btn btn-quiet btn-block" data-act="close-modal">Cancel</button>';
    } else if (m.type === "remove-device") {
      inner = "<h2>Remove this " + esc(m.label || "device") + "?</h2><p class=\"muted\">Partners are told that activity monitoring was turned off on that device.</p>" +
        btn("Remove device", "do-remove-device", "btn btn-danger btn-block") +
        '<button class="btn btn-quiet btn-block" data-act="close-modal">Cancel</button>';
    } else if (m.type === "setpin") {
      inner = "<h2>Partner PIN</h2><p class=\"muted\">Set a 4-digit PIN for " + esc(m.name) + ". Share it verbally. It will be required to change protection settings.</p>" +
        '<input class="field num" data-field="newpin" inputmode="numeric" maxlength="4" value="' + esc(val("newpin")) + '" />' +
        (m.error ? '<p class="warn">' + esc(m.error) + "</p>" : "") +
        btn("Save PIN", "save-pin") + '<button class="btn btn-quiet btn-block" data-act="close-modal">Cancel</button>';
    }
    return '<div class="scrim" data-act="scrim"><div class="sheet" id="sheet"><div class="sheet-inner">' + inner + "</div></div></div>";
  }

  function tourHTML() {
    const step = state.tour[state.tourIndex];
    if (!step) return "";
    const last = state.tourIndex === state.tour.length - 1;
    return '<div class="scrim"><div class="sheet"><div class="sheet-inner"><p class="kicker num">How It Works ' + (state.tourIndex + 1) + " of " + state.tour.length + "</p><h2>" + esc(step.title) + "</h2><p class=\"muted\">" + esc(step.body) + "</p>" +
      '<div class="split"><button class="btn btn-quiet" data-act="tour-skip">Skip</button><button class="btn btn-primary" data-act="tour-next">' + (last ? "Done" : "Next") + "</button></div></div></div></div>";
  }

  function show(screen) {
    state.screen = screen;
    state.modal = null;
    state.tour = null;
    state.error = "";
    state.info = "";
    state.anim = true;
    render();
  }

  async function boot() {
    try {
      state.boot = await call("Bootstrap");
    } catch (err) {
      state.screen = "welcome";
      state.error = errText(err);
      render();
      return;
    }
    if (!state.boot.signedIn) {
      show("welcome");
      return;
    }
    if (!state.boot.onboardingComplete) {
      await loadCatalog();
      show("onboard");
      return;
    }
    state.screen = "main";
    await enterApp(!!state.boot.showTour);
  }

  async function enterApp(withTour) {
    state.screen = "main";
    state.anim = true;
    await loadHome();
    if (state.tab === "settings") await loadSettings();
    if (state.tab === "group") await loadGroups();
    render();
    if (withTour) await startTour(false);
  }

  async function loadHome() {
    state.home = await call("Home");
  }

  async function loadSettings() {
    state.settings = await call("Settings");
    state.deviceCaption = state.settings.deviceCaption || "";
    field.profileName = state.settings.name || "";
    field.profileEmail = state.settings.email || "";
    (state.settings.apps || []).forEach(function (app) {
      field["app:" + app.key] = !!app.selected;
    });
    state.catalog = state.settings.apps || [];
  }

  async function loadCatalog() {
    state.catalog = await call("Catalog");
    (state.catalog || []).forEach(function (app) {
      if (field["app:" + app.key] == null) field["app:" + app.key] = !!app.selected;
    });
  }

  async function loadGroups() {
    state.groups = await call("ListGroups");
    state.pendingGroups = await call("PendingGroupInvites");
  }

  function selectedKeys() {
    return (state.catalog || []).filter(function (app) { return field["app:" + app.key]; }).map(function (app) { return app.key; });
  }

  async function gated(fn) {
    try {
      await fn(false);
      return true;
    } catch (err) {
      const msg = errText(err);
      if (msg.indexOf("Partner PIN") >= 0) {
        const ok = await pinGate();
        if (!ok) return false;
        await fn(true);
        return true;
      }
      if (msg.indexOf("Allow monitoring first") >= 0) {
        state.modal = { type: "consent" };
        render();
        return false;
      }
      throw err;
    }
  }

  function pinGate() {
    field.pin = "";
    return new Promise(function (resolve) {
      state.pinResolve = resolve;
      state.modal = { type: "pin", error: "" };
      render();
    });
  }

  function closeModalSoon() {
    const sheet = $("sheet");
    if (!sheet) {
      state.modal = null;
      render();
      return;
    }
    sheet.classList.add("leaving");
    setTimeout(function () {
      state.modal = null;
      render();
    }, 150);
  }

  async function startTour(force) {
    const inGroup = state.home ? !!state.home.inGroup : (state.groups || []).length > 0;
    state.tour = await call("TourSteps", inGroup);
    state.tourIndex = 0;
    state.tourForce = force;
    if (state.tour && state.tour[0]) state.tab = state.tour[0].tab === "group" ? "group" : state.tour[0].tab;
    render();
  }

  async function finishTour() {
    state.tour = null;
    if (!state.tourForce) {
      try { await call("FinishTour"); } catch (err) { toast(errText(err)); }
    }
    state.tourForce = false;
    render();
  }

  async function runAct(act, el) {
    if (state.busy) return;
    try {
      await handle(act, el);
    } catch (err) {
      const msg = errText(err);
      if (state.modal) {
        state.modal.error = msg;
        render();
      } else {
        state.error = msg;
        toast(msg);
        render();
      }
    }
  }

  async function onClick(event) {
    const el = event.target.closest("[data-act]");
    if (!el) return;
    if (el.dataset.act === "scrim" && event.target !== el) return;
    await runAct(el.dataset.act, el);
  }

  function enterAction() {
    if (state.tour) return "";
    if (state.modal) {
      if (state.modal.type === "pin") return "submit-pin";
      if (state.modal.type === "setpin") return "save-pin";
      if (state.modal.type === "delete") return "do-delete";
      if (state.modal.type === "feedback") return "send-feedback";
      return "";
    }
    if (state.screen === "login") return "do-login";
    if (state.screen === "register") return "do-register";
    if (state.screen === "forgot") return "do-forgot";
    if (state.screen === "onboard" && state.onStep === 0) return "onboard-invite";
    return "";
  }

  async function handle(act, el) {
    if (act === "show-welcome") return show("welcome");
    if (act === "show-register") return show("register");
    if (act === "show-login") return show("login");
    if (act === "show-forgot") { state.info = ""; return show("forgot"); }
    if (act === "open-terms") return openURL("https://remainfaithful.com/terms");
    if (act === "open-privacy") return openURL("https://remainfaithful.com/privacy");
    if (act === "mailto-support") return openURL("mailto:support@remainfaithful.com");
    if (act === "do-register") {
      state.busy = true; state.error = ""; render();
      try {
        await call("Register", val("name"), val("email"), val("password"), !!field.adult);
      } finally {
        state.busy = false;
      }
      field.password = "";
      await boot();
      return;
    }
    if (act === "do-login") {
      state.busy = true; state.error = ""; render();
      try {
        await call("Login", val("email"), val("password"));
      } finally {
        state.busy = false;
      }
      field.password = "";
      await boot();
      return;
    }
    if (act === "do-forgot") {
      state.info = await call("ForgotPassword", val("email"));
      state.error = "";
      render();
      return;
    }
    if (act === "onboard-next") {
      state.onStep += 1;
      state.error = "";
      state.info = "";
      state.anim = true;
      if (state.onStep === 2) await loadCatalog();
      render();
      return;
    }
    if (act === "onboard-invite") {
      state.info = await call("InvitePartner", val("partnerEmail"));
      state.error = "";
      render();
      return;
    }
    if (act === "share-partner") { copyText(shareText); return; }
    if (act === "onboard-consent") { await call("GrantConsent"); state.onStep = 2; await loadCatalog(); state.anim = true; render(); return; }
    if (act === "onboard-apps") { await call("SetSelection", selectedKeys(), false); state.onStep = 3; state.anim = true; render(); return; }
    if (act === "onboard-finish") { await call("CompleteOnboarding"); await enterApp(true); return; }
    if (act === "tab") {
      state.tab = el.dataset.tab;
      state.group = null;
      state.error = "";
      state.info = "";
      state.anim = true;
      if (state.tab === "home") await loadHome();
      if (state.tab === "group") await loadGroups();
      if (state.tab === "settings") await loadSettings();
      render();
      return;
    }
    if (act === "status-action") return statusAction();
    if (act === "discussed") {
      try {
        await call("MarkDiscussed", Number(el.dataset.id));
        discussed[el.dataset.id] = true;
      } catch (err) {
        toast("Couldn't save. Try again.");
        return;
      }
      render();
      return;
    }
    if (act === "open-panic") {
      let partner = "";
      try { partner = await call("PanicPartner"); } catch (err) { partner = ""; }
      state.modal = { type: "panic", partner: partner, sent: false, error: "" };
      render();
      return;
    }
    if (act === "send-panic") {
      state.modal.busy = true; state.modal.error = ""; render();
      try {
        await call("Panic");
        state.modal.busy = false;
        state.modal.sent = true;
      } catch (err) {
        state.modal.busy = false;
        state.modal.error = errText(err);
      }
      render();
      return;
    }
    if (act === "resend") { toast(await call("ResendVerification")); return; }
    if (act === "create-group") {
      const card = await call("CreateGroup", val("groupName"));
      field.groupName = "";
      state.group = await call("GetGroup", Number(card.id));
      state.anim = true;
      render();
      return;
    }
    if (act === "open-group") { state.group = await call("GetGroup", Number(el.dataset.id)); state.anim = true; render(); return; }
    if (act === "close-group") { state.group = null; await loadGroups(); render(); return; }
    if (act === "invite-group") {
      state.info = await call("InviteToGroup", state.group.id, val("inviteEmail"));
      state.error = "";
      field.inviteEmail = "";
      render();
      return;
    }
    if (act === "copy-link") { copyText(await call("GroupLink", Number(el.dataset.id))); return; }
    if (act === "leave-group") { await call("LeaveGroup", Number(el.dataset.id)); state.group = null; await loadGroups(); render(); return; }
    if (act === "encourage") {
      await call("Encourage", Number(el.dataset.id), Number(el.dataset.user));
      toast("Encouragement Sent");
      return;
    }
    if (act === "accept-group") { await call("AcceptGroupInvite", Number(el.dataset.id)); await loadGroups(); await loadHome(); render(); return; }
    if (act === "decline-group") { await call("DeclineGroupInvite", Number(el.dataset.id)); await loadGroups(); render(); return; }
    if (act === "save-profile") { toast(await call("UpdateProfile", val("profileName"), val("profileEmail"))); await loadSettings(); render(); return; }
    if (act === "toggle-alerts") { await gated(function (ok) { return call("SetAlerts", !state.settings.alertsOn, ok); }); await loadSettings(); await loadHome(); render(); return; }
    if (act === "toggle-blocking") { await gated(function (ok) { return call("SetBlocking", !state.settings.blockingOn, ok); }); await loadSettings(); await loadHome(); render(); return; }
    if (act === "toggle-lockout") { await gated(function (ok) { return call("SetLockout", !state.settings.lockoutOn, ok); }); await loadSettings(); await loadHome(); render(); return; }
    if (act === "save-apps") { await gated(function (ok) { return call("SetSelection", selectedKeys(), ok); }); await loadSettings(); toast("Saved."); render(); return; }
    if (act === "revoke") { await gated(function (ok) { return call("RevokeConsent", ok); }); await loadSettings(); render(); return; }
    if (act === "open-consent") { state.modal = { type: "consent" }; render(); return; }
    if (act === "grant-consent") { await call("GrantConsent"); state.modal = null; if (state.screen === "main") await loadSettings(); render(); return; }
    if (act === "invite-partner") { state.info = await call("InvitePartner", val("partnerEmail")); await loadSettings(); render(); return; }
    if (act === "accept-partner") { await call("AcceptPartnerInvite", Number(el.dataset.id)); await loadSettings(); await loadHome(); render(); return; }
    if (act === "decline-partner") { await call("DeclinePartnerInvite", Number(el.dataset.id)); await loadSettings(); render(); return; }
    if (act === "primary") { await call("SetPrimary", Number(el.dataset.id)); await loadSettings(); render(); return; }
    if (act === "remove-partner") { await call("RemovePartner", Number(el.dataset.id)); await loadSettings(); render(); return; }
    if (act === "remove-device") {
      state.modal = { type: "remove-device", id: el.dataset.id, label: el.dataset.label || "device" };
      render();
      return;
    }
    if (act === "do-remove-device") {
      const id = state.modal && state.modal.id;
      const left = await gated(function (ok) { return call("RemoveDevice", id, ok); });
      if (!left) return;
      state.modal = null;
      await loadSettings();
      await loadHome();
      toast("Removed. Partners were told.");
      render();
      return;
    }
    if (act === "set-pin") {
      field.newpin = "";
      state.modal = { type: "setpin", id: Number(el.dataset.id), name: el.dataset.name, replacing: el.dataset.replacing === "1" };
      render();
      return;
    }
    if (act === "save-pin") {
      await call("SetPIN", state.modal.id, val("newpin"), !!state.modal.replacing);
      state.modal = null;
      toast("Saved.");
      await loadSettings();
      render();
      return;
    }
    if (act === "remove-pin") { await call("RemovePIN", Number(el.dataset.id)); await loadSettings(); render(); return; }
    if (act === "open-feedback") {
      field.kind = field.kind || "Idea";
      try { state.deviceCaption = await call("DeviceCaption"); } catch (err) { state.deviceCaption = ""; }
      state.modal = { type: "feedback" };
      render();
      return;
    }
    if (act === "kind") { field.kind = el.dataset.kind; render(); return; }
    if (act === "send-feedback") {
      try {
        await call("SendFeedback", field.kind || "Idea", val("summary"), val("details"), !!field.device);
      } catch (err) {
        const msg = errText(err);
        state.modal.failed = msg.indexOf("required") < 0 && msg.indexOf("type must") < 0;
        state.modal.error = state.modal.failed ? "Couldn't send. Try again or email support@remainfaithful.com" : msg;
        render();
        return;
      }
      field.summary = "";
      field.details = "";
      state.modal = null;
      toast("Sent.");
      render();
      return;
    }
    if (act === "replay-tour") { await startTour(true); return; }
    if (act === "ask-signout") { state.modal = { type: "signout" }; render(); return; }
    if (act === "do-signout") {
      const left = await gated(function (ok) { return call("SignOut", ok); });
      if (!left) return;
      state.home = null;
      state.settings = null;
      show("welcome");
      return;
    }
    if (act === "ask-delete") { state.modal = { type: "delete" }; render(); return; }
    if (act === "do-delete") {
      await call("DeleteAccount", val("confirm").trim());
      state.modal = null;
      show("welcome");
      return;
    }
    if (act === "submit-pin") {
      const view = await call("VerifyPIN", val("pin"));
      if (!view.success) {
        state.modal.error = view.message || "Incorrect PIN. Your partner is the one who knows it.";
        render();
        return;
      }
      const resolve = state.pinResolve;
      state.pinResolve = null;
      state.modal = null;
      if (resolve) resolve(true);
      return;
    }
    if (act === "close-modal" || act === "scrim") {
      if (state.pinResolve) {
        const resolve = state.pinResolve;
        state.pinResolve = null;
        resolve(false);
      }
      closeModalSoon();
      return;
    }
    if (act === "tour-skip") return finishTour();
    if (act === "tour-next") {
      if (state.tourIndex >= state.tour.length - 1) return finishTour();
      state.tourIndex += 1;
      const step = state.tour[state.tourIndex];
      if (step && step.tab) {
        state.tab = step.tab;
        if (step.tab === "group") await loadGroups();
        if (step.tab === "settings") await loadSettings();
      }
      render();
    }
  }

  async function statusAction() {
    const action = state.home && state.home.statusAction;
    if (action === "Allow monitoring") {
      state.modal = { type: "consent" };
      render();
      return;
    }
    if (action === "Choose apps") {
      state.tab = "settings";
      await loadSettings();
      state.anim = true;
      render();
      return;
    }
    if (action === "Turn monitoring on") {
      try {
        await gated(function (ok) { return call("SetAlerts", true, ok); });
      } catch (err) {
        toast(errText(err));
      }
      await loadHome();
      render();
    }
  }

  document.addEventListener("input", function (event) {
    const el = event.target;
    if (el.dataset.field) {
      field[el.dataset.field] = el.value;
      if (el.dataset.field === "password") {
        call("PasswordHint", el.value).then(function (hint) {
          state.hint = hint;
          const node = $("password-hint");
          if (node) node.textContent = hint;
        }).catch(function () {});
      }
      if (el.dataset.field === "summary" || el.dataset.field === "details") {
        const node = $(el.dataset.field + "-count");
        const max = el.dataset.field === "summary" ? 120 : 4000;
        if (node) node.textContent = el.value.length + " of " + max;
      }
    }
  });

  document.addEventListener("change", function (event) {
    const el = event.target;
    if (el.dataset.check) field[el.dataset.check] = el.checked;
    if (el.dataset.app) field["app:" + el.dataset.app] = el.checked;
  });

  document.addEventListener("click", function (event) {
    onClick(event);
  });

  document.addEventListener("keydown", function (event) {
    if (event.key === "Escape" && (state.modal || state.tour)) {
      if (state.tour) finishTour();
      else runAct("close-modal", event.target);
      return;
    }
    if (event.key === "Enter" && (!event.target || event.target.tagName !== "TEXTAREA")) {
      const act = enterAction();
      if (!act) return;
      event.preventDefault();
      runAct(act, event.target);
    }
  });

  function listen() {
    if (!window.runtime || !window.runtime.EventsOn) return;
    window.runtime.EventsOn("shield", function (name) {
      state.modal = { type: "shield", name: name };
      state.anim = false;
      render();
    });
    window.runtime.EventsOn("alerts", function (count) {
      if (state.home) state.home.unread = count;
      const badge = document.querySelector('[data-tab="home"] .badge');
      if (count && badge) badge.textContent = String(count);
    });
    window.runtime.EventsOn("device-removed", function () {
      toast("This computer was removed. Partners were told that monitoring stopped.");
      if (state.screen !== "main") return;
      Promise.resolve().then(function () { return loadSettings(); }).then(function () { return loadHome(); }).then(render).catch(function () { render(); });
    });
    window.runtime.EventsOn("signed-out", function () {
      state.home = null;
      state.settings = null;
      show("welcome");
      toast("Not signed in. Please log in again.");
    });
  }

  function whenReady(fn) {
    const ready = window.go && window.go.main && window.go.main.App;
    if (ready) { fn(); return; }
    let tries = 0;
    const timer = setInterval(function () {
      tries += 1;
      if (window.go && window.go.main && window.go.main.App) {
        clearInterval(timer);
        fn();
      } else if (tries > 50) {
        clearInterval(timer);
        const root = $("root");
        if (root) root.textContent = "Remain Faithful could not start.";
      }
    }, 100);
  }

  whenReady(function () {
    listen();
    boot().catch(function (err) {
      state.screen = "welcome";
      state.error = errText(err);
      render();
    });
  });
})();
