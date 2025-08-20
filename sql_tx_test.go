package yasdb

import (
	"testing"
)

func TestTxCommit(t *testing.T) {
	db := newSqlTest(t)
	defer db.Close()
	db.mustExec("drop table if exists test_tx_commit")
	db.mustExec("create table test_tx_commit(c1 int,c2 varchar(20))")

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("begin tx err: %v", err)
	}
	_, err = tx.Exec("insert into test_tx_commit values(1,'000001')")
	if err != nil {
		t.Fatalf("exec tx err: %v", err)
	}
	_, err = tx.Exec("insert into test_tx_commit values(2,'000001')")
	if err != nil {
		t.Fatalf("exec tx err: %v", err)
	}
	_, err = tx.Query("select * from test_tx_commit")
	if err != nil {
		t.Fatalf("exec tx err: %v", err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("commit tx err: %v", err)
	}
}

func TestTxRollback(t *testing.T) {
	db := newSqlTest(t)
	defer db.Close()
	db.mustExec("drop table if exists test_tx_rollback")
	db.mustExec("create table test_tx_rollback(c1 int,c2 varchar(20))")

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("begin tx err: %v", err)
	}
	_, err = tx.Exec("insert into test_tx_rollback values(1,'000001')")
	if err != nil {
		t.Fatalf("exec tx err: %v", err)
	}
	_, err = tx.Exec("insert into test_tx_rollback values(2,'000001')")
	if err != nil {
		t.Fatalf("exec tx err: %v", err)
	}
	_, err = tx.Query("select * from test_tx_rollback")
	if err != nil {
		t.Fatalf("exec tx err: %v", err)
	}
	if err := tx.Rollback(); err != nil {
		t.Fatalf("rollback tx err: %v", err)
	}
}

func TestTxAutoCommit(t *testing.T) {
	// set autocommit true
	db := newSqlAutoCommitTest(t)
	defer db.Close()
	db.mustExec("drop table if exists t1")
	db.mustExec("create table t1(c1 int)")

	// start tx
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("begin tx err: %v", err)
	}
	_, err = tx.Exec("insert into t1 values(1)")
	if err != nil {
		t.Fatalf("exec tx err: %v", err)
	}
	_, err = tx.Exec("insert into t1 values(2)")
	if err != nil {
		t.Fatalf("exec tx err: %v", err)
	}
	_, err = tx.Query("select * from t1")
	if err != nil {
		t.Fatalf("exec tx err: %v", err)
	}
	if err := tx.Rollback(); err != nil {
		t.Fatalf("rollback tx err: %v", err)
	}
	var count int
	err = db.QueryRow("select count(*) from t1").Scan(&count)
	if err != nil {
		t.Fatalf("query after rollback tx err: %v", err)
	}
	if count != 0 {
		t.Errorf("after rollback expect count=0,in fact count=%d", count)
	}

}

func TestTxAutoCommit2(t *testing.T) {
	// set autocommit true
	db := newSqlAutoCommitTest(t)
	defer db.Close()
	db.mustExec("drop table if exists t1")
	db.mustExec("create table t1(c1 int)")

	// start tx
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("begin tx err: %v", err)
	}
	_, err = tx.Exec("insert into t1 values(1)")
	if err != nil {
		t.Fatalf("exec tx err: %v", err)
	}
	_, err = tx.Exec("insert into t1 values(2)")
	if err != nil {
		t.Fatalf("exec tx err: %v", err)
	}
	_, err = tx.Query("select * from t1")
	if err != nil {
		t.Fatalf("exec tx err: %v", err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("commit tx err: %v", err)
	}
	var count int
	err = db.QueryRow("select count(*) from t1").Scan(&count)
	if err != nil {
		t.Fatalf("query after commit tx err: %v", err)
	}
	if count != 2 {
		t.Errorf("after commit expect count=2,in fact count=%d", count)
	}

	// test auto commit open now
	_, err = db.Exec("insert into t1 values (3)")
	if err != nil {
		t.Fatalf("query err: %v", err)
	}

	db2 := newSqlAutoCommitTest(t)
	err = db2.QueryRow("select count(*) from t1").Scan(&count)
	if err != nil {
		t.Fatalf("query after commit tx err: %v", err)
	}
	if count != 3 {
		t.Errorf("after commit expect count=3,in fact count=%d", count)
	}
}
