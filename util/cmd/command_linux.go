//go:build linux
// +build linux

package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"syscall"
	"time"
)

// 杀死指定命令对应的进程
func KillCommand(cmdName string) (exitCode uint32, err error) {
	exitCode = 0
	cmd := exec.Command("killall", cmdName)
	err = cmd.Start()
	if err = cmd.Wait(); err != nil {
		if cmd.ProcessState == nil {
			panic(errors.New("killall command error Please check whether it is installed correctly"))
		}
		waitStatus := cmd.ProcessState.Sys().(syscall.WaitStatus)
		if waitStatus.Exited() {
			exitCode = 128
		}
		return exitCode, err
	} else {
		return 0, nil
	}
}

func MvCommand(cmdName string) (exitCode uint32, err error) {
	exitCode = 0
	st := time.Now()
	stStr := fmt.Sprintf("%d-%d-%d", st.Year(), st.Month(), st.Day())
	cmd := exec.Command("mv", cmdName, cmdName+"-backup-"+stStr)
	err = cmd.Start()
	if err = cmd.Wait(); err != nil {
		if cmd.ProcessState == nil {
			panic(errors.New("killall command error Please check whether it is installed correctly"))
		}
		waitStatus := cmd.ProcessState.Sys().(syscall.WaitStatus)
		if waitStatus.Exited() {
			exitCode = 128
		}
		return exitCode, err
	} else {
		return 0, nil
	}
}
