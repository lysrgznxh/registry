package core

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	BaseConfig  *BaseConfig  `mapstructure:"base"`  //nxnos
	MysqlConfig *MysqlConfig `mapstructure:"mysql"` //mysql
}

type BaseConfig struct {
	LocalServerUrl string `mapstructure:"local_server_url"` //本地可以访问的服务地址,用于合成图片访问地址
}

type MysqlConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Charset  string `mapstructure:"charset"`
}

// 同步配置以便更新写入文件
func SyncConfig() {
	config.BaseConfig = BaseConf
	config.MysqlConfig = MysqlConf
}

// 更新配置变量
func updateConfigVariable(conf *Config) error {
	config = conf
	BaseConf = conf.BaseConfig
	MysqlConf = conf.MysqlConfig
	return nil
}

// 从指定目录中加载配置文件
func LoadConfig(configFile string) error {
	log := BuildLog(GetPackageFuncName(), LmCore)
	var conf *Config
	err := initConfigTemplate()
	if err != nil {
		return err
	}

	workingDirectory, err := os.Getwd()
	if err != nil {
		return err
	}
	if configFile != "" {
		dir, _ := filepath.Split(configFile)
		workingDirectory = dir
	} else {
		configFile = filepath.Join(workingDirectory, "config.toml")
	}
	fmt.Println("configFile", configFile)
	fmt.Println("workingDirectory", workingDirectory)
	conf = GetDefaultConfig()
	_, err = os.Stat(configFile)
	if os.IsNotExist(err) { //如果配置文件不存在,使用默认配置
		log.Info("config file not found, use default config")
		err = writeConfigFile(configFile, conf)
		if err != nil {
			return err
		}
		//更新包变量
		return updateConfigVariable(conf)
	}

	log.Info("load config file from ", configFile)
	rootViper := viper.New()
	rootViper.SetConfigType("toml")
	rootViper.SetConfigName("config")
	rootViper.AddConfigPath(workingDirectory)

	if err := rootViper.ReadInConfig(); err != nil {
		return err
	}

	if err := rootViper.Unmarshal(conf); err != nil {
		return err
	}

	//更新包变量
	return updateConfigVariable(conf)
}

// 从配置文件中读取配置
func ReadConfig(configFile string) (*Config, error) {
	log := BuildLog(GetPackageFuncName(), LmCore)
	var conf *Config
	err := initConfigTemplate()
	if err != nil {
		return nil, err
	}

	workingDirectory, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if configFile != "" {
		dir, _ := filepath.Split(configFile)
		workingDirectory = dir
	} else {
		configFile = filepath.Join(workingDirectory, "config.toml")
	}

	conf = GetDefaultConfig()
	log.Info("load config file from ", configFile)
	rootViper := viper.New()
	rootViper.SetConfigType("toml")
	rootViper.SetConfigName("config")
	rootViper.AddConfigPath(workingDirectory)

	if err := rootViper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := rootViper.Unmarshal(conf); err != nil {
		return nil, err
	}

	return conf, nil
}

// 保存配置文件
func SaveConfig() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	configFile := filepath.Join(pwd, "config.toml")
	return writeConfigFile(configFile, config)
}

func writeConfigFile(configFilePath string, config *Config) error {
	//将路径中的\替换成\\
	//if !strings.Contains(config.Base.ChainRepoPath, "\\\\") {
	//	config.Base.ChainRepoPath = strings.ReplaceAll(config.Base.ChainRepoPath, "\\", "\\\\")
	//}
	var buffer bytes.Buffer
	if err := configTemplate.Execute(&buffer, config); err != nil {
		return err
	}
	//对替换完的内容进行编码转换  &#34;  替换为  ”
	replaceString := strings.ReplaceAll(buffer.String(), "&#34;", `\"`)

	return os.WriteFile(configFilePath, []byte(replaceString), 0644)
}
