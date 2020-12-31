// 进行一些基本的文件操作
// 参考：
// https://colobu.com/2016/10/12/go-file-operations/
package dm_1

import (
	"crypto/md5"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"time"
)

type (
	DMResource struct {
		// 一个资源的绝对路径
		Path string
	}
	// 详见os包中的 types.go
	// 用于方便转化为json对象
	DMFileInfo struct {
		Name string       // base name of the file
		Size int64        // length in bytes for regular files; system-dependent for others
		Mode os.FileMode     // file mode bits
		ModTime time.Time // modification time
		IsDir bool        // abbreviation for Mode().IsDir()
		Sys interface{}   // underlying data source (can return nil)
	}
)

// 获取文件的基本信息
func ( R DMResource ) GetBasicFileInfo() ( dfi *DMFileInfo, e error ) {
	var fi os.FileInfo
	if fi, e = os.Stat(R.Path); e != nil {
		return
	} else {
		dfi = &DMFileInfo{
			Name:    fi.Name(),
			Size:    fi.Size(),
			Mode:    fi.Mode(),
			ModTime: fi.ModTime(),
			IsDir:   fi.IsDir(),
			Sys:     fi.Sys(),
		}
		return
	}
}

func ( R DMResource ) IsDir() bool {
	s, err := os.Stat( R.Path )
	if err != nil {
		return false
	}
	return s.IsDir()
}

func ( R DMResource ) IsFile() bool {
	return !R.IsDir()
}

// 获取某个文件的md5
func ( R DMResource ) GetMD5() ( md5str string, e error ) {
	e = nil
	tMd5 := md5.New()

	if R.IsDir() {
		e = errors.New("cannot get md5 from a directory")
		return
	} else {
		var f *os.File
		if f, e = os.Open( R.Path ); e != nil {
			return
		} else {
			io.Copy( tMd5, f )
			md5str = string( tMd5.Sum([]byte("")) )
			return
		}
	}
}

// 枚举所有的子文件夹
func ( R DMResource ) Ls() ( rs []DMResource, e error ) {
	rs = []DMResource{}
	var fs []os.FileInfo

	if R.IsFile() {
		e = errors.New("cannot ls a file")
		return
	}

	if fs, e = ioutil.ReadDir( R.Path ); e != nil {
		return
	} else {
		for _, f := range fs {
			rs = append(rs, DMResource{
				Path: R.Path +"/"+ f.Name(),
			})
		}
		return
	}
}

// 获取文件扩展名
// 返回形如 ".ext"
func ( R DMResource ) GetExt() ( ext string, e error ) {
	if R.IsFile() {
		ext = path.Ext( R.Path )
		return
	} else {
		e = errors.New("cannot get ext from a directory")
		return
	}
}
