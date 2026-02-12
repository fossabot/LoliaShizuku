//go:build !windows

package services

import "os/exec"

func configureBackgroundProcess(cmd *exec.Cmd) {
	_ = cmd
}
