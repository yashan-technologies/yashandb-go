package main

import (
    "database/sql"
    "fmt"

    _ "git.yasdb.com/cod-noah/yasdb-go"
    "git.yasdb.com/cod-noah/yasdb-go/example"
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
