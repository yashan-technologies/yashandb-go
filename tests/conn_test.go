package tests_test

import (
    "fmt"
    "testing"

    _ "cod-git.sics.com/cod-noah/yasdb-go"
    "github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
    ast := assert.New(t)
    _, err := ConnectDsn("sys/sys@127.0.0.1:1600")
    ast.NotNil(err, "connect to sys/sys@127.0.0.1:1600")

    db, err := Connect()
    ast.Nil(err, "connect to ", DSN)

    err = db.Close()
    ast.Nil(err, "close ", DSN)
}

func TestAutoCommitTrue(t *testing.T) {
    ast := NewTcheck(t, fmt.Sprintf("%s?autoCommit=1", DSN))

    tableName := "t_g_autocommit_true"
    ast.Exec(GenDropTableSQL(tableName))

    columns := Column{"id": "int"}
    ast.Exec(GenCreateTableSQL(tableName, columns))

    ast.Exec(fmt.Sprintf("insert into %s values(?)", tableName), 1)

    ast.DB.Close()

    ast.ConnectDB(DSN)
    count := 0
    err := ast.DB.QueryRow(fmt.Sprintf("select count(*) from %s", tableName)).Scan(&count)
    ast.Ast.Nil(err, "count 1")
    ast.Ast.Equal(count, 1)

    ast.DB.Close()
}

func TestAutoCommitFalse(t *testing.T) {
    ast := NewTcheck(t, fmt.Sprintf("%s?autoCommit=0", DSN))

    tableName := "t_g_autocommit_false"
    ast.Exec(GenDropTableSQL(tableName))

    columns := Column{"id": "int"}
    ast.Exec(GenCreateTableSQL(tableName, columns))

    ast.Exec(fmt.Sprintf("insert into %s values(?)", tableName), 1)

    db1 := ast.DB
    // close connection，will be commited.
    // so do NOT close it before check it.
    defer db1.Close()

    ast.ConnectDB(DSN)
    count := 0
    err := ast.DB.QueryRow(fmt.Sprintf("select count(*) from %s", tableName)).Scan(&count)
    ast.Ast.Nil(err, "count 1")
    ast.Ast.Equal(count, 0)

    ast.DB.Close()
}
