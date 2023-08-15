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

#### 安装GCC
编译cgo需要使用64位的gcc编译器，可以从[tdm-gcc](https://jmeubank.github.io/tdm-gcc/download)或者[mingw-w64](https://www.mingw-w64.org/downloads/)下载并安装。
以通过tdm-gcc安装gcc为例：

1. 运行tdm-gcc的安装程序，选择安装目录。
![安装目录](images/gcc%E5%AE%89%E8%A3%85.png)

2. 完成安装后，tdm-gcc安装程序会自动将目录下的bin文件夹（C:\TDM-GCC-64\bin）添加到环境变量中，如图所示。如果是其他方法安装mingw，则需要注意手动添加环境变量。
![gcc添加环境变量](images/gcc-%E7%8E%AF%E5%A2%83%E5%8F%98%E9%87%8F.png)

3. 重新打开一个CMD窗口，输入`gcc --version`，没有报错则安装成功。

#### LD_LIBRARY_PATH

为了正确使用go驱动，需要安装YashanDB C驱动客户端并设置环境变量。
1. 在产品软件包或安装目录的Drivers文件夹中，查找yashandb-client-版本号-win-x86_64.zip。
2. 下载并解压到本地路径，如D:\yashandb\yashandb-client。
3. 设置环境变量PATH，指向该文件路径下的lib文件夹，如D:\yashandb\yashandb-client\lib。

![c驱动添加环境变量](/images/c驱动-环境变量.png)



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

go mod edit -replace git.yasdb.com/go/yasdb-go@v1.0.1=<yasdb-go源码包解压后的绝对路径>
go mod tidy

go run main.go
```

### 出包注意事项

> 出包需要将 `libyas_infra.so.0` 和 `libyascli.so.0` 同时打包进去，才能在其他机器运行。