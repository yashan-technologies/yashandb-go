package yasdb

import "testing"

func TestRmCommnetAndlSemicolon(t *testing.T) {
	testDatas := []struct {
		caseName  string
		plsql     string
		wantPlsql string
	}{
		{"case01", `/*创建存储过程*/
        create or replace procedure myl1.procAdd2(p1 in int, p2 in int) IS
         v3 int;
        begin
         v3 := p1 + p2;
        end;
        /`, "create or replace procedure myl1.procAdd2(p1 in int, p2 in int) IS\n         v3 int;\n        begin\n         v3 := p1 + p2;\n        end;\n        /"},
		{"case02", `-- 创建存储过程
        create or replace procedure myl1.procAdd2(p1 in int, p2 in int) IS
         v3 int;
        begin
         v3 := p1 + p2;
        end;
        /`, "create or replace procedure myl1.procAdd2(p1 in int, p2 in int) IS\n         v3 int;\n        begin\n         v3 := p1 + p2;\n        end;\n        /"},
	}

	for _, v := range testDatas {
		t.Run(v.caseName, func(t *testing.T) {
			plsql := rmCommnetAndlSemicolon(v.plsql)
			if plsql != v.wantPlsql {
				t.Fatalf("rm psql failed, want: %s; actual: %s", v.wantPlsql, plsql)
			}
		})

	}
}
