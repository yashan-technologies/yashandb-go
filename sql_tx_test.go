package yasdb

import "testing"

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
