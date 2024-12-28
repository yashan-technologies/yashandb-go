



第三方库`github.com/jmoiron/sqlx`基于结构体封装了查询等方法，较好的提升了易用性，yasdb-go驱动支持通过该依赖库连接并操作yashandb，使用示例如下：

```go
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "git.yasdb.com/go/yasdb-go"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

type user struct {
	Id   int    `db:"ID"`
	Age  int    `db:"AGE"`
	Name string `db:"NAME"`
}

func InitDB(dsn string) (err error) {
	db, err = sqlx.Open("yasdb", dsn)
	if err != nil {
		fmt.Printf("connect server failed, err:%v\n", err)
		return
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(3)
	return
}

func exec(sqlstr string, args ...interface{}) (sql.Result, error) {
	fmt.Println(sqlstr)
	result, err := db.Exec(sqlstr, args...)
	if err != nil {
		log.Fatal(err)
	}
	return result, err
}

func Create() {
	exec("drop table if exists test_users")
	exec("create table test_users (id int PRIMARY KEY, name varchar(256), age int)")
}

func Insert(id int, name string, age int) {
	exec("INSERT INTO test_users(id, name, age) VALUES(?, ?, ?)", id, name, age)

}

func QueryMultiRow() {
	sqlStr := "SELECT id, age, name FROM test_users"
	var users []*user
	if err := db.Select(&users, sqlStr); err != nil {
		fmt.Printf("get data failed, err:%v\n", err)
		return
	}
	for i := 0; i < len(users); i++ {
		fmt.Printf("id:%d, name:%s, age:%d\n", users[i].Id, users[i].Name, users[i].Age)
	}
}

func main() {
	dsn := flag.String("dsn", "sys/sys@127.0.0.1:1688", "input you dsn(data source name) to connect yashandb.")
	flag.Parse()
	if err := InitDB(*dsn); err != nil {
		log.Fatal(err)
	}
	Create()
	Insert(1, "a", 3)
	db.MustBegin().Commit()
	QueryMultiRow()
	Insert(2, "b", 4)
	Insert(3, "张飞", 5)
	db.MustBegin().Commit()
	QueryMultiRow()
}
```

