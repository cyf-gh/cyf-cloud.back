package dm_1

type (
	DMImage struct {
		DMResource
	}
	DMImageFile struct {

	}
)

// 检查是否为图片文件
func ( R DMImage ) IsValid() {

}

// func ( R DMImage ) Open() ( e error ) {
// 	// ext, e := R.GetExt()
// }
//
// func ( R DMImageFile ) ToBase64() ( b64 string, e error ) {
//
// }
//
// func ( R DMImageFile ) ToFile() ( e error ) {
//
// }