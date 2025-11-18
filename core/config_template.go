package core

import (
	"html/template"
	"strings"
)

const defaultConfigTemplate = `# This is a TOML config file.


[nxnos]
local_server_url = "{{.NxnosConfig.LocalServerUrl}}"

[mysql]
host = "{{.MysqlConfig.Host}}"
port = "{{.MysqlConfig.Port}}"
user = "{{.MysqlConfig.User}}"
password = "{{.MysqlConfig.Password}}"
database = "{{.MysqlConfig.Database}}"
charset = "{{.MysqlConfig.Charset}}"
`

var configTemplate *template.Template

func initConfigTemplate() error {
	var err error
	tmpl := template.New("configFileTemplate").Funcs(template.FuncMap{
		"StringsJoin": strings.Join,
	})
	if configTemplate, err = tmpl.Parse(defaultConfigTemplate); err != nil {
		return err
	}
	return nil
}
