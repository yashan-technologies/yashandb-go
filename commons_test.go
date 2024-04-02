package yasdb

import "testing"

func TestRmCommnet(t *testing.T) {
	wantPlsql := "create or replace procedure myl1.procAdd2(p1 in int, p2 in int) IS\nv3 int;\nbegin\nv3 := p1 + p2;\nend;\n/\n"
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
        /`, wantPlsql},
		{"case02", `-- 创建存储过程
        create or replace procedure myl1.procAdd2(p1 in int, p2 in int) IS
         v3 int;
        begin
         v3 := p1 + p2;
        end;
        /`, wantPlsql},
		{"case03", `
		/*创建
		存储过
		程*/
		
		/*创*建
		存储*过
		程*/
	        create or replace procedure myl1.procAdd2(p1 in int, p2 in int) IS
         v3 int;
		 /*创*建
		存储*过
		程*/
        begin
         v3 := p1 + p2;
        end;
        /`, wantPlsql},
		{"case04", `
		-- 创建存储过程
		/*创*建
		存储*过
		程*/
        create or replace procedure myl1.procAdd2(p1 in int, p2 in int) IS
         v3 int;
        begin
         v3 := p1 + p2;
        end;
        /`, wantPlsql},
	}

	for _, v := range testDatas {
		t.Run(v.caseName, func(t *testing.T) {
			plsql := rmComment(v.plsql)
			if plsql != v.wantPlsql {
				t.Fatalf("rm psql failed, want: %s; actual: %s", v.wantPlsql, plsql)
			}
		})

	}
}
