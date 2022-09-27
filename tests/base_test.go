package tests_test

import (
    "database/sql"
    "fmt"
    "log"
    "strings"
    "testing"

    _ "git.yasdb.com/cod-noah/yasdb-go"
    "github.com/stretchr/testify/assert"
)

const DSN string = "sys/yasdb_123@192.168.6.177:1688"

func init() {
    log.SetFlags(log.Ltime)
}

type Column map[string]string

func Connect() (*sql.DB, error) {
    return ConnectDsn(DSN)
}

func ConnectDsn(dsn string) (*sql.DB, error) {
    db, _ := sql.Open("yasdb", dsn)
    return db, db.Ping()
}

func GenCreateTableSQL(name string, columns Column) string {
    cols := []string{}
    for k, t := range columns {
        cols = append(cols, fmt.Sprintf("%s %s", k, t))
    }
    return fmt.Sprintf("create table %s (%s)", name, strings.Join(cols, ", "))
}

func GenDropTableSQL(name string) string {
    return fmt.Sprintf("drop table if exists %s", name)
}

func LogSQL(query string, args ...interface{}) {
    if len(args) > 0 {
        log.Printf("\tsql: %s\t%v\n", query, args)
    } else {
        log.Printf("\tsql: %s\n", query)
    }
}

func LogRes(res sql.Result, query string, args ...interface{}) {
    c, _ := res.RowsAffected()
    if len(args) > 0 {
        log.Printf("\trows affect: %d\tsql: %s\t%v\n", c, query, args)
    } else {
        log.Printf("\trows affect: %d\tsql: %s\n", c, query)
    }
}

type Tcheck struct {
    DB  *sql.DB
    Ast *assert.Assertions
    Tx  *sql.Tx
}

func NewTcheck(t *testing.T, dsn ...string) *Tcheck {
    ast := &Tcheck{Ast: assert.New(t)}
    if len(dsn) > 0 {
        ast.ConnectDB(dsn[0])
    }
    return ast
}

func (t *Tcheck) ConnectDB(dsn string) {
    db, err := ConnectDsn(dsn)
    t.Ast.Nil(err, fmt.Sprintf("connect to %s", dsn))
    if err != nil {
        log.Fatal(err)
    }
    t.DB = db
}

func (t *Tcheck) Begin() *sql.Tx {
    tx, err := t.DB.Begin()
    t.Ast.Nil(err, "tx begin")
    t.Tx = tx
    return tx
}

func (t *Tcheck) GenTable(tableName string, columns Column) {
    t.Exec(GenDropTableSQL(tableName))
    t.Exec(GenCreateTableSQL(tableName, columns))
}

func (t *Tcheck) Exec(query string, args ...interface{}) (sql.Result, error) {
    res, err := t.DB.Exec(query, args...)
    t.Ast.Nil(err, fmt.Sprint(query, args))
    LogRes(res, query, args...)
    return res, err
}

func (t *Tcheck) TxExec(query string, args ...interface{}) (sql.Result, error) {
    res, err := t.Tx.Exec(query, args...)
    t.Ast.Nil(err, fmt.Sprint(query, args))
    LogRes(res, query, args...)
    return res, err
}

func (t *Tcheck) Query(query string, args ...interface{}) (*sql.Rows, error) {
    rows, err := t.DB.Query(query, args...)
    t.Ast.Nil(err, fmt.Sprint(query, args))
    LogSQL(query, args...)
    return rows, err
}

func (t *Tcheck) TxQuery(query string, args ...interface{}) (*sql.Rows, error) {
    rows, err := t.Tx.Query(query, args...)
    t.Ast.Nil(err, fmt.Sprint(query, args))
    LogSQL(query, args...)
    return rows, err
}

func (t *Tcheck) QueryRow(query string, args ...interface{}) *sql.Row {
    row := t.DB.QueryRow(query, args...)
    LogSQL(query, args...)
    return row
}

func (t *Tcheck) TxQueryRow(query string, args ...interface{}) *sql.Row {
    row := t.Tx.QueryRow(query, args...)
    LogSQL(query, args...)
    return row
}
