package api

import (
	"agent-network-protocol/registry/core"
	"agent-network-protocol/registry/db/agent"
	"agent-network-protocol/registry/util/crypto"
	"agent-network-protocol/registry/util/serializer"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"regexp"
	"time"
)

// agent 注册
func agentRegister(w http.ResponseWriter, r *http.Request) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmServApi)

	req := &core.AgentRegisterReq{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		serializer.ErrorRes(w, err)
		return
	}
	if req.Sign == "" {
		log.Error("sign参数为空")
		serializer.StringErrorRes(w, errors.New("sign参数为空"))
		return
	}
	if req.Nonce == "" {
		log.Error("nonce参数为空")
		serializer.StringErrorRes(w, errors.New("nonce参数为空"))
		return
	} else if len(req.Nonce) < 32 {
		log.Error("nonce参数长度不得小于32位")
		serializer.StringErrorRes(w, errors.New("nonce参数长度不得小于32位"))
		return
	}

	regex := regexp.MustCompile(`^[a-zA-Z0-9_\-\.]{3,50}$`)
	if !regex.MatchString(req.Name) {
		serializer.StringErrorRes(w, errors.New("name字段只能包含数字、字母、下划线、减号 长度为3-50个字节"))
		return
	}
	//版本无效
	if !core.IsValidVersion(req.Version) {
		log.WithError(err).Error("IsValidVersion")
		serializer.StringErrorRes(w, errors.New("传入的版本格式无效,参考格式: 1.0.1"))
		return
	}
	remotes := ""
	if len(req.Remotes) > 0 {
		remotes, err = req.AnalysisRemotes()
		if err != nil {
			log.WithError(err).Error("AnalysisRemotes")
			serializer.StringErrorRes(w, err)
			return
		}
	}
	//检查是否已经存在,如果已经存在则更新,否则插入
	newOb := agent.Agent{
		Category:    req.Category,
		PublicAddr:  req.PublicAddr,
		Author:      req.Author,
		Title:       req.Title,
		Name:        req.Name,
		Logo:        req.Logo,
		Version:     req.Version,
		Description: req.Description,
		WebsiteUrl:  req.WebsiteUrl,
		Remotes:     remotes,
		Status:      string(req.Status),
		Type:        string(req.Type),
	}

	exist, existObject, err := agent.GetAgentByName(req.Name)
	if err != nil {
		log.WithError(err).Error("GetMcpByNameVersion")
		serializer.StringErrorRes(w, errors.New("GetMcpByNameVersion"))
		return
	}

	// 如果记录已经存在,取出id,放入newMcp中,这样就可以触发save自动更新
	if exist {
		content := []byte(fmt.Sprintf("%s_%s", req.Name, req.Nonce))
		if !crypto.VerifySignature(content, req.Sign, existObject.PublicAddr) {
			log.WithFields(logrus.Fields{"req addr": req.PublicAddr, "owner addr": existObject.PublicAddr}).Error("签名验证失败,您无权修改.")
			serializer.StringErrorRes(w, errors.New("签名验证失败,您无权修改."))
			return
		}
		log.WithFields(logrus.Fields{"name": req.Name, "publicAddr": req.PublicAddr, "existObject": existObject.PublicAddr}).Info("update")
		newOb.Id = existObject.Id
		newOb.UpdateDate = time.Now()
	} else {
		log.WithFields(logrus.Fields{"name": req.Name, "publicAddr": req.PublicAddr}).Info("insert")
		newOb.UpdateDate = time.Now()
		newOb.PublishDate = time.Now()
		newOb.ServerId = uuid.New().String()
		newOb.VersionId = uuid.New().String()
		newOb.UseTimes = 0
	}
	log.Info(newOb)
	log.Info("AdvertSlot:", newOb.AdvertSlot)
	err = newOb.Save()
	if err != nil {
		log.WithError(err).Error("Save")
		serializer.ErrorRes(w, err)
		return
	}

	serializer.SuccessRes(w, struct {
	}{})
}
