//go:build windows

package services

import (
	"os/exec"
	"syscall"
)

const windowsCreateNoWindow uint32 = 0x08000000

func configureBackgroundProcess(cmd *exec.Cmd) {
	if cmd == nil {
		return
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: windowsCreateNoWindow,
	}
}
