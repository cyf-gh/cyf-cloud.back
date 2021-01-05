package orm

import (
	"../../dm_1"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	InitEngine("L:/cyf-cloud.db/")
	ec := m.Run()
	os.Exit(ec)
}

func TestDMAddResource(t *testing.T) {
	dmRoot := dm_1.DMResource{
		Path: "L:/mount",
	}
	recRes := dmRoot.LsRecruit()

	e := DMAddResource( recRes ); if e != nil { t.Log( e ); t.FailNow() }
}

func TestDMTargetResource_GetClones(t *testing.T) {
	dmRes, e := DMGetTargetResourceById(31 ); if e != nil { t.Log( e ); t.FailNow() }
	rs, e := dmRes.GetClones(); if e != nil { t.Log( e ); t.FailNow() }
	for _, r := range rs {
		t.Log( r.Path )
	}
}

func TestDMGetTargetResourceByPath(t *testing.T) {
	dmRes, e := DMGetTargetResourceByPath("L:/mount/HHD_SIMUL/1/1-2/c.txt" ); if e != nil { t.Log( e ); t.FailNow() }
	t.Log( dmRes )
}

func TestDMIsTargetResourceExist(t *testing.T) {
	exist, e := DMIsTargetResourceExist( "L:/mount/HHD_SIMUL/1/1-2/c.txt" ); if e != nil { t.Log( e ); t.FailNow() }
	t.Log( exist )
}

func TestDMTargetResource_Update(t *testing.T) {
	dmRes, e := DMGetTargetResourceById(31 ); if e != nil { t.Log( e ); t.FailNow() }
	dmRes.Description = "hello from TestDMTargetResource_Update"
	e = dmRes.Update(); if e != nil { t.Log( e ); t.FailNow() }
}

func TestDMIsTagExist(t *testing.T) {
	t.Log( DMIsTagExist( "aaa" ) )
}

func TestDMTag_Insert(t *testing.T) {
	tag := DMTag{
		Name: "aaa",
	}
	if !DMIsTagExist( tag.Name ) {
		tag.Insert()
	}
}
