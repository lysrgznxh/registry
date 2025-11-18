package agent

import (
	"agent-network-protocol/registry/core"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
	"xorm.io/builder"
)

// 根据主键查询
func GetAgentByName(name string) (exist bool, data Agent, err error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmDbUser)
	var where = builder.NewCond()
	where = where.And(builder.Eq{"name": name})
	exist, err = agentDb.Where(where).Get(&data)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"name": name}).Error("db.get")
	}
	return
}

// 根据主键查询agent
func GetAgentById(id int) (exist bool, data Agent, err error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmDbUser)
	exist, err = agentDb.ID(id).Get(&data)
	if err != nil {
		log.WithError(err).WithField("id", id).Error("db.get")
	}
	return
}

// 搜索agent
func SearchAgent(page int, pagesize int, type_, advertSlot, categoryStr string, names []string, status string) (target []Agent, err error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmDbUser)
	where := ""
	whereCases := []string{}
	if type_ != "" {
		whereCases = append(whereCases, fmt.Sprintf("type='%s'", type_))
	}
	if advertSlot != "" {
		whereCases = append(whereCases, fmt.Sprintf("advert_slot=%s", advertSlot))
	}
	if categoryStr != "" {
		whereCases = append(whereCases, fmt.Sprintf("category=%s", categoryStr))
	}
	if len(names) > 0 {
		nameStrCase := ""
		for k, name := range names {
			if k > 0 {
				nameStrCase = nameStrCase + ","
			}
			nameStrCase += fmt.Sprintf("'%s'", name)
		}
		whereCases = append(whereCases, fmt.Sprintf("name in(%s)", nameStrCase))
	}
	if status != "" {
		whereCases = append(whereCases, fmt.Sprintf("status='%s'", status))
	}
	if page == 0 {
		return nil, errors.New("page 必须大于 0")
	}
	if pagesize < 10 {
		return nil, errors.New("pagesize 必须大于等于 10")
	}
	begin := (page - 1) * pagesize

	if len(whereCases) > 0 {
		where = ""
		for _, whereCase := range whereCases {
			if where != "" {
				where += " and "
			}
			where = where + whereCase
		}
		where = fmt.Sprintf("where %s", where)
	}
	sql := fmt.Sprintf(`SELECT * FROM agent %s limit %d,%d`, where, begin, pagesize)
	list, err := agentDb.Query(sql)
	if err != nil {
		log.WithError(err).Error("db.find")
		return
	}
	layout := "2006-01-02 15:04:05"
	for _, v := range list {
		id, _ := strconv.Atoi(string(v["id"]))
		category, _ := strconv.Atoi(string(v["category"]))
		useTimes, _ := strconv.Atoi(string(v["use_times"]))
		aiLevel, _ := strconv.Atoi(string(v["ai_level"]))
		evalScore, _ := strconv.Atoi(string(v["eval_score"]))
		advertSlot, _ := strconv.Atoi(string(v["advert_slot"]))
		price, _ := strconv.ParseFloat(string(v["price"]), 64)
		// 解析为time.Time类型
		publishDate, _ := time.Parse(layout, string(v["publish_date"]))
		updateDate, _ := time.Parse(layout, string(v["update_date"]))
		target = append(target, Agent{
			Id:          id,
			Name:        string(v["name"]),
			Title:       string(v["title"]),
			Logo:        string(v["logo"]),
			Author:      string(v["author"]),
			AdvertSlot:  advertSlot,
			Category:    category,
			UseTimes:    useTimes,
			AiLevel:     aiLevel,
			EvalScore:   evalScore,
			Type:        string(v["type"]),
			RepoUrl:     string(v["repo_url"]),
			RepoSource:  string(v["repo_source"]),
			Version:     string(v["version"]),
			Description: string(v["description"]),
			WebsiteUrl:  string(v["website_url"]),
			Packages:    string(v["packages"]),
			Remotes:     string(v["remotes"]),
			Status:      string(v["status"]),
			PublishDate: publishDate,
			UpdateDate:  updateDate,
			ServerId:    string(v["server_id"]),
			VersionId:   string(v["version_id"]),
			Price:       price,
		})
	}
	return
}
