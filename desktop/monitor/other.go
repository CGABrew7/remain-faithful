//go:build !darwin && !windows

package monitor

import (
	"errors"

	"remain-faithful/desktop/core"
)

func Platform() string { return "other" }

func Ready() bool { return false }

func Sample() (Front, error) { return Front{}, nil }

func Quit(key string) error {
	if !core.AllowQuit(Platform(), key) {
		return errors.New("that app is not chosen for blocking")
	}
	return errors.New("blocking is not available on this system")
}
