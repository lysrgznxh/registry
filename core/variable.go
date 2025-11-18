package core

// 配置项，暴露给外部进行使用,使用之前 需要先调用 LoadConfig 进行加载
var (
	config *Config

	BaseConf *BaseConfig

	//mysql配置
	MysqlConf *MysqlConfig
)
