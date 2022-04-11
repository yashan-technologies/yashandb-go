# yasdb-go

## quick started

### 设置git

```bash
vim ~/.gitconfig

# 输入，我们的gitlab不支持https，所以需要把所有https转成http，http转成git就可以免密

[url "git@cod-git.sics.com:"]
    insteadOf = http://cod-git.sics.com/
[url "http://cod-git.sics.com"]
    insteadOf = https://cod-git.sics.com

```

### 设置环境变量

由于该工程采用的是cgo的方式进行开发，需要用到yasdb v0.21.1的两个so库

- libyas_infra.so.0
- libyascli.so.0

第三方库

- libcrypto.so.1.1

```bash
go env -w GOPRIVATE=cod-git.sics.com
export LD_LIBRARY_PATH=$GOPATH/pkg/mod/cod-git.sics.com/cod-noah/yasdb-go@v0.21.1/deps/lib:$LD_LIBRARY_PATH
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

    _ "cod-git.sics.com/cod-noah/yasdb-go"
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
go mod tidy

go run main.go
```

### 出包注意事项

> 出包需要将 `libyas_infra.so.0` 和 `libyascli.so.0` 同时打包进去，才能在其他机器运行。
