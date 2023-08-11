# windows online version instructions

## quick started

### 设置git

在git安装目录的etc文件夹下的config文件中添加：

```
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

为了正确使用go驱动，需要安装YashanDB C驱动客户端并设置环境变量。
1. 在产品软件包或安装目录的Drivers文件夹中，查找yashandb-client-版本号-win-x86_64.zip。
2. 下载并解压到本地路径，如D:\yashandb\yashandb-client。
3. 设置环境变量PATH，指向该文件路径下的lib文件夹，如D:\yashandb\yashandb-client\lib。

![Alt text](/images/添加环境变量.png)


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

# 必须在init后执行。
# 注意必须设置私有仓：go env -w GOPRIVATE=git.yasdb.com

go get git.yasdb.com/go/yasdb-go@Commit # 即当前最新commit号
go mod tidy

go run main.go
```

### 出包注意事项

> 出包需要将 `libyas_infra.so.0` 和 `libyascli.so.0` 同时打包进去，才能在其他机器运行。