package yasdb

import (
	"database/sql"
	"testing"
)

func TestNull(t *testing.T) {

	db, err := sql.Open("yasdb", testDsn)
	if err != nil {
		t.Fatalf("open database err: %v", err)
		return
	}
	defer db.Close()

	_, err = db.Exec("drop table if exists test_null;")
	if err != nil {
		t.Fatalf(err.Error())
		return
	}

	_, err = db.Exec("create table test_null(id int,c1 int,c2 boolean,c3 date,c4 float) ;")
	if err != nil {
		t.Fatalf(err.Error())
		return
	}

	_, err = db.Exec("insert into test_null(id) values(?);  ", 1)
	if err != nil {
		t.Fatalf(err.Error())
		return
	}

	rows, err := db.Query("select c1, c2, c3, c4 from test_null;")
	if err != nil {
		t.Fatalf(err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			c1 sql.NullInt32
			c2 sql.NullBool
			c3 sql.NullTime
			c4 sql.NullFloat64
		)

		if err := rows.Scan(&c1, &c2, &c3, &c4); err != nil {
			t.Fatalf(err.Error())
			return
		}
		if c1.Valid != false || c2.Valid != false || c3.Valid != false || c4.Valid != false {
			t.Fatalf("null doesn't work correctly")
			return
		}
	}

}
