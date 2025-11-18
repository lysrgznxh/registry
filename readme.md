
# 编译
cd cmd/registry
编译window版本: buildWin.bat
编译linux版本:  buildLinux.sh

# 启动
运行目录下创建或编辑config.toml

```
[base]
# 填写本地可以访问到的公网地址和端口
local_server_url = "http://192.168.3.35:52320"



[mysql]
# 以下为链接的mysql数据库配置
host = "localhost"
port = "53306"
user = "agent"
password = "VMDCVyDFJ3bsx7gjGboZk739tEH1lr"
database = "agent"
charset = "utf8"
```

registry.exe start


# 查看版本号
registry.exe version





# ------ agent登记 ------

## 说明
注册agent时,需要使用私钥进行签名,同时需带上公钥地址,系统会将agent(以name为依据)的所有权绑定到该地址上,后续的修改只能由该公钥地址的私钥签名才可以.


## 测试用例
```services/api/agent_register_test.go```

# ------ 图片上传 ------

## 说明
注册agent时,需要的图片从此接口进行上传.


## 测试用例
```services/api/agent_logo_upload_test.go```

# ------ agent列表 ------

## 说明
搜索登记的agent


## 测试用例
```services/api/agent_list_test.go```