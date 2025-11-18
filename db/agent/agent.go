package agent

import (
	"agent-network-protocol/registry/core"
	"errors"
	"github.com/sirupsen/logrus"
	"time"
)

// agent服务列表
type Agent struct {
	Id          int       `json:"id" xorm:"not null pk autoincr INT(11)"`
	Name        string    `json:"name" xorm:"unique VARCHAR(100)"` // agent标识 , 不可重复
	Title       string    `json:"title" xorm:"default ''"`         // 标题,用于app显示
	Author      string    `json:"author" xorm:"default ''"`        // 作者
	Logo        string    `json:"logo" xorm:"default ''"`          // agent 图片
	RepoUrl     string    `json:"repo_url" xorm:"default ''"`      // 仓库地址
	RepoSource  string    `json:"repo_source" xorm:"default ''"`   // 仓库来源
	Version     string    `json:"version"`
	Description string    `json:"description"`
	WebsiteUrl  string    `json:"website_url"` //官方网站ad slot
	Packages    string    `json:"packages" xorm:"longtext"`
	Remotes     string    `json:"remotes" xorm:"longtext"`
	Status      string    `json:"status"`
	PublishDate time.Time `json:"publish_date" xorm:"TIMESTAMP"`
	UpdateDate  time.Time `json:"update_date" xorm:"TIMESTAMP"`
	ServerId    string    `json:"server_id"`
	VersionId   string    `json:"version_id"`
	Type        string    `json:"type"`                                 // 类型可能是 agent 或者 mcp
	Category    int       `json:"category" xorm:"INT(10) default 0"`    // 分类  0:默认 1:聊天机器人 2:智能写作 3:办公自动化 4:智能投资 5:个性推荐 6:医疗诊断 7:企业战略 8:金融风控 9:智能家居 10:智能客服  11:客户管理
	AiLevel     int       `json:"ai_level" xorm:"INT(10) default 0"`    // ai级别 可能是  0  1(l1) 2(l2) 3(l3) 4(l4) 5(l5)
	EvalScore   int       `json:"eval_score" xorm:"INT(10) default 0"`  // 评价分数
	PublicAddr  string    `json:"public_addr" xorm:"default ''"`        // 发布者的公钥地址
	UseTimes    int       `json:"use_times" xorm:"INT(10) default 0"`   // 使用次数
	AdvertSlot  int       `json:"advert_slot" xorm:"INT(10) default 0"` // 推荐位置  0:没有推荐位置  非0值为对应的推荐位
	Price       float64   `json:"price" xorm:"decimal(10,2) default 0"` // 价格
}

func (c *Agent) Save() error {
	log := core.BuildLog(core.GetStructFuncName(c), core.LmDbUser)
	if c.Id == 0 {
		return c.Insert()
	} else {
		_, err := agentDb.ID(c.Id).Update(c)
		if err != nil {
			log.WithError(err).Error("db.update")
		}
		return err
	}
}

func (c *Agent) Insert() error {
	log := core.BuildLog(core.GetStructFuncName(c), core.LmDbUser)
	_, err := agentDb.Insert(c)
	if err != nil {
		log.WithError(err).Error("db.insert")
	}
	return err
}

func (c *Agent) Delete() (err error) {
	log := core.BuildLog(core.GetStructFuncName(c), core.LmDbUser)
	if c.Id != 0 {
		_, err = agentDb.ID(c.Id).Delete(c)
		if err != nil {
			log.WithError(err).WithFields(logrus.Fields{
				"id": c.Id,
			}).Error("db.delete")
		}
	} else {
		err = errors.New("请设置有效的id再进行删除")
	}
	return
}
