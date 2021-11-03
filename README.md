# yasdb-go

## quick started

### 下载代码

```bash
go get -insecure cod-git.sics.com/cod-noah/yasdb-go@v0.1.0
```

 - v0.1.0 主要是br分支，对应的TP版本号为v0.1.x
 - v0.2.0 主要是master分支，对应的TP版本号为v0.2.x

### 设置环境变量

由于该工程采用的是cgo的方式进行开发，需要用到TP的两个so库

- libyas_infra.so.0
- libyascli.so.0

```bash
export LD_LIBRARY_PATH=$GOPATH/pkg/mod/cod-git.sics.com/cod-noah/yasdb-go@v0.1.0/deps/lib:$LD_LIBRARY_PATH
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
    db, err := sql.Open("yasdb", "sys/sys@192.168.30.219:16001")
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

### 出包注意事项

> 出包需要将 `libyas_infra.so.0` 和 `libyascli.so.0` 同时打包进去，才能在其他机器运行。
