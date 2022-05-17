# yasdb-go

## quick started

### 设置git

```bash
vim ~/.gitconfig

# 输入，我们的gitlab不支持https，所以需要把所有https转成http，再把http转成git就可以免密

[url "git@git.yasdb.com:"]
    insteadOf = https://git.yasdb.com/
```

### 设置环境变量

由于该工程采用的是cgo的方式进行开发，需要用到yasdb v0.21.1的两个so库

- libyas_infra.so.0
- libyascli.so.0

第三方库

- libcrypto.so.1.1

```bash
# 将我们的gitlab设置为私有仓
go env -w GOPRIVATE=git.yasdb.com

# 设置yasdb 客户端所需的lib加载的路径，yasdb-go默认带了v21.1版本的lib库
# 也可自行指向其他带lib库的路径
# 注意：
#    1. v21.1  支持访问v22.1的yasdb
#    2. v22.1不支持访问v21.1的yasdb

export LD_LIBRARY_PATH=$GOPATH/pkg/mod/git.yasdb.com/cod-noah/yasdb-go@v1.0.1/deps/lib:$LD_LIBRARY_PATH
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

    _ "git.yasdb.com/cod-noah/yasdb-go"
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
# 我们的gitlab不支持https，所以使用go mod tidy会导致失败
# 需要提前使用go get -insecure，以http的方式下载，
# 注意必须设置私有仓：go env -w GOPRIVATE=git.yasdb.com

go get -insecure git.yasdb.com/cod-noah/yasdb-go
go mod tidy

go run main.go
```

### 出包注意事项

> 出包需要将 `libyas_infra.so.0` 和 `libyascli.so.0` 同时打包进去，才能在其他机器运行。
