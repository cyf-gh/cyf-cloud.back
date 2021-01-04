package dm_1

import (
	"bufio"
	"encoding/base64"
	"github.com/nfnt/resize"
	"image"
	"io/ioutil"
	"os"
)

type (
	DMImage struct {
		DMResource
	}
	DMImageFile struct {
		Image image.Image
		Ext string
	}
)

// 检查是否为图片文件
func ( R DMImage ) IsValid() bool {
	if R.IsDire() {
		return false
	}
	ext, _ := R.GetExt()
	if ext != "" {
		for _, ie := range DMExts["image"] {
			if ie == ext {
				return true
			}
		}
	}
	return false
}

func ( R DMImage ) Open() ( If DMImageFile, e error ) {
	var img image.Image
	f, e := os.Open( R.Path )
	defer f.Close()

	ext, e := R.GetExt()
	img, e = DMImageDecode(ext, f)
	if e != nil {
		return
	}
	If = DMImageFile{ Image:img, Ext:ext }
	return
}

func ( R DMImageFile ) ToBase64() ( b64 string, er error ) {
	buf, er := DMImageEncode( R.Ext, R.Image )
	r := bufio.NewReader( buf )

	if c, e := ioutil.ReadAll(r); e != nil {
		er = e
		return
	} else {
		b64 = base64.StdEncoding.EncodeToString( c )
		return
	}
}

func ( R DMImageFile ) Resize( maxW, maxH uint ) DMImageFile {
	return DMImageFile{ Image: resize.Thumbnail( maxW, maxH, R.Image, resize.Lanczos3 ), Ext: R.Ext }
}