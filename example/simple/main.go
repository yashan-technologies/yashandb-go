package main

import (
    "database/sql"
    "fmt"
    "log"
    "time"

    _ "git.yasdb.com/cod-noah/yasdb-go"
)

func Connect() *sql.DB {
    db, err := sql.Open("yasdb", "sys/yasdb_123@192.168.31.139:1688")
    if err != nil {
        log.Fatalf("some error %s", err.Error())
    }
    return db
}

var sql_1_row string = `select version from v$instance`

func test2() {
    db := Connect()
    a := ""
    // var err error

    for {
        rows, err := db.Query(sql_1_row)

        // time.Sleep(300 * time.Microsecond)

        for rows.Next() {
            err = rows.Scan(&a)
            if err != nil {
                log.Fatal("some wrong for query", err.Error())
            }
            fmt.Println(a)
            if a != "Release 22.1.B105 x86_64 178549b" {
                panic("no equal")
            }

        }
        time.Sleep(10 * time.Microsecond)
    }
}

func main() {

    test2()
}
