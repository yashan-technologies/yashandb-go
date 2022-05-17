package tests_test

import (
    "fmt"
    "testing"
    "time"

    _ "git.yasdb.com/cod-noah/yasdb-go"
)

func TestSelect(t *testing.T) {
    ast := NewTcheck(t, DSN)
    ast.Begin()
    dual := 0
    err := ast.TxQueryRow("select 1 from dual").Scan(&dual)
    ast.Ast.Nil(err, "select 1")
    ast.Ast.Equal(dual, 1, "select 1 from dual")
}

func TestReUse(t *testing.T) {
    ast := NewTcheck(t, DSN)

    ast.Begin()
    rows, err := ast.TxQuery("select * from v$instance")
    ast.Ast.Nil(err, "select v$instance")
    count := 0
    for rows.Next() {
        count++
    }
    ast.Ast.Equal(count, 1, "select v$instance")

    rows, err = ast.TxQuery("select * from v$instance")
    ast.Ast.Nil(err, "select v$instance")
    count = 0
    for rows.Next() {
        count++
    }
    ast.Ast.Equal(count, 1, "select v$instance")
}

func TestBind(t *testing.T) {
    ast := NewTcheck(t, DSN)
    defer ast.DB.Close()

    ast.Begin()
    rows, err := ast.TxQuery("select SID, USERNAME, STATUS from v$session where sid=?", 16)
    ast.Ast.Nil(err, "select with bind")
    for rows.Next() {
        sid, username, status := 0, "", ""
        err := rows.Scan(&sid, &username, &status)
        ast.Ast.Nil(err, "scan bind")
        ast.Ast.Equal(sid, 16)
    }
}

func TestBindParam(t *testing.T) {
    ast := NewTcheck(t, DSN)
    defer ast.DB.Close()

    tableName := "t_g_bind_param"
    ast.GenTable(tableName, Column{"a": "int", "b": "varchar(255)", "c": "double"})

    ast.Begin()
    ast.TxExec(fmt.Sprintf("insert into %s (a, b, c) values(:1, :2, 30)", tableName), 1, "aaa")
    ast.TxExec(fmt.Sprintf("insert into %s (a, b, c) values(:1, :2, 40)", tableName), 2, "bbb")
    ast.Tx.Commit()

    ast.Begin()
    rows, _ := ast.TxQuery(fmt.Sprintf("select a, b, c from %s", tableName))
    index := 1
    for rows.Next() {
        a, b, c := 0, "", 0
        rows.Scan(&a, &b, &c)
        ast.Ast.Equal(a, index)
        index++
    }
    ast.Begin()
    ast.TxExec(fmt.Sprintf("delete from %s where a=:1", tableName), 1)
    ast.Tx.Commit()

    ast.Begin()
    a, b, c := 0, "", 0
    ast.TxQueryRow(fmt.Sprintf("select a, b, c from %s", tableName)).Scan(&a, &b, &c)
    ast.Ast.Equal(a, 2)
    ast.Ast.Equal(b, "bbb")
    ast.Ast.Equal(c, 40)

    ast.Begin()
    ast.TxExec(fmt.Sprintf("update %s set c=:1", tableName), 50)
    ast.Tx.Commit()

    ast.Begin()
    a, b, c = 0, "", 0
    ast.TxQueryRow(fmt.Sprintf("select a, b, c from %s", tableName)).Scan(&a, &b, &c)
    ast.Ast.Equal(a, 2)
    ast.Ast.Equal(b, "bbb")
    ast.Ast.Equal(c, 50)

}

func TestCount(t *testing.T) {
    ast := NewTcheck(t, DSN)
    defer ast.DB.Close()

    tableName := "t_g_count"
    ast.GenTable(tableName, Column{"id": "int", "name": "varchar(255)"})

    ast.Begin()
    ast.TxExec(fmt.Sprintf("insert into %s (id, name) values(?, ?)", tableName), 1, "张三")
    ast.TxExec(fmt.Sprintf("insert into %s (id, name) values(?, ?)", tableName), 2, "李四")
    ast.Ast.Nil(ast.Tx.Commit())

    ast.Begin()
    count := 0
    err := ast.TxQueryRow(fmt.Sprintf("select count(*) from %s", tableName)).Scan(&count)
    ast.Ast.Nil(err, "test count")
    ast.Ast.Equal(count, 2)
}

func TestAllNum(t *testing.T) {
    ast := NewTcheck(t, DSN)
    defer ast.DB.Close()

    tableName := "t_g_all_num"
    columns := Column{
        "int1":    "tinyint",
        "int2":    "smallint",
        "int3":    "integer",
        "int4":    "bigint",
        "num5":    "number(10,5)",
        "float6":  "float",
        "double7": "double",
    }
    ast.GenTable(tableName, columns)

    ast.Begin()
    ast.TxExec(fmt.Sprintf("insert into %s (int1,int2,int3,int4,num5,float6,double7) values(?,?,?,?,?,?,?)", tableName), 1, 2, 3, 4, 5.0, 6.0, 7.0)
    ast.Tx.Commit()

    ast.Begin()
    var a, b, c, d int64
    var e, f, g float64
    rows, _ := ast.TxQuery(fmt.Sprintf("select int1, int2, int3, int4, num5, float6, double7 from %s limit ?", tableName), 1)
    for rows.Next() {
        rows.Scan(&a, &b, &c, &d, &e, &f, &g)
        ast.Ast.EqualValues(a, 1)
        ast.Ast.EqualValues(b, 2)
        ast.Ast.EqualValues(c, 3)
        ast.Ast.EqualValues(d, 4)
        ast.Ast.EqualValues(e, 5.0)
        ast.Ast.EqualValues(f, 6.0)
        ast.Ast.EqualValues(g, 7.0)
    }
}

func TestNameParam(t *testing.T) {
    ast := NewTcheck(t, DSN)
    defer ast.DB.Close()

    ast.Begin()
    rows, _ := ast.TxQuery("select * from v$session where sid=:id", 16)
    count := 0
    for rows.Next() {
        count++
    }
    ast.Ast.Equal(count, 1)
}

func TestFetchAll(t *testing.T) {
    ast := NewTcheck(t, DSN)
    defer ast.DB.Close()

    tableName := "t_g_fetch_all"
    columns := Column{
        "id":    "integer",
        "name":  "varchar(255)",
        "birth": "date",
    }
    ast.GenTable(tableName, columns)

    data := [][]interface{}{
        {1, "张三", time.Now()},
        {2, "李四", time.Now()},
        {3, "王五", time.Now()},
    }

    ast.Begin()
    for _, d := range data {
        ast.TxExec(fmt.Sprintf("insert into %s (id, name, birth) values(?, ?, ?)", tableName), d[0], d[1], d[2])
    }
    ast.Tx.Commit()

    ast.Begin()
    rows, _ := ast.TxQuery(fmt.Sprintf("select id, name, birth from %s", tableName))
    cols, _ := rows.Columns()
    fmt.Println(cols)
    index := 0
    for rows.Next() {
        id, name, birth := 0, "", time.Time{}
        rows.Scan(&id, &name, &birth)
        fmt.Println(id, name, birth.Format("2006-01-02 15:04:05"))
        ast.Ast.EqualValues(id, data[index][0])
        ast.Ast.EqualValues(name, data[index][1])
        ast.Ast.EqualValues(birth.Format("2006-01-02 15:04:05"), data[index][2].(time.Time).Format("2006-01-02 15:04:05"))
        index++
    }
}
