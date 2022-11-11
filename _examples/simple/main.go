package main

import (
    "database/sql"
    "fmt"

    _ "git.yasdb.com/go/yasdb-go"
    example "git.yasdb.com/go/yasdb-go/_examples"
)

func getYasdbConn(dsn string) (*sql.DB, error) {
    return sql.Open("yasdb", dsn)
}

func main() {
    dsn := example.GetDsn()
    db, err := sql.Open("yasdb", dsn)
    if err != nil {
        fmt.Println("failed to connect yashandb, err:", err)
        return
    }
    defer db.Close()
    rows, err := db.Query("select version from v$instance")
    if err != nil {
        fmt.Println(err)
        return
    }
    for rows.Next() {
        var version string
        err = rows.Scan(&version)
        if err != nil {
            fmt.Println(err)
            return
        }
        fmt.Println("YashanDB version:", version)
    }
}
