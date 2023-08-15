# offline version instructions

## quick started

### 获取yasdb-go源码包

- 先从内网的git仓获取yasdb-go的源码包: [http://git.yasdb.com/go/yasdb-go](http://cod-git.sics.com/go/yasdb-go)

- 将源码解压到本地路径

### 设置环境变量

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

# 使用 go mod的replace命令，将安装包指向 本地路径
go mod edit -replace git.yasdb.com/go/yasdb-go@v1.0.1=<yasdb-go源码包解压后的绝对路径>
go mod tidy

go run main.go
```

### 出包注意事项

> 出包需要将 `libyas_infra.so.0` 和 `libyascli.so.0` 同时打包进去，才能在其他机器运行。