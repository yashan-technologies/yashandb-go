package yasdb

import (
    "database/sql"
)

func init() {
    sql.Register("yasdb", &YasdbDriver{})
}
