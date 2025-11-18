package agent

import (
	"agent-network-protocol/registry/core"
	"fmt"
	"testing"
)

func TestGetAllAgent(t *testing.T) {
	//加载配置文件
	err := core.LoadConfig("E:\\gopath\\src\\github.com\\agent-network-protocol\\AgentCenter\\cmd\\registry\\config.toml")
	if err != nil {
		panic(err.Error())
	}
	err = Init(core.RunCtrl)
	if err != nil {
		t.Fatal(err)
		return
	}

	fmt.Println("LocalServerUrl:", core.BaseConf.LocalServerUrl)
	defer agentDb.Close()
	list, err := SearchAgent(1, 10, "mcp", "", "", []string{}, "")
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	for _, v := range list {
		fmt.Println("name:", v.Name, " version:", v.Version)
	}
}

func TestGetMcpById(t *testing.T) {
	Init(core.RunCtrl)
	fmt.Println(GetAgentById(99))
}
