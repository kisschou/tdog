# tdog

![go](https://img.shields.io/github/go-mod/go-version/kisschou/tdog?color=green&style=flat-square) ![commit](https://img.shields.io/github/last-commit/kisschou/tdog) ![licence](https://img.shields.io/github/license/kisschou/tdog?color=green) ![tdog](https://img.shields.io/badge/kisschou-tdog-green)

Just start in a smaile.


## About

[中文文档](./README_CN.md)

Core for all my project of golang. It's sort of a framework, but I feel like it's actually a scaffolding.

The Restful-Server section is referred to GIN. It's not powerful, but it's relatively complete.


## Installation

To install Tdog package, you need to install Go and set your Go workspace first.

  1. Got it by go command:
  ```
    $ go get -u github.com/kisschou/tdog
  ```
  2. Import it in your code:
  ```
    import "github.com/kisschou/tdog"
  ```


## Quick start

#### Use it for server:

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

#### Use it for scaffolding:

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


## Functions

Functions are at the heart of scaffolding.

#### MySQL Handler

#### Redis Handler

#### Util Handler

#### Snowflake Handler

#### Excel Handler

#### Config Handler

#### Logger Handler



## Licence

Copyright (c) 2020-present Kisschou.
