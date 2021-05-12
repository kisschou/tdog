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

#### 1. MySQL Handler



#### 2. Redis Handler



#### 3. Util Handler



#### 4. Snowflake Handler



#### 5. Excel Handler



#### 6. Config Handler



#### 7. Logger Handler



#### 8. Validation Handler

This is a parameter automatic validation module.

In a set of rules, fast, convenient, automatic calibration of the corresponding fields according to the rules, and feedback calibration results.

So, rules are important.



##### 8.1 Rule Struct

```go
Rule struct {
		Name      string   `json:"name"`     // key's name
		ParamType string   `json:"type"`     // val's type
		IsMust    bool     `json:"is_must"`  // is must be set
		Rule      []string `json:"validate"` // vaildate rule
}
```

| field     | type     | desc                                                         |
| --------- | -------- | ------------------------------------------------------------ |
| Name      | string   | The key name of the query in the parameter list.             |
| ParamType | string   | The type of value to be queried in the parameter list.       |
| IsMust    | bool     | Specifies that the key name must exist in the parameter list. |
| Rule      | []string | The validation rules used can be found in the Rule Description section below. |

> IsMust takes precedence over all rules. All rules take precedence and are executed in the specified order.



##### 8.2 Rule description

| keyword             | description                                                  | example           |
| ------------------- | ------------------------------------------------------------ | ----------------- |
| empty               | Determine if the value is empty.If it's a number, it's less than or equal to 0 | -                 |
| phone               | Determine whether the content of the value is a phone number. Currently, only Chinese mobile phone numbers are supported. | -                 |
| email               | Determine whether the contents of the value are mailbox addresses. | -                 |
| scope(x, y)         | Set a range to specify a reasonable length for a string or number. | scope(1,10)       |
| enum(str1,str2,...) | Set an enumeration to constrain the contents of the value.   | enum(header,body) |
| date                | Determine whether the value conforms to the date format, which is yyyy-mm-dd. | -                 |
| datetime            | Check whether the value conforms to the date and time format, which is: yyyy-mm-dd HH: MM :ss. | -                 |
| sensitive-word      | The desensitized word list retrieves whether the content contains sensitive words. This item is temporarily invalid. | -                 |

> You can select multiple rules to constrain a field at the same time, such as' []string{" empty ", 'email', 'scope(10,)'} 'to specify a value that also satisfies:
> · Can't be empty
> · Is the correct email address
> · Length greater than 10

> The scope(x,y) rule can set both x and y values, or only one of them. When you set the value of x, you're adding a condition that's greater than x, and the same thing for the value of y, when you set it, you're adding a condition that's less than y.



##### 8.3 Functions and Usage

In fact, what is said here is some obvious things, but afraid of a long time, they have forgotten, so record.



###### 8.3.1 NewValidate()

This function is used to initialize the validate structure, which is the starting point and the core of everything.

After import tdog, this function is used as `tdog.newvalidate()`



###### 8.3.1 (*validate) Rule(input []\*rule) *validate

This function is mainly used to inject a list of rules. The list of rules is derived from the `Rule Struct`.



###### 8.3.2 (*validate) Json(input string) *validate

This Function is also used to inject a list of rules, the difference is that the parameter passed to this function is `json`. If used in conjunction with the `Rule Function`, the latter will override the former.



###### 8.3.3 (*validate) Check(needle map[string]string) (output *validReport, err error)



###### 8.3.4 (*validate) UninterruptedCheck(needle map[string]string) (output *validReportCenter, err error)



##### 8.4 Validate Report Center and Validate Report



## Licence

Copyright (c) 2020-present Kisschou.
