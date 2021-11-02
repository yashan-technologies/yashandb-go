package yasdb_test

import (
    "database/sql"
    "log"
    "testing"
    "time"

    _ "yasdb"
)

func init() {
    log.SetFlags(log.Ldate | log.Llongfile)
}

func Connect() *sql.DB {
    db, err := sql.Open("yasdb", "sys/sys@192.168.30.219:16001")
    if err != nil {
        log.Fatalf("some error %s", err.Error())
    }
    return db
}

type DataT struct {
    Id      int
    Name    string
    Created time.Time
}

type Database struct {
    Status string
    Role   string
    Point  string
}

func TestQueryString(t *testing.T) {
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

func TestDate(t *testing.T) {
    db := Connect()
    tx, err := db.Begin()
    if err != nil {
        t.Fatal(err)
    }
    _, err = tx.Exec("drop table if exists go_test_t2")
    if err != nil {
        t.Fatal(err)
    }
    _, err = tx.Exec("create table go_test_t2 (id int, name varchar(20), created timestamp)")
    if err != nil {
        t.Fatal(err)
    }
    res, err := tx.Exec("insert into go_test_t2 values (1, 'oldman', '2021-11-02 01:00:00.0')")
    if err != nil {
        log.Panic(err)
    }
    affect, _ := res.RowsAffected()
    log.Println(affect, " row affected")
    res, err = tx.Exec("insert into go_test_t2 values (?, ?, ?)", 2, "fireman", time.Now())
    if err != nil {
        t.Fatal(err)
    }
    affect, _ = res.RowsAffected()
    log.Println(affect, " row affected")
    stmt, err := tx.Prepare(`insert into go_test_t2 values(?, ?, ?)`)
    if err != nil {
        t.Fatal(err)
    }
    res, err = stmt.Exec(4, "youngman", time.Now())
    if err != nil {
        t.Fatal(err)
    }
    affect, _ = res.RowsAffected()
    log.Println(affect, " row affected")
    if err := tx.Commit(); err != nil {
        t.Fatal(err)
    }

    rows, err := db.Query("select id, name, created from go_test_t2")
    if err != nil {
        log.Fatal("some wrong for query", err.Error())
    }
    for rows.Next() {
        var date DataT
        if err := rows.Scan(&date.Id, &date.Name, &date.Created); err != nil {
            log.Fatal(err)
        }
        log.Println(date.Name, date.Created.Format("2006-01-02 15:04:05.999999"))
    }
}
