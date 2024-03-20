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
	// select OBJECT_ID,SUBPROGRAM_ID from dba_procedures where  object_name = 'PROCADD';
	// select VERSION from DBA_SOURCE where OWNER = UPPER(?) AND NAME = UPPER(?)
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
	p, err := NewPlsqlDebug(testDsn, callPlSql_1, int64(1), int64(2), sql.Out{Dest: &out})
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
	p, err := NewPlsqlDebug(testDsn, callPlSql_1, v1, v2, sql.Out{Dest: &out})
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
	p, err := NewPlsqlDebug(testDsn, callPlSql_1, v1, v2, sql.Out{Dest: &out})
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

func TestPdgStepNextStepInto(t *testing.T) {
	createProcedute(t)
	out := int64(0)
	p, err := NewPlsqlDebug(testDsn, callPlSql_1, int64(1), int64(2), sql.Out{Dest: &out})
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

func TestPdgStepNext(t *testing.T) {
	createProcedute(t)
	out := int64(0)
	p, err := NewPlsqlDebug(testDsn, callPlSql_1, int64(1), int64(2), sql.Out{Dest: &out})
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
	p, err := NewPlsqlDebug(testDsn, callPlSql_1, int64(1), int64(2), sql.Out{Dest: &out})
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

func TestPdbgGetRunningData(t *testing.T) {
	createProcedute(t)
	out := 0
	v1 := 102
	v2 := 100
	p, err := NewPlsqlDebug(testDsn, callPlSql_1, v1, v2, sql.Out{Dest: &out})
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	var status DebuggerStatus
	if err := PdbgGetRunningData(p.Stmt, DBG_RUNNING_ATTR_STATUS, &status); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("running debugger status: %v\n", status)

	var objId uint64
	if err := PdbgGetRunningData(p.Stmt, DBG_RUNNING_ATTR_OBJ_ID, &objId); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("running obj id: %v\n", objId)

	var className string
	if err := PdbgGetRunningData(p.Stmt, DBG_RUNNING_ATTR_CLASS_NAME, &className); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("running class name: %v\n", className)

	var methodName string
	if err := PdbgGetRunningData(p.Stmt, DBG_RUNNING_ATTR_METHOD_NAME, &methodName); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("running method name: %v\n", methodName)

	if err := p.Continue(); err != nil {
		t.Fatal(err)
	}

	if err := PdbgGetRunningData(p.Stmt, DBG_RUNNING_ATTR_STATUS, &status); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("running debugger status: %v\n", status)

	if err := p.GetBindOutValue(); err != nil {
		t.Fatal(err)
	}
	if out != (v1 + v2) {
		t.Fatalf("bind out value %d != %d", out, v1+v2)
	}
	fmt.Println(out)
}

func TestPdbgGetFrameData(t *testing.T) {
	createProcedute(t)
	out := 0
	v1 := 102
	v2 := 100
	p, err := NewPlsqlDebug(testDsn, callPlSql_1, v1, v2, sql.Out{Dest: &out})
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	if err := PdbgStepInto(p.Stmt); err != nil {
		t.Fatal(err)
	}

	count, err := PdbgGetAllFrames(p.Stmt)
	if err != nil {
		t.Fatal(err)
	}

	for i := uint32(0); i < count; i++ {
		var objId uint64
		if err := PdbgGetFrameData(p.Stmt, i, DBG_FRAME_ATTR_OBJ_ID, &objId); err != nil {
			t.Fatal(err)
		}
		fmt.Printf("frame: %d ,object id: %d\n", i, objId)

		var subId uint16
		if err := PdbgGetFrameData(p.Stmt, i, DBG_FRAME_ATTR_SUB_ID, &subId); err != nil {
			t.Fatal(err)
		}
		fmt.Printf("frame: %d ,subid: %d\n", i, subId)

		var lineNo uint32
		if err := PdbgGetFrameData(p.Stmt, i, DBG_FRAME_ATTR_LINE_NO, &lineNo); err != nil {
			t.Fatal(err)
		}
		fmt.Printf("frame: %d ,lineNo: %d\n", i, subId)

		var stackNo uint32
		if err := PdbgGetFrameData(p.Stmt, i, DBG_FRAME_ATTR_STACK_NO, &stackNo); err != nil {
			t.Fatal(err)
		}
		fmt.Printf("frame: %d ,stackNo: %d\n", i, subId)

		var className string
		if err := PdbgGetFrameData(p.Stmt, i, DBG_FRAME_ATTR_CLASS_NAME, &className); err != nil {
			t.Fatal(err)
		}
		fmt.Printf("frame: %d ,className: %s\n", i, className)

		var methodName string
		if err := PdbgGetFrameData(p.Stmt, i, DBG_FRAME_ATTR_METHOD_NAME, &methodName); err != nil {
			t.Fatal(err)
		}
		fmt.Printf("frame: %d ,methodName: %s\n", i, methodName)

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

func TestPdbgGetAllVars(t *testing.T) {
	createProcedute(t)
	out := 0
	v1 := 102
	v2 := 100
	p, err := NewPlsqlDebug(testDsn, callPlSql_1, v1, v2, sql.Out{Dest: &out})
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}
	varCount, err := p.GetAllVars()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("values count is %v\n", varCount)
}

func TestPdbgGetAllFrames(t *testing.T) {
	createProcedute(t)
	out := 0
	v1 := 102
	v2 := 100
	p, err := NewPlsqlDebug(testDsn, callPlSql_1, v1, v2, sql.Out{Dest: &out})
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}
	varCount, err := p.GetAllFrames()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("frame count is %v\n", varCount)
}

func TestPdbgAddBreakpoint(t *testing.T) {
	createProcedute(t)
	out := int64(0)
	p, err := NewPlsqlDebug(testDsn, callPlSql_1, int64(1), int64(2), sql.Out{Dest: &out})
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	bpId, err := p.AddBreakpoint(2361, 1, 4) // select OBJECT_ID,SUBPROGRAM_ID from dba_procedures where  object_name = 'PROCADD';
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("breakpoint id is %v\n", bpId)

	bpId, err = p.AddBreakpoint(2361, 1, 5)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("breakpoint id is %v\n", bpId)

	if err := p.Continue(); err != nil {
		t.Fatal(err)
	}

	if err := p.Continue(); err != nil {
		t.Fatal(err)
	}
}

func TestPdbgGetBreakpointData(t *testing.T) {
	createProcedute(t)
	out := int64(0)
	p, err := NewPlsqlDebug(testDsn, callPlSql_1, int64(1), int64(2), sql.Out{Dest: &out})
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	bpId, err := p.AddBreakpoint(2361, 1, 4) // select OBJECT_ID,SUBPROGRAM_ID from dba_procedures where  object_name = 'PROCADD';
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("breakpoint id is %v\n", bpId)

	bpId, err = p.AddBreakpoint(2361, 1, 5)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("breakpoint id is %v\n", bpId)

	bpCount, err := p.GetBreakpointsCount()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("breakpoint count is %v\n", bpCount)

	for i := uint32(0); i < bpCount; i++ {
		var objId uint64
		if err := p.GetBreakpointData(i, DBG_BP_ATTR_OBJ_ID, &objId); err != nil {
			t.Fatal(err)
		}
		fmt.Printf("breakpoit %d objectId is %v\n", i, objId)

		var subId uint16
		if err := p.GetBreakpointData(i, DBG_BP_ATTR_SUB_ID, &subId); err != nil {
			t.Fatal(err)
		}
		fmt.Printf("breakpoit %d subId is %v\n", i, subId)

		var lineNo uint32
		if err := p.GetBreakpointData(i, DBG_BP_ATTR_LINE_NO, &lineNo); err != nil {
			t.Fatal(err)
		}
		fmt.Printf("breakpoit %d lineNo is %v\n", i, lineNo)

		var bpId uint32
		if err := p.GetBreakpointData(i, DBG_BP_ATTR_BP_ID, &bpId); err != nil {
			t.Fatal(err)
		}
		fmt.Printf("breakpoit %d bpId is %v\n", i, bpId)

	}
}

func TestPdbgDeleteBrakPoint(t *testing.T) {
	createProcedute(t)
	out := int64(0)
	p, err := NewPlsqlDebug(testDsn, callPlSql_1, int64(1), int64(2), sql.Out{Dest: &out})
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	bpId, err := p.AddBreakpoint(2361, 1, 4)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("breakpoint id is %v\n", bpId)

	bpId, err = p.AddBreakpoint(2361, 1, 5)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("breakpoint id is %v\n", bpId)

	bpCount, err := p.GetBreakpointsCount()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("breakpoint count is %v\n", bpCount)

	i := uint32(1)
	var objId uint64
	if err := p.GetBreakpointData(i, DBG_BP_ATTR_OBJ_ID, &objId); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("breakpoit %d objectId is %v\n", i, objId)

	var subId uint16
	if err := p.GetBreakpointData(i, DBG_BP_ATTR_SUB_ID, &subId); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("breakpoit %d subId is %v\n", i, subId)

	var lineNo uint32
	if err := p.GetBreakpointData(i, DBG_BP_ATTR_LINE_NO, &lineNo); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("breakpoit %d lineNo is %v\n", i, lineNo)

	if err := p.DeleteBreakpoint(objId, subId, lineNo); err != nil {
		t.Fatal(err)
	}

	bpCount, err = p.GetBreakpointsCount()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("breakpoint count is %v\n", bpCount)
}

func TestPdbgDeleteAllBreakpoints(t *testing.T) {
	createProcedute(t)
	out := int64(0)
	p, err := NewPlsqlDebug(testDsn, callPlSql_1, int64(1), int64(2), sql.Out{Dest: &out})
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	bpId, err := p.AddBreakpoint(2361, 1, 4)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("breakpoint id is %v\n", bpId)

	bpId, err = p.AddBreakpoint(2361, 1, 5)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("breakpoint id is %v\n", bpId)

	bpCount, err := p.GetBreakpointsCount()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("after delete all breakpoint count is %v\n", bpCount)

	if err := p.DeleteAllBreakpoint(); err != nil {
		t.Fatal(err)
	}

	bpCount, err = p.GetBreakpointsCount()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("breakpoint count is %v\n", bpCount)
}

func TestPdbgGetVarData(t *testing.T) {
	createProcedute(t)
	out := 0
	v1 := 102
	v2 := 100
	p, err := NewPlsqlDebug(testDsn, callPlSql_1, v1, v2, sql.Out{Dest: &out})
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}
	varCount, err := p.GetAllVars()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("values count is %v\n", varCount)

	for i := uint32(0); i < varCount; i++ {
		var blockNo uint32
		if err := p.GetVarData(i, DBG_VAR_ATTR_BLOCK_NO, &blockNo); err != nil {
			t.Fatal(err)
		}
		fmt.Printf("var %d blockNo is %v\n", i, blockNo)

		var dataType uint8
		if err := p.GetVarData(i, DBG_VAR_ATTR_TYPE, &dataType); err != nil {
			t.Fatal(err)
		}
		fmt.Printf("var %d dataType number is %v; dataType name is %s\n", i, dataType, GetDatabaseTypeName(uint32(dataType)))

		var isGlobal bool
		if err := p.GetVarData(i, DBG_VAR_ATTR_IS_GLOBAL, &isGlobal); err != nil {
			t.Fatal(err)
		}
		fmt.Printf("var %d isGlobal is %v\n", i, isGlobal)

		var name string
		if err := p.GetVarData(i, DBG_VAR_ATTR_NAME, &name); err != nil {
			t.Fatal(err)
		}
		fmt.Printf("var %d name is %v\n", i, name)

		var valueSize uint32
		if err := p.GetVarData(i, DBG_VAR_ATTR_VALUE_SIZE, &valueSize); err != nil {
			t.Fatal(err)
		}
		fmt.Printf("var %d valueSize is %v\n", i, valueSize)

	}
}

func TestPdbgGetVarValue(t *testing.T) {
	createProcedute(t)
	out := 0
	v1 := 102
	v2 := 100
	p, err := NewPlsqlDebug(testDsn, callPlSql_1, v1, v2, sql.Out{Dest: &out})
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}
	varCount, err := p.GetAllVars()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("values count is %v\n", varCount)

	for i := uint32(0); i < varCount; i++ {
		value, err := PdbgGetVarValue(p.Stmt, i)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("var %d value is %v\n", i, value)

	}
}

func TestPdbgGetAllData(t *testing.T) {
	createProcedute(t)
	out := 0
	v1 := 102
	v2 := 100
	p, err := NewPlsqlDebug(testDsn, callPlSql_1, v1, v2, sql.Out{Dest: &out})
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	bpId, err := p.AddBreakpoint(2361, 1, 4)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("breakpoint id is %v\n", bpId)

	bpId, err = p.AddBreakpoint(2361, 1, 5)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("breakpoint id is %v\n", bpId)

	runningData, err := p.GetAllRunningData()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(*runningData)

	varData, err := p.GetAllVarData()
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range varData {
		fmt.Println(*v)
	}

	frameData, err := p.GetAllFrameData()
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range frameData {
		fmt.Println(*v)
	}

	bPData, err := p.GetAllBreakpointData()
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range bPData {
		fmt.Println(*v)
	}
}
