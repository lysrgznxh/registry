package core

import (
	"agent-network-protocol/registry/core/model"
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"
)

//运行控制公共变量,可通过包名直接使用,该变量的初始化由

func init() {

}

// 对应数据库中的对象
type AgentHttpRes struct {
	Id          int               `json:"id"`
	Name        string            `json:"name"`   // 唯一标识
	Title       string            `json:"title"`  // 名称
	Author      string            `json:"author"` //作者
	Logo        string            `json:"logo"`
	RepoUrl     string            `json:"repo_url"`    // 仓库地址
	RepoSource  string            `json:"repo_source"` // 仓库来源
	Version     string            `json:"version"`
	Description string            `json:"description"` // 描述
	WebsiteUrl  string            `json:"website_url"` // 官方网站
	Packages    []model.Package   `json:"packages"`
	Remotes     []model.Transport `json:"remotes"`
	Status      string            `json:"status"`
	PublishDate time.Time         `json:"publish_date"` // 发布时间
	UpdateDate  time.Time         `json:"update_date"`  // 更新时间
	ServerId    string            `json:"server_id"`
	VersionId   string            `json:"version_id"`
	Type        string            `json:"type"`        // 类型可能是 agent 或者 mcp
	AiLevel     int               `json:"ai_level"`    // ai级别 可能是  0  1(l1) 2(l2) 3(l3) 4(l4) 5(l5)
	EvalScore   int               `json:"eval_score"`  // 评价分数
	UseTimes    int               `json:"use_times"`   // 使用次数
	Category    int               `json:"category"`    // 所属分类
	AdvertSlot  int               `json:"advert_slot"` //推荐位置  0:没有推荐位置  非0值为对应的推荐位
	Price       float64           `json:"price"`       //价格
}

// 对外提供的结构
type RepoStatus struct {
	RepoSize   int64  `json:"RepoSize"`
	StorageMax int64  `json:"StorageMax"`
	NumObjects int64  `json:"NumObjects"`
	RepoPath   string `json:"RepoPath"`
}

type ServmOutput struct {
	Status int    `json:"status"`
	Info   string `json:"info"`
	Data   string `json:"data"`
}

// 返回对象的名称/类型,只能传递非指针类型，如果是指针类型需要手工用*转回来
func GetObjectName(ob interface{}) string {
	pv1 := reflect.ValueOf(ob)
	var object reflect.Type
	if pv1.Kind() == reflect.Ptr {
		object = pv1.Type().Elem()
		//fmt.Println("a:",object,"n:",object.Name())
	} else {
		object = reflect.TypeOf(ob)
		//fmt.Println("b:",object,"n:",object.Name())
	}
	return object.Name()
}

// 返回对应结构的当前函数名(包含包名)
func GetStructFuncName(object interface{}) string {
	structName := reflect.TypeOf(object).String() //结构体名称
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	arry := strings.Split(f.Name(), ".")
	funcName := arry[len(arry)-1] //函数名称
	return structName + "." + funcName
}

// 返回当前函数名称(不含包名)
// 适合在中小模块在单层包结构下进行使用
func GetFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	arry := strings.Split(f.Name(), ".")
	return arry[len(arry)-1]
}

// 返回当前函数名称(含包名)
// 适合在大模块下多层包结构下进行使用
// 返回最后一层 包名 +  函数名
func GetPackageFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	list := strings.Split(f.Name(), "/")
	return list[len(list)-1]
}

// 检查str是否存在于slice中
func ContainsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// IsValidVersion 验证版本号是否符合 0.0.0 格式（由数字和点构成，且每个数字段为 0）
func IsValidVersion(version string) bool {
	// 定义正则表达式：匹配 "0.0.0" 格式（三个数字段均为 0，由点分隔）
	pattern := `^\d+\.\d+\.\d+$`
	matched, err := regexp.MatchString(pattern, version)
	if err != nil {
		// 正则表达式编译错误（理论上不会发生，除非 pattern 语法错误）
		return false
	}
	return matched
}

func AgentTypeCheck(typeStr string) bool {
	if typeStr != string(model.AgentTypeAgent) && typeStr != string(model.AgentTypeMcp) {
		return false
	}
	return true
}

func GetMysqlDbConnInfo() string {
	if MysqlConf == nil {
		return ""
	}
	//"agent:@tcp(localhost:53306)/agent?charset=utf8"
	return fmt.Sprintf("%s@tcp(%s:%s)/%s?charset=%s", MysqlConf.User, MysqlConf.Host, MysqlConf.Port, MysqlConf.Database, MysqlConf.Charset)
}

func GetMysqlDbPasswordConnInfo() string {
	if MysqlConf == nil {
		return ""
	}
	//"agent:password@tcp(localhost:53306)/agent?charset=utf8"
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", MysqlConf.User, MysqlConf.Password, MysqlConf.Host, MysqlConf.Port, MysqlConf.Database, MysqlConf.Charset)
}
