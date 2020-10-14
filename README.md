# tdog-core
仅仅是一个微笑开始

## 说明
本项目主要是用于wordgame项目的核心包. 主要是参考gin改造而成，改动太多，难以详尽描述。使用go-mod做包管理。

## 引用包
包名 | 版本 | 用途
:--: | :--: | :--:
[go-redis](https://github.com/go-redis/redis) | v7.2.0 | 处理redis连接和redis方法实现
[xorm](https://github.com/xormplus/xorm) | v0.7.9 | 数据库的orm处理，当前主要针对mysql数据库使用
[httprouter](https://github.com/julienschmidt/httprouter) | v1.3.0 | 据说能提高40倍效率的玩意，用于路由处理
[logrus](https://github.com/sirupsen/logrus) | v1.5.0 | 日志处理
[viper](https://github.com/spf13/viper) | v1.6.3 | 配置文件获取，配置项读取，感觉用起来会有卡顿和偶发性的取key失败的问题，有更好的替代方案时考虑更换
[go-swagger](https://github.com/go-swagger/go-swagger) | v0.23.0 | 因为当前版本只针对RESTFul所以swagger还是一个不错的工具
[godaemon](https://github.com/icattlecoder/godaemon) | v0.0.0-20190426080617-f87981e709a1 | 以后台形式运行golang

## 结构
```
└── tdog                                // 核心功能实现
    ├── core                            // 核心包
    │   ├── controller.go               // 基础控制器类，希望每个控制器都引用他
    │   ├── error.go                    // 错误处理
    │   ├── feign.go                    // Feign请求转发处理脚本
    │   ├── jwt.go                      // 一个垃圾到可能会被嘲讽的东西
    │   ├── model.go                    // 基础模型类，希望每个模型都引用他
    │   ├── request.go                  // 请求数据类，所谓的统一入口
    │   ├── response.go                 // 返回数据类，所谓的统一出口
    │   ├── router.go                   // 路由解析类，httprouter作用于此
    │   ├── service.go                  // 基础服务类，希望每个服务都引用他
    │   └── websocket.go                // WebSocket服务类
    │   ├── captcha.go                  // 验证码图片生成
    │   ├── config.go                   // 配置文件获取类，viper作用于此
    │   ├── crypt.go                    // 加密方法都放在这里
    │   ├── file.go                     // 上传文件接收存储等
    │   ├── http_request.go             // 对外模拟请求方法
    │   ├── logger.go                   // 日志类, logrus作用于此
    │   ├── mysql.go                    // MySQL操作类, xorm作用于此
    │   ├── redis.go                    // Redis操作类, go-redis作用于此
    │   ├── snowflake.go                // 雪花算法
    │   └── util.go                     // 基础方法都在这里
    └──  tests                          // 单元测试
        └── main.go
```

## 使用说明
```
## 安装golang环境等问题不再重复叙述

## clone项目
shell> git clone https://github.com/kisschou/word-game-go.git
shell> cd word-game-go

## 运行项目
shell> go run .
## 后台运行
shell> go run -d=true .

## 项目打包
## ...Linux or Mac OS:
shell> go build .
## 会在目录下生成一个名为wordgame的二进制文件
## 运行:
shell> ./wordgame
## 后台运行
shell> ./wordgame -d=true
## ...Windows:
shell> CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build .
## 会生成一个名为 wordgame.exe 的可执行文件
## 在Windows下双击即可运行

```
> 本项目在windows下运行会报config找不到的问题，这是因为Linux和Windows下路径分隔符"/"和"\\"的问题。

#### 业务添加说明
* 路由
> 路由添加请参考 [app/routers/README.md](#)，里面有详细描述。

* 控制器
> 不懂怎么给描述，索性直接上demo
```
package controllers

import (
    "net/http"

    "wordgame/tdog/services"
    "wordgame/tdog/core"
)

type Demo struct() {
    Base core.Controller // 这里需引用基础控制器文件
}

func (demo *Demo) Hello() {
    name := ""

    // 因为目前没有做自动验证，所以得自己做参数判断
    if _, ok := demo.Base.Req.Params["name"]; ok {
        if len(demo.Base.Req.Params["name"]) > 0 {
            name = demo.Base.Req.Params["name"][0]
        }
    }

    // 请求services
    DemoService := new(services.Demo)
    name = DemoService.GetName(name)

    // 数据返回
    // 返回json
    member.Base.Res.JSON(http.StatusOK, core.H{
        "success" => "Hello " + name,
    })
    // 返回string
    member.Base.Res.String(http.StatusOK, "Hello " + name)
    // 返回xml
    ...
    // 返回data
    ...
}
```

* 服务
> 不懂怎么给描述，索性直接上demo
```
package services

import (
    "time"

    "wordgame/app/models"
    "wordgame/tdog/core"
)

type (
    DemoInfo string

    Demo struct {
        Base core.Service // 这里需引用基础服务文件
    }
)

func (demo *Demo) GetName(name string) (retStr DemoInfo) {
    demo.Base.Redis.NewEngine() // 初始化redis控件
    if demo.Base.Redis.Engine.Exists(name).Val() > 0 { // 判断key是否存在redis
        retStr = demo.Base.Redis.Engine.Get(name).Val() // 存在取出数据返回
        return
    }

    // 不存在，从数据库中获取
    DemoModel := new(models.DemoModel)
    retStr = DemoModel.GetName(name)

    // 数据存储到redis
    // 更多的redis操作请参考go-redis文档
    demo.Base.Redis.Engine.SetNX(name, retStr, time.Duration(0)*time.Second)

    return
}
```

* 模型
> 模型的生成建议使用xorm通过表结构自动生成模型后编写;
> xorm生成模型参考: https://github.com/go-xorm/cmd/blob/master/README.md
```
package models

import (
    "wordgame/tdog/core"
)

type (
    // 这个是模块定义的结构体
    DemoModel struct {
        Base core.Model // 这里需引用基础模型文件
    }

    // 数据库Column的结构体
    // 可以通过xorm工具生成
    Demo struct {
        Name       string      `xorm:"comment('测试名') unique CHAR(40)"`
        RetName    string      `xorm:"comment('测试返回名') CHAR(40)"`
    }
)

func (demoModel *DemoModel) GetName(name) (retStr string) {
    demoInfo := new(Demo)
    demoModel.Base.Sql.NewEngine() // 初始化数据库驱动
    // 从数据库中读取数据
    // 更多的数据库操作请参考xorm文档
    result, _ := demoModel.Base.Sql.Engine.Where("name=?", name).Get(demoInfo)
    if !result {
        // 数据不存在
        return
    }
    retStr = demoInfo.RetName
    return
}
```

## Swagger
> 这里使用的是`go-swagger`包

#### 安装
```
Windows:
Command> go get -u github.com/go-swagger/go-swagger/cmd/swagger
## 然后设置swagger目录到环境变量,如果swagger成功唤醒请忽略

Linux:
shell> go get -u github.com/go-swagger/go-swagger/cmd/swagger
shell> export PATH=$GOPATH/bin:$PATH

MacOS:
shell> brew tap go-swagger/go-swagger
shell> brew install go-swagger
```
> 也可以直接去[官方下载](https://github.com/go-swagger/go-swagger/releases)

#### 使用
```
## 1. 在对外的路由指向的controller中用注释的方式标识接口信息
// swagger:operation POST /member/login member login
// ---
// summary: 用户登录
// description: 用户登录
// parameters:
// - name: Authorization
//   in: header
//   description: 授权信息
//   type: string
//   required: true
// - name: username
//   in: body
//   description: 用户名
//   type: string
//   required: true
// - name: password
//   in: body
//   description: 密码
//   type: string
//   required: true
// responses:
//   200: repoResp
//   401: badReq

#### 解释
// swagger:operaion [POST:请求方式(可以是GET\PUT\DELETE...)] [url:请求地址] [标签] [用于此端点的请求]
// 注: 最后两个理解为id和节点, 用于标注地址
// --- 这个部分下面是YAML格式的swagger规范.确保您的缩进是一致的和正确的
// summary: 标题
// description: 描述
// parametres:   下面是参数了
// - name: 参数名
    in: [header|body|query] 参数的位置
    description: 描述
    type: 类型
    required: 是否必须
// responses: 响应

## 2. 根目录下运行, 会自动遍历go文件生成接口文件
shell> swagger generate spec -o ./swagger.json

## 3. 启动swagger服务
shell> swagger serve -F=swagger ./swagger.json
```

## 其他
Kisschou&copy;2020.All Rights.
