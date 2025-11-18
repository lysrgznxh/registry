package core

//项目常量在此登记

const (
	MysqlDir     = "mysql"
	ClientLogDir = "log"

	APiServicePort = 52320 //本地接口服务端口号
)

type ModelName string

const (
	CategroyBase       int = 1 // 基础
	CategroyYingXiao   int = 2 // 营销
	CategroyShuJuFenXi int = 3 // 数据分析
)
