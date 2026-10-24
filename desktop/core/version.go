package core

// Version is the desktop marketing version. It tracks iOS 1.0.2 behavior
// and uses its own build so a desktop install is easy to tell apart.
const (
	Version     = "1.0.0"
	BuildNumber = "1"
)

func VersionLabel() string {
	return "Remain Faithful " + Version + " (" + BuildNumber + ")"
}
