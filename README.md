# yasdb-go

### quick started

yasdb-go驱动包可以通过以下两种方式进行安装和配置：

#### 在线下载

- 通过go get命令在`https://git.yasdb.com/go/yasdb-go`仓库中下载
- [linux](./linux_online.md) 
- [windows](./win_online.md) 

#### 离线获取

- 如果无法连接上`https://git.yasdb.com/go/yasdb-go`，可以使用离线方案，直接获得源码包
- [linux](./linux_offline.md) 
- [windows](./win_offline.md) 

### DSN的填写说明
参考教程: [DSN format](./DSNFormat.md)


### yacapi更新

git subtree pull --prefix=yacapi  git@git.yasdb.com:cod-x/yacapi.git  master --squash


### 使用示例
yasdb-go驱动包中提供了一些使用示例。其中通过标准库`database/sql`连接操作yashandb示例请参考[examples](./_examples)；通过第三方库`github.com/jmoiron/sqlx`连接操作yashandb示例请参考[sqlx示例](./sqlx示例.md) 。