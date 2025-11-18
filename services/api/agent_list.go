package api

import (
	"agent-network-protocol/registry/core"
	"agent-network-protocol/registry/core/model"
	"agent-network-protocol/registry/db/agent"
	"agent-network-protocol/registry/util/serializer"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

// agent 信息查询
func agentList(w http.ResponseWriter, r *http.Request) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmServApi)
	// 设置跨域支持
	w.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有来源访问
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	params := r.URL.Query()
	nameStr := params.Get("name")         // 根据唯一键 name  查询
	typeStr := params.Get("type")         // 类型
	advertSlot := params.Get("ad")        // 广告位
	categoryStr := params.Get("category") // 分类
	pageStr := params.Get("page")
	pagesizeStr := params.Get("pagesize")
	if pageStr == "" {
		pageStr = "1"
	}
	if pagesizeStr == "" {
		pagesizeStr = "10"
	}
	names := []string{}
	if nameStr != "" {
		names = strings.Split(nameStr, ",")
	}
	if typeStr != "" {
		ok := core.AgentTypeCheck(typeStr)
		if !ok {
			log.Error("无效的参数type:", typeStr)
			serializer.ErrorRes(w, errors.New("参数type只能为nxn或者mcp"))
			return
		}
	}
	if advertSlot != "" {
		_, err := strconv.Atoi(advertSlot)
		if err != nil {
			log.WithError(err).Error("strconv.Atoi")
			serializer.ErrorRes(w, errors.New("ad参数无效"))
			return
		}
	}
	if categoryStr != "" {
		_, err := strconv.Atoi(categoryStr)
		if err != nil {
			log.WithError(err).Error("strconv.Atoi")
			serializer.ErrorRes(w, errors.New("category参数无效"))
			return
		}
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		log.WithError(err).Error("strconv.Atoi")
		serializer.ErrorRes(w, errors.New("page参数无效"))
		return
	}
	pagesize, err := strconv.Atoi(pagesizeStr)
	if err != nil {
		log.WithError(err).Error("strconv.Atoi")
		serializer.ErrorRes(w, errors.New("pagesize参数无效"))
		return
	}

	list, err := agent.SearchAgent(page, pagesize, typeStr, advertSlot, categoryStr, names, "active")
	if err != nil {
		log.WithError(err).Error("user.SearchAgent")
		serializer.ErrorRes(w, err)
		return
	}
	serializer.SuccessRes(w, dbAgent2Agent(list))
}

func dbAgent2Agent(list []agent.Agent) []core.AgentHttpRes {
	servers := []core.AgentHttpRes{}
	for i := 0; i < len(list); i++ {

		packages := []model.Package{}
		if list[i].Packages != "" {
			_ = json.Unmarshal([]byte(list[i].Packages), &packages)
		}
		remotes := []model.Transport{}
		if list[i].Remotes != "" {
			_ = json.Unmarshal([]byte(list[i].Remotes), &remotes)
		}
		servers = append(servers, core.AgentHttpRes{
			Id:          list[i].Id,
			Name:        list[i].Name,
			RepoUrl:     list[i].RepoUrl,
			RepoSource:  list[i].RepoSource,
			Version:     list[i].Version,
			Description: list[i].Description,
			WebsiteUrl:  list[i].WebsiteUrl,
			Packages:    packages,
			Remotes:     remotes,
			Status:      list[i].Status,
			PublishDate: list[i].PublishDate,
			UpdateDate:  list[i].UpdateDate,
			ServerId:    list[i].ServerId,
			VersionId:   list[i].VersionId,
			Type:        list[i].Type,
			AiLevel:     list[i].AiLevel,
			EvalScore:   list[i].EvalScore,
			Category:    list[i].Category,
			Title:       list[i].Title,
			Logo:        list[i].Logo,
			Author:      list[i].Author,
			UseTimes:    list[i].UseTimes,
			AdvertSlot:  list[i].AdvertSlot,
			Price:       list[i].Price,
		})
	}
	return servers
}
