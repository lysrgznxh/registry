package core

import (
	"agent-network-protocol/registry/core/model"
	"encoding/json"
	"errors"
)

// 对应数据库中的对象
type AgentRegisterReq struct {
	Category    int               `json:"category"`    // 分类
	PublicAddr  string            `json:"public_addr"` // 公钥地址
	Sign        string            `json:"token"`       // 签名
	Author      string            `json:"author"`      // 作者
	Title       string            `json:"title"`       // 中文名称
	Logo        string            `json:"logo"`
	Name        string            `json:"name"`        // 英文名称(唯一)
	Version     string            `json:"version"`     // 版本
	Description string            `json:"description"` // 说明
	WebsiteUrl  string            `json:"website_url"` // 官方网站
	Remotes     []model.Transport `json:"remotes"`
	Status      model.Status      `json:"status"` // 固定 active
	Type        model.AgentType   `json:"type"`   //agent 、 mcp 、group
	Nonce       string            `json:"nonce"`  //随机数,参与签名和验签
}

// 分析remotes参数并返回可以入库的string
func (this *AgentRegisterReq) AnalysisRemotes() (string, error) {
	if len(this.Remotes) == 0 {
		return "", errors.New("远程链接信息不可为空")
	}
	remotes_bytes, err := json.Marshal(this.Remotes)
	if err != nil {
		return "", err
	}
	return string(remotes_bytes), nil
}
