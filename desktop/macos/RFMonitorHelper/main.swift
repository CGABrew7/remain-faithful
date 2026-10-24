import AppKit

let args = CommandLine.arguments
_ = NSApplication.shared

if args.count >= 3 && args[1] == "quit" {
    let bid = args[2]
    for app in NSRunningApplication.runningApplications(withBundleIdentifier: bid) {
        app.terminate()
    }
    exit(0)
}

let front = NSWorkspace.shared.frontmostApplication?.bundleIdentifier ?? ""
print(front)
exit(0)
