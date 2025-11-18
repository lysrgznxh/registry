package core

// 返回默认配置
func GetDefaultConfig() *Config {

	baseConfig := &BaseConfig{
		LocalServerUrl: "",
	}

	mysqlConfig := &MysqlConfig{
		Host:     "localhost",
		Port:     "53306",
		User:     "agent",
		Password: "VMDCVyDFJ3bsx7gjGboZk739tEH1lr",
		Database: "agent",
		Charset:  "utf8",
	}

	return &Config{
		BaseConfig:  baseConfig,
		MysqlConfig: mysqlConfig,
	}
}
