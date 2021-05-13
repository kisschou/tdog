# tdog

![go](https://img.shields.io/github/go-mod/go-version/kisschou/tdog?color=green&style=flat-square) ![commit](https://img.shields.io/github/last-commit/kisschou/tdog) ![licence](https://img.shields.io/github/license/kisschou/tdog?color=green) ![tdog](https://img.shields.io/badge/kisschou-tdog-green)

Just start in a smaile.

## Contents
-   [tdog](#tdog)
    -   [Contents](#contents)
    -   [About](#about)
    -   [Installation](#installation)
    -   [Quick start](#quick-start)
        -   [Use it for server:](#use-it-for-server)
        -   [Use it for scaffolding:](#use-it-for-scaffolding)
    -   [Functions](#functions)
        -   [1. MySQL Handler](#mysql-handler)
            -   [1.1 Struct](#struct)
            -   [1.2 NewMySQL() \*mySql](#newmysql-mysql)
            -   [1.3 (*mySql) Change(name string) *xorm.Engine](#mysql-changename-string-xorm.engine)
            -   [1.4 (*mySql) New(name string, conf *MySqlConf)\*xorm.Engine](#mysql-newname-string-conf-mysqlconf-xorm.engine)
            -   [1.5 Example](#example)
        -   [2. Redis Handler](#redis-handler)
        -   [3. Util Handler](#util-handler)
        -   [4. Snowflake Handler](#snowflake-handler)
        -   [5. Excel Handler](#excel-handler)
        -   [6. Config Handler](#config-handler)
        -   [7. Logger Handler](#logger-handler)
        -   [8. Validation Handler](#validation-handler)
            -   [8.1 Rule Struct](#rule-struct)
            -   [8.2 Rule description](#rule-description)
            -   [8.3 Functions and Usage](#functions-and-usage)
            -   [8.4 Validate Report Center and Validate Report](#validate-report-center-and-validate-report)
    -   [Contributing](#contributing)
    -   [Licence](#licence)


## About

[中文文档](./README_CN.md)

Core for all my project of golang. It's sort of a framework, but I feel like it's actually a scaffolding.

The Restful-Server section is referred to GIN. It's not powerful, but it's relatively complete.

<br />


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

<br />

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

<br />

## Functions

Functions are at the heart of scaffolding.

<br />

#### 1. MySQL Handler

The [XORM](https://github.com/go-xorm/xorm) package is used, so the detailed functions can be referred to its [documentation](https://pkg.go.dev/github.com/go-xorm/xorm).

Here are just a few things I feel I need:

- Eliminating the need to write initialization over and over again when using the same configuration for the same library within a framework

- Made a simple link pool, automatic loading configuration initialization engine, improve the reusability

<br />

##### 1.1 Struct

```go
mySql struct {
	engineList map[string]*xorm.Engine // engine pool
	Engine     *xorm.Engine            // current engine
}

MySqlConf struct {
	engine       string
	Host         string
	Port         string
	User         string
	Pass         string
	Db           string
	Charset      string
	Prefix       string
	dsn          string
	Debug        bool
	MaxIdleConns int
	MaxOpenConns int
}
```
MySql configuration instructions:

| field         | type     | desc                                                         |
| ------------- | -------- | ------------------------------------------------------------ |
| engine        | string   | Operating database engine.             |
| Host          | string   | Database connection address.       |
| Port          | string   | Database connection port. |
| User          | string   | Database connection account. |
| Pass          | string   | Database connection password. |
| Db            | string   | Database connection database. |
| Charset       | string   | Character set used by the database. |
| Prefix        | string   | The prefix of the table in the database. |
| dsn           | string   | Data source name. |
| Debug         | string   | Whether to enable debugging mode. |
| MaxIdleConns  | string   | The max idle connections on pool. |
| MaxOpenConns  | string   | The max open connections on pool. |

> When you customize the MySQL configuration file, you only need to set the following values:
> - Host
> - Port
> - User
> - Pass
> - Db
> - Charset
> - Prefix
> - Debug
> - MaxIdleConns
> - MaxOpenConns

If you are importing a database connection from a configuration file, the configuration file will be styled using TOML. The configuration items are as follows:

```toml
$ cat database.toml
## Whether to enable debugging mode.
debug = true

# master
[master] # engine name
host = "127.0.0.1" # connection address
port = "3306" # connection port
user = "root" # connection account
pass = "root" # connection password
db = "test_db" # connection database
charset = "utf8mb4" # Character set used by the database
prefix = "" # The prefix of the table in the database

# db read only
[master_read] # engine name
host = "127.0.0.1" # connection address
port = "3306" # connection port
user = "root" # connection account
pass = "root" # connection password
db = "test_db" # connection database
charset = "utf8mb4" # Character set used by the database
prefix = "" # The prefix of the table in the database
```

> See the configuration section for how to get them from the toml file.

<br />

##### 1.2 NewMySQL() *mySql

This function is used to initialize the mySql structure, which is the starting point and the core of everything.

After import tdog, this function is used as `tdog.NewMySQL()`

<br />

##### 1.3 (*mySql) Change(name string) *xorm.Engine

Use the tag name to switch the current database engine.

If the tag name does not exist in the engine group, it will go to the configuration file to obtain the corresponding configuration, and then return to the engine after construction.

<br />

##### 1.4 (*mySql) New(name string, conf *MySqlConf) *xorm.Engine

You can use this function to generate an engine through a custom configuration file.

<br />

##### 1.5 Example

```go
import "github.com/kisschou/tdog"

engine := tdog.NewMySQL().Engine // init a mysql engine use default configuation.

tab1Impl := new(tab1) // init table struct.
result, err := engine.Where("id=?", 1).Get(tab1Impl)
// Query SQL: SELECT * FROM tab1 WHERE id = 1 LIMIT 1;

// About Transaction
trans := engine.Session() // init a new transaction.
defer trans.Close()
affected, err := trans.InsertMulti([]*tab1)
if affected < 1 && err != nil {
	trans.Rollback()
}
trans.Commit()
```

> For more orm operation methods, please refer to the [document](https://pkg.go.dev/github.com/go-xorm/xorm).

<br />

#### 2. Redis Handler

The [go-redis](https://github.com/go-redis/redis) package is used, so the detailed functions can be referred to its [documentation](https://pkg.go.dev/github.com/go-redis/redis/v8).

Here are just a few things I feel I need:

- Eliminating the need to write initialization over and over again when using the same configuration for the same library within a framework

- Made a simple link pool, automatic loading configuration initialization engine, improve the reusability

<br />

##### 1.1 Struct

```go
type redisModel struct {
	engineList map[string]*redisImpl.Client // Engine Pool
	Engine     *redisImpl.Client            // Current Engine
	db         int                          // Current Db
}
```

If you are importing a redis connection from a configuration file, the configuration file will be styled using TOML. The configuration items are as follows:

```toml
$ cat cache.toml
[master]
host = "127.0.0.1" # Connection address
port = "6379" # Connection port
pass = "" # Connection password
pool_size = 10 # Connection pool size
```

> See the configuration section for how to get them from the toml file.

<br />

##### 1.2 NewRedis() *redisModel

This function is used to initialize the redisModel structure, which is the starting point and the core of everything.

After import tdog, this function is used as `tdog.NewRedis()`

<br />

##### 1.3 (*redisModel) Change(name string) *redisImpl.Client

Use the tag name to switch the current redis engine.

If the tag name does not exist in the engine group, it will go to the configuration file to obtain the corresponding configuration, and then return to the engine after construction.

<br />

##### 1.4 (*redisModel) New(name, host, port, pass string, poolSize int) *redisImpl.Client

You can use this function to generate an engine through a custom configuration file.

Passing parameter description:

| param    | type   | desc                           |
| -------- | ------ | ------------------------------ |
| name     | string | Set the connection name.       |
| host     | string | Set the connection host.       |
| port     | string | Set the connection port.       |
| pass     | string | Set the connection password.   |
| poolSize | int    | Sets the connection pool size. |

<br />

##### 1.5 (*redisModel) Db(index int) *redisImpl.Client

Accept an index to switch the libraries used by the current engine.

<br />

##### 1.6 Example

```go
import "github.com/kisschou/tdog"

engine := tdog.NewRedis().Engine // init a redis engine use default configuation.
result, err := engine.SetNX(tdog.Ctx, "test:key", "Hello World", time.Duration(60)*time.Second).Result() // Set key
engine.Get(tdog.Ctx, "test:key").String() // Get key
```

> For more orm operation methods, please refer to the [document](https://pkg.go.dev/github.com/go-xorm/xorm).

<br />

#### 3. Util Handler

<br />

#### 4. Crypt Handler

<br />

#### 5. Excel Handler

The use of [XLSX](https://github.com/tealeg/xlsx) package to do a convenient use of Excel processing tools, because do not know how to face the Excel processing, so just do a simple function.

<br />

###### 5.1 NewExcel(file string) *excel

Initialize an Excel module with the file, which is the starting point and the core of everything.

Note: that the file passed in here must contain the path.

<br />

###### 5.2 ( *excel) Get() \[\]\[\]\[\]string

Read all the data out of the file and return it.

<br />

###### 5.3 (*excel) Open() (excelImpl *xlsx.File)

Return to the operation controller of Excel file, refer to the detailed tutorial of XLSX for the specific operation functions.

<br />

#### 6. Config Handler



<br />

#### 7. Logger Handler

<br />

#### 8. Validation Handler

This is a parameter automatic validation module.

In a set of rules, fast, convenient, automatic calibration of the corresponding fields according to the rules, and feedback calibration results.

So, rules are important.

<br />

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

<br />

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

<br />

##### 8.3 Functions and Usage

In fact, what is said here is some obvious things, but afraid of a long time, they have forgotten, so record.

<br />

###### 8.3.1 NewValidate()

This function is used to initialize the validate structure, which is the starting point and the core of everything.

After import tdog, this function is used as `tdog.NewValidate()`

<br />

###### 8.3.1 (*validate) Rule(input []\*rule) *validate

This function is mainly used to inject a list of rules. The list of rules is derived from the `Rule Struct`.

<br />

###### 8.3.2 (*validate) Json(input string) *validate

This Function is also used to inject a list of rules, the difference is that the parameter passed to this function is `json`. If used in conjunction with the `Rule Function`, the latter will override the former.

<br />

###### 8.3.3 (*validate) Check(needle map[string]string) (output *validReport, err error)

Begin validation and return the `validation report structure` as soon as any data fails.

The received parameter is `map[string]string`. This is the collection of uplink data. If multidimensional array is involved, it needs to build its own loop processing.

What is returned is a structure of the `Validate report` and an `error interface`. Errors need to be handled by your own judgment. This is an old thing in golang. The use of the structure of the `Validate report` can be viewed in more detail below.

<br />


###### 8.3.4 (*validate) UninterruptedCheck(needle map[string]string) (output *validReportCenter, err error)

Start the Validate, regardless of whether it encounters an object that has an error, it will stubbornly insist on running all the data.

The received parameter is `map[string]string`. This is the collection of uplink data. If multidimensional array is involved, it needs to build its own loop processing.

It will pack all the Validate reports into the `Validate report center` and return. The detailed description of the calibration report center will look down. It also returns an `error interface` that needs to be handled by itself.

<br />


##### 8.4 Validate Report Center and Validate Report

The processing related to the Validate report cannot escape these two little guys.

<br />

###### 8.4.1 Validate Report

This structure has no sub-functions, but it has some parameters that can be used:

| param     | type     | desc                                                         |
| --------- | -------- | ------------------------------------------------------------ |
| Name      | string   | The key name of the query in the parameter list.             |
| Rule      | []string | The validation rules used can be found in the Rule Description section below. |
| Result    | bool     | The result of the verification, `True` means the verification is successful, `False` means the verification failed |
| Message   | string   | The text message feedback of the verification result is currently only fixed in Chinese, and customization is not supported. If necessary, you can handle it yourself according by `Result`. |

> You can use these parameters directly.
>
> If after a round of inspection and found that all have passed the verification, a report structure will be returned at this time, its Name and Rule are empty, Result is True, and Message is Success.
> So if you see a similar validate report, This can continue the following process.

<br />

###### 8.4.2 Validate Report Center

- (*validReportCenter) ReportList() []*validReport 

  > get all report from report center.

- (*validReportCenter) ReportByIndex(index int) *validReport

  > get the report by index. so given int index, returns `*report`.

- (*validReportCenter) ReportByName(name string) *validReport

  > get the report by name. so must given string name, and will returns `*report`

- (*validReportCenter) BuildTime() string

  > get build time from report center.

- (*validReportCenter) ElapsedTime() int64

  > get elapsed time from report center. The return value is nanoseconds.

- (*validReportCenter) ToJson() string

  > convert to json and return.

  <br />


## Contributing

Let's have a good time together!!!

- Fork the Project
- Create your Feature Branch (git checkout -b feature/AmazingFeature)
- Commit your Changes (git commit -m 'Add some AmazingFeature')
- Push to the Branch (git push origin feature/AmazingFeature)
- Open a Pull Request

<br />


## Licence

Copyright (c) 2020-present Kisschou.