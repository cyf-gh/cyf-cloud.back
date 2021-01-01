package dm_1

import "testing"

func TestDMImageFile_ToBase64(t *testing.T) {
	dm := DMImage{ DMResource{Path:"C:/Users/cyf-m/Pictures/go.gif"}}
	img, e := dm.Open()
	if e != nil {
		t.FailNow()
	}
	b64, e := img.ToBase64()
	if e != nil {
		t.FailNow()
	} else {
		t.Log( b64 )
	}
}