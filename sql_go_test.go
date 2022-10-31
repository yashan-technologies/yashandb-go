//
// The following test cases are from https://github.com/bradfitz/go-sql-test
//
package yasdb

import (
    "database/sql"
    "fmt"
    "math/rand"
    "strings"
    "testing"
)

func TestBlob_gosql(t *testing.T) {
    t.Parallel()
    runSqlTest(t, testBlob_gosql)
}

func testBlob_gosql(t *sqlTest) {
    t.sqlGenInfo = &sqlGenInfo{}
    t.tableName = tablePrefix + "blob"
    t.columnNameType = [][2]string{
        {"id", "integer primary key"},
        {"bar", "blob"},
    }
    t.genTableTest()

    blob := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
    t.mustExec(fmt.Sprintf("insert into %s (id,bar) values(?, ?)", t.tableName), 1, blob)

    want := fmt.Sprintf("%x", blob)
    b := make([]byte, 16)
    err := t.QueryRow(fmt.Sprintf("select bar from %s where id = ?", t.tableName), 1).Scan(&b)
    got := fmt.Sprintf("%x", b)
    if err != nil {
        t.Errorf("[]byte scan: %v", err)
    } else if got != want {
        t.Errorf("for []byte, got %q; want %q", got, want)
    }

    got = ""
    err = t.QueryRow(fmt.Sprintf("select bar from %s where id = ?", t.tableName), 1).Scan(&got)
    want = string(blob)
    if err != nil {
        t.Errorf("string scan: %v", err)
    } else if got != want {
        t.Errorf("for string, got %q; want %q", got, want)
    }
}

func TestClob_gosql(t *testing.T) {
    t.Parallel()
    runSqlTest(t, testClob_gosql)
}

func testClob_gosql(t *sqlTest) {
    t.sqlGenInfo = &sqlGenInfo{}
    t.tableName = tablePrefix + "clob"
    t.columnNameType = [][2]string{
        {"id", "integer primary key"},
        {"bar", "clob"},
    }
    t.genTableTest()

    want := strings.Repeat("你好，YashanDB！", 10) + "......"
    t.mustExec(fmt.Sprintf("insert into %s (id,bar) values(?, ?)", t.tableName), 1, want)

    got := ""
    err := t.QueryRow(fmt.Sprintf("select bar from %s where id = ?", t.tableName), 1).Scan(&got)
    if err != nil {
        t.Errorf("string scan: %v", err)
    } else if got != want {
        t.Errorf("for string, got %q; want %q", got, want)
    }
}

func TestManyQueryRow_gosql(t *testing.T) {
    t.Parallel()
    runSqlTest(t, testManyQueryRow_gosql)
}

func testManyQueryRow_gosql(t *sqlTest) {
    if testing.Short() {
        t.Logf("skipping in short mode")
        return
    }
    t.sqlGenInfo = &sqlGenInfo{}
    t.tableName = tablePrefix + "manyqueryrow"
    t.columnNameType = [][2]string{
        {"id", "integer primary key"},
        {"name", "varchar(50)"},
    }
    t.genTableTest()

    t.mustExec(fmt.Sprintf("insert into %s (id, name) values(?,?)", t.tableName), 1, "bob")
    var name string
    for i := 0; i < 10000; i++ {
        err := t.QueryRow(fmt.Sprintf("select name from %s where id = ?", t.tableName), 1).Scan(&name)
        if err != nil || name != "bob" {
            t.Fatalf("on query %d: err=%v, name=%q", i, err, name)
        }
    }
}

func TestTxQuery_gosql(t *testing.T) {
    t.Parallel()
    runSqlTest(t, testTxQuery_gosql)
}

func testTxQuery_gosql(t *sqlTest) {
    t.sqlGenInfo = &sqlGenInfo{}
    t.tableName = tablePrefix + "txquery"
    t.columnNameType = [][2]string{
        {"id", "integer primary key"},
        {"name", "varchar(50)"},
    }
    t.genTableTest()

    tx, err := t.Begin()
    if err != nil {
        t.Fatal(err)
    }
    defer tx.Rollback()

    want := "你好，YahsanDB！"
    _, err = tx.Exec(fmt.Sprintf("insert into %s (id, name) values(?,?)", t.tableName), 1, want)
    if err != nil {
        t.Fatal(err)
    }

    r, err := tx.Query(fmt.Sprintf("select name from %s where id = ?", t.tableName), 1)
    if err != nil {
        t.Fatal(err)
    }
    defer r.Close()

    if !r.Next() {
        if r.Err() != nil {
            t.Fatal(err)
        }
        t.Fatal("expected one rows")
    }

    var got string
    err = r.Scan(&got)
    if err != nil {
        t.Errorf("string scan: %v", err)
    } else if got != want {
        t.Errorf("for string, got %q; want %q", got, want)
    }
}

func TestPrepareStmt_gosql(t *testing.T) {
    t.Parallel()
    runSqlTest(t, testPrepareStmt_gosql)
}

func testPrepareStmt_gosql(t *sqlTest) {
    t.sqlGenInfo = &sqlGenInfo{}
    t.tableName = tablePrefix + "preparestmt"
    t.columnNameType = [][2]string{{"count", "int"}}
    t.genTableTest()

    sel, err := t.Prepare(fmt.Sprintf("SELECT count FROM %s ORDER BY count DESC", t.tableName))
    if err != nil {
        t.Fatalf("prepare 1: %v", err)
    }
    ins, err := t.Prepare(fmt.Sprintf("INSERT INTO %s (count) VALUES (?)", t.tableName))
    if err != nil {
        t.Fatalf("prepare 2: %v", err)
    }

    for n := 1; n <= 3; n++ {
        if _, err := ins.Exec(n); err != nil {
            t.Fatalf("insert(%d) = %v", n, err)
        }
    }

    const nRuns = 10
    ch := make(chan bool)
    for i := 0; i < nRuns; i++ {
        go func() {
            defer func() {
                ch <- true
            }()
            for j := 0; j < 10; j++ {
                count := 0
                if err := sel.QueryRow().Scan(&count); err != nil && err != sql.ErrNoRows {
                    t.Errorf("Query: %v", err)
                    return
                }
                if _, err := ins.Exec(rand.Intn(100)); err != nil {
                    t.Errorf("Insert: %v", err)
                    return
                }
            }
        }()
    }
    for i := 0; i < nRuns; i++ {
        <-ch
    }
}

func TestEmoji_gosqltest(t *testing.T) {
    t.Parallel()
    runSqlTest(t, testEmoji_gosqltest)
}

func testEmoji_gosqltest(t *sqlTest) {
    t.sqlGenInfo = &sqlGenInfo{}
    t.tableName = tablePrefix + "emoji"
    t.columnNameType = [][2]string{
        {"id", "integer primary key"},
        {"c1", "varchar(10)"},
        {"c2", "varchar(10)"},
    }
    t.genTableTest()

    wantC1 := "😁"
    wantC2 := "😮"
    t.mustExec(fmt.Sprintf("insert into %s (id,c1,c2) values(?, ?, ?)", t.tableName), 1, wantC1, wantC2)

    gotC1 := ""
    gotC2 := ""
    err := t.QueryRow(fmt.Sprintf("select c1,c2 from %s where id = ?", t.tableName), 1).Scan(&gotC1, &gotC2)
    if err != nil {
        t.Errorf("[]byte scan: %v", err)
    } else if gotC1 != wantC1 || gotC2 != wantC2 {
        t.Errorf("for emoji, got %s, %s; want %s, %s", gotC1, gotC2, wantC1, wantC2)
    }
}
