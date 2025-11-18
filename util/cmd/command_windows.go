//go:build windows
// +build windows

package cmd

import (
	"errors"
	"golang.org/x/text/encoding/simplifiedchinese"
	"os/exec"
	"syscall"
)

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

// 杀死指定命令对应的进程
func KillCommand(cmdName string) (exitCode uint32, err error) {
	cmd := exec.Command("taskkill.exe", "/f", "/im", cmdName)
	output, err := cmd.CombinedOutput()
	cmdRe := ConvertByte2String(output, "GB18030")
	if err != nil {
		if cmd.ProcessState == nil {
			panic(errors.New("taskkill.exe error"))
		}
		exitCode = cmd.ProcessState.Sys().(syscall.WaitStatus).ExitCode
		return exitCode, errors.New(cmdRe)
	} else {
		return 0, nil
	}
}

func MvCommand(cmdName string) (exitCode uint32, err error) {
	return 0, nil
}

type Charset string

func ConvertByte2String(byte []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}
