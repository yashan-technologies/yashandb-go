# online version instructions

## quick started

### 设置git

```bash
vim ~/.gitconfig

# https转成git，结合gitlab的SSH密钥配置就可以免密

[url "git@git.yasdb.com:"]
    insteadOf = https://git.yasdb.com/
```

### 设置环境变量

#### go env

设置go的环境变量

```bash
# 将我们的gitlab设置为私有仓
go env -w GOPRIVATE=git.yasdb.com
```

#### LD_LIBRARY_PATH

由于该工程采用的是cgo的方式进行开发，需要用到YashanDB的两个so库

- libyas_infra.so.0
- libyascli.so.0

第三方库

- libcrypto.so.1.1
- libpcre.so

```bash
# 设置yasdb客户端所需的lib加载的路径
export LD_LIBRARY_PATH=<yasdb的lib库文件绝对路径>:$LD_LIBRARY_PATH
```

### 创建项目

```bash
mkdir yasdb_connect && cd yasdb_connect

vim main.go
# 输入以下代码
```

### 编码开发

```go
package main

import (
    "database/sql"
    "log"

    _ "git.yasdb.com/go/yasdb-go"
)

func Connect() *sql.DB {
    db, err := sql.Open("yasdb", "sys/sys@127.0.0.1:1688")
    if err != nil {
        log.Fatalf("some error %s", err.Error())
    }
    return db
}

type Database struct {
    Status string
    Role   string
    Point  string
}

func main() {
    db := Connect()
    var s Database
    err := db.QueryRow("select STATUS, DATABASE_ROLE, FLUSH_POINT from V$DATABASE where STATUS = ?", "NORMAL").Scan(&s.Status, &s.Role, &s.Point)
    if err != nil {
        log.Fatal("some wrong for query", err.Error())
    }
    if s.Status != "NORMAL" {
        log.Fatal(s.Status, " is not equal")
    }
    log.Println(s.Status, s.Role, s.Point)
}
```

### 编译执行

```bash
go mod init yasdb_connect

# 必须在init后执行，go get默认使用https下载包。

go get git.yasdb.com/go/yasdb-go@Commit # 即当前最新commit号
go mod tidy

go run main.go
```

### 出包注意事项

> 出包需要将 `libyas_infra.so.0` 和 `libyascli.so.0` 同时打包进去，才能在其他机器运行。