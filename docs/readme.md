# 接口使用说明

`` 接口服务器的地址 : http://127.0.0.1:52320``

## 1.获取agent列表

### 接口路径 /v1/agent/list

### 参数列表

- name **[参数说明]**:Agent的唯一编号 **[是否必填项]**:`否` **[默认值]**:`空` **[类型]**:`字符串`
- type **[参数说明]**:类型,可选范围(agent|mcp) **[是否必填项]**:`否` **[默认值]**:`空` **[类型]**:`字符串`
- page **[参数说明]**:页数 **[是否必填项]**:`否` **[默认值]**:`1` **[类型]**:`字符串`
- pagesize **[参数说明]**:每页记录数 **[是否必填项]**:`否` **[默认值]**:`10` **[类型]**:`字符串`

### 返回结构
```
{
    "status": 1,
    "info": "",
    "data": [
            {
            "id": 1823,
            "name": "myname.shop",
            "title": "商城助手",
            "author": "作者",
            "logo": "http://127.0.0.1：52320/logos/335ca168-65e0-445b-8963-524643997590.png",
            "repo_url": "",
            "repo_source": "",
            "version": "0.0.1",
            "description": "商品发布与管理",
            "website_url": "",
            "packages": [],
            "remotes": [
                {
                    "type": "sse",
                    "url": "http://123.123.123.123:9000/v1/chat-messages",
                    "headers": [
                        {
                            "description": "api key",
                            "isRequired": true,
                            "format": "string",
                            "value": "Bearer app-DEpbuy3UttZc3hRR7PIKhMjY",
                            "default": "Bearer app-DEpbuy3UttZc3hRR7PIKhMjY",
                            "name": "Authorization"
                        }
                    ],
                    "jsons": [
                        {
                            "description": "要查询的内容",
                            "isRequired": true,
                            "format": "string",
                            "name": "query"
                        },
                        {
                            "description": "用户自定义文件列表",
                            "format": "json",
                            "variables": {
                                "variable_name1": {
                                    "format": "json",
                                    "isArray": true,
                                    "variables": {
                                        "transfer_method": {
                                            "description": "传输方法",
                                            "format": "string",
                                            "choices": [
                                                "remote_url",
                                                "local_file"
                                            ]
                                        },
                                        "type": {
                                            "description": "文件类型",
                                            "format": "string",
                                            "choices": [
                                                "image",
                                                "document",
                                                "audio",
                                                "video",
                                                "custom"
                                            ]
                                        },
                                        "upload_file_id": {
                                            "description": "上传文件id",
                                            "format": "string"
                                        },
                                        "url": {
                                            "description": "图片地址",
                                            "format": "file_path"
                                        }
                                    }
                                }
                            },
                            "name": "inputs"
                        }
                    ]
                }
            ],
            "members": [],
            "status": "active",
            "publish_date": "2025-09-27T16:50:00Z",
            "update_date": "2025-10-08T17:21:30Z",
            "server_id": "1b382f70-4420-463c-97ae-4de737cdac09",
            "version_id": "4d97d4e1-5cb9-426e-87b8-626abf4e5285",
            "type": "nxn",
            "ai_level": 0,
            "eval_score": 0,
            "use_times": 0,
            "category": 2,
            "advert_slot": 0,
            "price": 0
        }
    ]
}
```

### 返回结构说明
```
{
"status": 1,
"info": "",
"data": []
}
```
status 表示是否请求成功,1 成功 , 0失败 ，当失败时请读取 info 字段获取错误信息

data 为返回的agent列表

agent 结构参考
```
{
    "id": 系统流水号,
    "name": "智能体应用ID,类似安卓的包名,保持全局唯一,不可重复,格式:英文+数字+.",
    "title": "商城助手",
    "author": "作者",
    "logo": "图片地址,可通过上传接口获取路径,或者使用互联网的图片",
    "repo_url": "",
    "repo_source": "",
    "version": "版本号",
    "description": "智能体说明",
    "website_url": "",
    "packages": [],
    "remotes": [{
        "type": "sse",
        "url": "智能体地址",
        "headers": [
            {
                "description": "头信息的说明",
                "isRequired": true,
                "format": "string",
                "value": "固定值",
                "default": "默认值",
                "name": "头信息字段名"
            }
        ],
        "jsons": [
            {
                "description": "字符串字段说明",
                "isRequired": true,
                "format": "string",
                "name": "字符串字段名称"
            },
            {
                "description": "结构体字段说明",
                "format": "json",
                "name": "inputs"
                "variables": {
                    "结构体字段子属性名1": {
                        "format": "string",
                        "isArray": false,
                    },
                    "结构体字段子属性名2": {
                        "format": "string",
                        "isArray": false,
                    }
                }
            }
        ]
    }]
}
```



## 2.上传agent形象图

### 接口路径 /v1/agent/logo/upload

### 参数列表

- image **[参数说明]**:以form格式保存的图片字段 **[是否必填项]**:`是` **[类型]**:`form字段`


### 返回结构说明
```
{
    "status": 1,
    "info": "",
    "data": {
        "filename":"可访问的图片路径",
        "size":文件大小,
    }
}
```
status 表示是否请求成功,1 成功 , 0失败 ，当失败时请读取 info 字段获取错误信息



## 3.发布agent

### 接口路径 /v1/agent/list

### 请求参数

```
{
    "category": 2,
    "token": "身份信息",
    "author": "智能体作者",
    "title": "智能体名称",
    "logo": "智能体图片",
    "name": "智能体应用名称(只能是英文字母+数字)",
    "version": "1.0.1",
    "description": "根据关键字搜索音乐",
    "website_url": "",
    "remotes": [
        {
            "type": "streamable-http",
            "url": "智能体访问地址",
            "headers": [
                {
                    "description": "智能体访问令牌",
                    "isRequired": true,
                    "format": "string",
                    "value": "令牌值",
                    "default": "令牌值",
                    "name": "令牌字段名"
                }
            ],
            "jsons": [
                {
                    "description": "要查询的内容",
                    "isRequired": true,
                    "format": "string",
                    "name": "query"
                },
                {
                    "description": "用户id",
                    "isRequired": true,
                    "format": "string",
                    "name": "user"
                },
                {
                    "description": "会话id",
                    "isRequired": true,
                    "format": "string",
                    "name": "conversation_id"
                },
                {
                    "description": "输出模式",
                    "isRequired": true,
                    "format": "string",
                    "value": "streaming",
                    "name": "response_mode"
                },
                {
                    "description": "用户自定义结构",
                    "format": "json",
                    "variables": {
                        "shopAccesstoken": {
                            "description":"身份令牌",
                            "format": "string",
                            "isRequired": true,
                        },
                        "team_workers": {
                            "description":"团队成员",
                            "format": "string",
                            "isRequired": true,
                        }
                    },
                    "name": "inputs"
                }
            ]
        }
    ],
    "status": "active",
    "type": "nxn"
}
```

### 返回结构说明
```
{
    "status": 1,
    "info": "",
    "data": {}
}
```
status 表示是否请求成功,1 成功 , 0失败 ，当失败时请读取 info 字段获取错误信息