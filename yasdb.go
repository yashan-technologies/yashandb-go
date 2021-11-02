package yasdb

import (
    "database/sql"
    "log"
)

func init() {
    log.Println("register yasdb driver")
    sql.Register("yasdb", &YasDriver{})
}
