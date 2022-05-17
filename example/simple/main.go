package main

import (
    "database/sql"
    "log"

    _ "git.yasdb.com/cod-noah/yasdb-go"
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
