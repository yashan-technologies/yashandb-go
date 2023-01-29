# DSN 
`数据源名称(data source name, dsn)`是用来连接YashanDB，具有一定格式的字符串，里面包含了数据库用户名、密码、主机IP地址、端口号、连接参数等信息。  

## FORMAT
`DSN`的格式为：`user/password@host:port[?param1=value1&...&paramN=valueN]`  
- `user`: 数据库用户名，必填项
- `password`: 数据库用户密码，必填项
- `host`: 主机的ip地址，必填项
- `port`: 数据库服务端口，必填项
- `parameters`: 数据库连接参数，非必填项；目前只能设置`autocommit=true`,表示自动提交事务  

特殊字符说明：
- 如果`user`和`password`中有特殊字符`/`、`@`、`\`，需要使用符号`\`进行转义
- eg: 数据库用户`user`为`s/ys@\`，在`DSN`中使用`s\/ys\@\\`来表示，完整示例参照下一节

## EXAMPLE
| DSN                                                | user    | password   | host      | port | parameters      |
| -------------------------------------------------- | ------- | ---------- | --------- | ---- | --------------- |
| `sys/yasdb_123@127.0.0.1:1688`                     | sys     | yasdb_123  | 127.0.0.1 | 1688 | -               |
| `sys/yasdb_123@127.0.0.1:1688?autocommit=true`     | sys     | yasdb_123  | 127.0.0.1 | 1688 | autocommit=true |
| `sys/yasdb\@_123@127.0.0.1:1688`                   | sys     | yasdb@_123 | 127.0.0.1 | 1688 | -               |
| `sys\//yasdb\@_123@127.0.0.1:1688`                 | sys/    | yasdb@_123 | 127.0.0.1 | 1688 | -               |
| `sys\/\\/yasdb_123@127.0.0.1:1688?autocommit=true` | `sys/\` | yasdb_123  | 127.0.0.1 | 1688 | autocommit=true |
其中，`-`表示该参数未设置