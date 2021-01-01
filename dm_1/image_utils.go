package dm_1

import (
	"bytes"
	"golang.org/x/image/bmp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
)

func DMImageDecode( ext string, f *os.File ) ( img image.Image, e error ){
	switch ext {
	case ".png":
		img, e = png.Decode( f )
		break
	case ".jpeg":
		img, e = jpeg.Decode( f )
		break
	case ".jpg":
		img, e = jpeg.Decode( f )
		break
	case ".bmp":
		img, e = bmp.Decode( f )
		break
	case ".gif":
		img, e = gif.Decode( f )
		break
	}
	return
}

// 有损格式均为默认选项
func DMImageEncode( ext string, img image.Image ) ( buf *bytes.Buffer, e error ) {
	buf = new(bytes.Buffer)

	switch ext {
	case ".png":
		e = png.Encode( buf, img )
		break
	case ".jpeg":
		e = jpeg.Encode( buf, img, nil )
		break
	case ".jpg":
		e = jpeg.Encode( buf, img, nil )
		break
	case ".bmp":
		e = bmp.Encode( buf, img )
		break
	case ".gif":
		e = gif.Encode( buf, img, nil )
		break
	}
	return
}

