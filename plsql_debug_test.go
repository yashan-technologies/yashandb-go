package yasdb

import (
	"database/sql"
	"fmt"
	"testing"
)

var (
	plsql_1 = `create or replace procedure procAdd(p1 in int, p2 in int, p3 out int) 
	           is 
	           begin
               p3 := 100; 
			   p3 := p1 + p2;
               end;`

	callPlSql_1 = `call procAdd(?,?,?)`
)

func createProcedute(t *testing.T) {
	db, err := sql.Open("yasdb", fmt.Sprintf("%s?%s", testDsn, "autocommit=true"))
	if err != nil {
		t.Fatalf("open database err: %v", err)
		return
	}
	defer db.Close()

	_, err = db.Exec(plsql_1)
	if err != nil {
		t.Fatalf(err.Error())
		return
	}
}

func TestNewPlsqlDebug(t *testing.T) {
	out := int64(0)
	p, err := NewPlsqlDebug(testDsn, callPlSql_1, "procAdd", int64(1), int64(2), sql.Out{Dest: &out})
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
}

func TestPdbgStart(t *testing.T) {
	createProcedute(t)
	out := 0
	v1 := 1
	v2 := 100
	p, err := NewPlsqlDebug(testDsn, callPlSql_1, "procAdd", v1, v2, sql.Out{Dest: &out})
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}
	if err := p.Abort(); err != nil {
		t.Fatal(err)
	}
}

func TestPdbgContinte(t *testing.T) {
	createProcedute(t)
	out := 0
	v1 := 102
	v2 := 100
	p, err := NewPlsqlDebug(testDsn, callPlSql_1, "procAdd", v1, v2, sql.Out{Dest: &out})
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}
	if err := p.Continue(); err != nil {
		t.Fatal(err)
	}
	if err := p.GetBindOutValue(); err != nil {
		t.Fatal(err)
	}
	if out != (v1 + v2) {
		t.Fatalf("bind out value %d != %d", out, v1+v2)
	}
	fmt.Println(out)
}

func TestPdbgAddBreakpoint(t *testing.T) {
	createProcedute(t)
	out := int64(0)
	p, err := NewPlsqlDebug(testDsn, callPlSql_1, "procAdd", int64(1), int64(2), sql.Out{Dest: &out})
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	if err := p.Start(); err != nil {
		t.Fatal(err)
	}
	bpId1, err := p.AddBreakpoint(1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("breakpoint1 is %d\n", bpId1)
	// bpId2, err := p.AddBreakpoint(4)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fmt.Printf("breakpoint2 is %d\n", bpId2)

	if err := p.Continue(); err != nil {
		t.Fatal(err)
	}
	// if err := p.Continue(); err != nil {
	// 	t.Fatal(err)
	// }
}

func TestPdgStepNextStepInto(t *testing.T) {
	createProcedute(t)
	out := int64(0)
	p, err := NewPlsqlDebug(testDsn, callPlSql_1, "procAdd", int64(1), int64(2), sql.Out{Dest: &out})
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	if err := p.Start(); err != nil {
		t.Fatal(err)
	}
	if err := p.StepInto(); err != nil {
		t.Fatal(err)
	}
	if err := p.StepInto(); err != nil {
		t.Fatal(err)
	}
	if err := p.StepInto(); err != nil {
		t.Fatal(err)
	}
}

func TestPdgStepNextStepNext(t *testing.T) {
	createProcedute(t)
	out := int64(0)
	p, err := NewPlsqlDebug(testDsn, callPlSql_1, "procAdd", int64(1), int64(2), sql.Out{Dest: &out})
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	if err := p.Start(); err != nil {
		t.Fatal(err)
	}
	if err := p.StepNext(); err != nil {
		t.Fatal(err)
	}
	if err := p.Abort(); err != nil {
		t.Fatal(err)
	}
}

func TestPdbgStepOut(t *testing.T) {
	createProcedute(t)
	out := int64(0)
	p, err := NewPlsqlDebug(testDsn, callPlSql_1, "procAdd", int64(1), int64(2), sql.Out{Dest: &out})
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	if err := p.Start(); err != nil {
		t.Fatal(err)
	}
	if err := p.StepOut(); err != nil {
		t.Fatal(err)
	}
}
