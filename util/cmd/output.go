package cmd

import (
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"time"
)

// 捕获明令行输出并保存日志
type CmdOutput struct {
	logWriter *rotatelogs.RotateLogs
	lastMsg   []byte
}

func (this *CmdOutput) Write(p []byte) (n int, err error) {
	this.lastMsg = p
	return this.logWriter.Write(p)
}

func (this *CmdOutput) GetLastMsg() []byte {
	return this.lastMsg
}

func NewCmdOutput(baseLogPath string, life, split time.Duration) (*CmdOutput, error) {
	writer, err := rotatelogs.New(
		baseLogPath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPath), // Generate a soft chain and point to the latest log file
		rotatelogs.WithMaxAge(life),
		rotatelogs.WithRotationTime(split),
	)
	if err != nil {
		return nil, err
	}
	cmdOutput := CmdOutput{
		logWriter: writer,
	}
	return &cmdOutput, nil
}
