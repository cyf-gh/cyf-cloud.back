package dm_1

import "testing"

func TestDMResource_GetExt(t *testing.T) {
	dm := DMResource{Path:"/a/b/c/1.png"}
	if ext, e := dm.GetExt(); e != nil {
		t.Log( e ); t.FailNow()
	} else {
		t.Log( ext )
	}
}

func TestDMResource_Ls(t *testing.T) {
	dm := DMResource{Path:"C:/.repos"}
	if rs, e := dm.Ls(); e != nil {
		t.Log( e ); t.FailNow()
	} else {
		t.Log( rs )
	}
}

func TestDMResource_IsDir(t *testing.T) {
	dm := DMResource{Path:"C:/.repos"}
	if dm.IsDire() {
		t.Log("is dir")
	} else {
		t.FailNow()
	}
}