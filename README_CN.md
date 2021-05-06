# tdog

![go](https://img.shields.io/github/go-mod/go-version/kisschou/tdog?color=green&style=flat-square) ![commit](https://img.shields.io/github/last-commit/kisschou/tdog) ![licence](https://img.shields.io/github/license/kisschou/tdog?color=green) ![tdog](https://img.shields.io/badge/kisschou-tdog-green)

仅仅是一个微笑的开始.


## 介绍

这是我所有的golang项目的核心包。也许这个项目可以算是一个框架，但是我更觉得它是一个脚手架。

`Restful-Server` 模块是参考 `GIN` 做的. 也许它不是很强大, 但是也基本上可以算是完善。


## 安装

首先你必须安装golang和golang的开发环境。

  1. 通过go的命令获取它:
  ```
    $ go get -u github.com/kisschou/tdog
  ```
  2. 在你的项目中导入它:
  ```
    $ import "github.com/kisschou/tdog"

  ```

## 快速使用

#### 用来起服务:

```
$ cat example.go
```

```go
package main

import (
    "net/http"

    "github.com/kisschou/tdog"
)

func main() {
    r := tdog.New()

    r.GET("/ping", func(httpUtil *tdog.HttpUtil) {
        httpUtil.Res.JSON(http.StatusOK, tdog.H{
            "message": "Pong GET",
        })
    })

    r.Run() // Start server.
}
```

#### 当脚手架使用:

```go
package main

import (
    "github.com/kisschou/tdog"
)

func main() {
    id, err := tdog.NewSnowflake(1, 1, 1).Get()
    if err != nil {
        log.Fatalln(err)
    }
    log.Printf(id)
}
```

## 函数们

函数是身为脚手架的核心所在.

#### MySQL

#### Redis

#### 基本函数

#### Snowflake

#### Excel

#### 配置文件

#### 日志


## 证书

Copyright (c) 2020-present Kisschou.
