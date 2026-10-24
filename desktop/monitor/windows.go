//go:build windows

package monitor

import (
	"errors"
	"os"
	"strings"
	"unsafe"

	"remain-faithful/desktop/core"

	"golang.org/x/sys/windows"
)

const processTerminate = 0x0001

var (
	user32                   = windows.NewLazySystemDLL("user32.dll")
	getForegroundWindow      = user32.NewProc("GetForegroundWindow")
	getWindowThreadProcessID = user32.NewProc("GetWindowThreadProcessId")
)

func Platform() string { return "windows" }

func Ready() bool { return true }

func Sample() (Front, error) {
	hwnd, _, _ := getForegroundWindow.Call()
	if hwnd == 0 {
		return Front{}, nil
	}
	var pid uint32
	_, _, _ = getWindowThreadProcessID.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
	if pid == 0 {
		return Front{}, nil
	}
	name, err := exeBase(pid)
	if err != nil || name == "" || core.IsSelf(Platform(), name) {
		return Front{}, err
	}
	return Front{Key: name}, nil
}

func Quit(key string) error {
	if !core.AllowQuit(Platform(), key) {
		return errors.New("that app is not chosen for blocking")
	}
	return quitChosen(strings.ToLower(key))
}

func exeBase(pid uint32) (string, error) {
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return "", err
	}
	defer windows.CloseHandle(snap)
	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))
	if err := windows.Process32First(snap, &entry); err != nil {
		return "", err
	}
	for {
		if entry.ProcessID == pid {
			return strings.ToLower(windows.UTF16ToString(entry.ExeFile[:])), nil
		}
		if err := windows.Process32Next(snap, &entry); err != nil {
			return "", nil
		}
	}
}

func quitChosen(key string) error {
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return err
	}
	defer windows.CloseHandle(snap)
	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))
	if err := windows.Process32First(snap, &entry); err != nil {
		return err
	}
	me := uint32(os.Getpid())
	for {
		name := strings.ToLower(windows.UTF16ToString(entry.ExeFile[:]))
		if name == key && entry.ProcessID != 0 && entry.ProcessID != me {
			if h, err := windows.OpenProcess(processTerminate, false, entry.ProcessID); err == nil {
				_ = windows.TerminateProcess(h, 1)
				_ = windows.CloseHandle(h)
			}
		}
		if err := windows.Process32Next(snap, &entry); err != nil {
			break
		}
	}
	return nil
}
