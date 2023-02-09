package yasdb

import (
    "database/sql"
    "testing"
)

func TestConnect(t *testing.T) {
    db, err := sql.Open("yasdb", testDsn)
    if err != nil {
        t.Fatalf("error connecting: %v", err)
    }
    if err = db.Close(); err != nil {
        t.Fatalf("error db close: %v", err)
    }
}

/*
must set env:
LD_LIBRARY_PATH
YASDB_DATA
*/
func TestConnectWithoutPassword(t *testing.T) {
    db, err := sql.Open("yasdb", "")
    if err != nil {
        t.Fatalf("error connecting: %v", err)
    }
    if err = db.Close(); err != nil {
        t.Fatalf("error db close: %v", err)
    }
}

func TestPing(t *testing.T) {
    db, err := sql.Open("yasdb", testDsn)
    if err != nil {
        t.Fatalf("error connecting: %v", err)
    }
    if err = db.Ping(); err != nil {
        t.Fatalf("error db ping: %v", err)
    }
    if err = db.Close(); err != nil {
        t.Fatalf("error db close: %v", err)
    }
}

func TestAutoCommitTrue(t *testing.T) {
    runsqlTestACTrue(t, testAutoCommitTrue)
}

func testAutoCommitTrue(t *sqlTest) {
    si := sqlGenInfo{}
    t.sqlGenInfo = &si

    si = sqlGenInfo{
        tableName: "test_auto_commit_true",
        columnNameType: [][2]string{
            {"id", "int"},
            {"name", "varchar(20)"},
        },
        execArgs: [][]interface{}{
            {1, "column1"},
            {2, "column2"},
            {3, "column3"},
        },
        queryResult: [][]interface{}{
            {int32(1), "column1"},
            {int32(2), "column2"},
            {int32(3), "column3"},
        },
    }
    t.genTableTest()
    t.runInsertTest()
    t.runSelectTest()
}

func TestAutoCommitFase(t *testing.T) {
    runsqlTestACFalse(t, testAutoCommitFalse)
}

func testAutoCommitFalse(t *sqlTest) {
    si := sqlGenInfo{}
    t.sqlGenInfo = &si

    si = sqlGenInfo{
        tableName: "test_auto_commit_false",
        columnNameType: [][2]string{
            {"id", "int"},
            {"name", "varchar(20)"},
        },
        execArgs: [][]interface{}{
            {1, "column1"},
            {2, "column2"},
            {3, "column3"},
        },
        queryResult: [][]interface{}{
            {int32(1), "column1"},
            {int32(2), "column2"},
            {int32(3), "column3"},
        },
    }
    t.genTableTest()
    t.runInsertTest()
    t.runSelectTest()
}
