package cmd

import (
	"agent-network-protocol/registry/core"
	"os/exec"
	"strings"
)

// 运行命令
func RunCommand(cmdPath string, args ...string) (cmd *exec.Cmd, err error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmNodeUtil)
	fullCmd := cmdPath + " " + strings.Join(args, " ")
	log.WithField("path", fullCmd).Info("run command")
	cmd = exec.Command(cmdPath, args...)
	return
}
