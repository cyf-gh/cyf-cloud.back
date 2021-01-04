package dm_1

import "testing"

func TestDMRecruitLs(t *testing.T) {
	r := DMResource{
		Path:       "L:/mount",
	}
	res := DMRecruitLs( r )
	for _, r := range res {
		t.Log( r.Path, "\t"+r.GetGenre() )
	}
}
