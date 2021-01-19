// 进行一些基本的文件操作
// 参考：
// https://colobu.com/2016/10/12/go-file-operations/
package dm_1

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"
)

type (
	DMResource struct {
		DMFileInfo
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
		IsDir bool        // abbreviation for Mode().IsDire()
		Sys interface{}   // underlying data source (can return nil)
	}
	// 详见os包中的 types.go
	// 用于方便转化为json对象
	DMFileInfoViewModel struct {
		Name, Path string
		Size int64
		Mode string
		ModTime string
		IsDir bool
		Sys interface{}
	}
)

func ( R DMResource ) ToReadable() *DMFileInfoViewModel {
	return &DMFileInfoViewModel{
		Path: 	 R.Path,
		Name:    R.Name,
		Size:    R.Size,
		Mode:    R.Mode.String(),
		ModTime: R.ModTime.Format("2006-01-02 15:04:05"),
		IsDir:   R.IsDir,
		Sys:     R.Sys,
	}
}

// 获取文件的基本信息
func ( pR *DMResource ) GetBasicFileInfo() ( dfi *DMFileInfo, e error ) {
	var fi os.FileInfo
	if fi, e = os.Stat( pR.Path); e != nil {
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
		pR.Name = dfi.Name
		pR.Size = dfi.Size
		pR.Mode = dfi.Mode
		pR.ModTime = dfi.ModTime
		pR.IsDir = dfi.IsDir
		pR.Sys = dfi.Sys
		return
	}
}

func ( R DMResource ) Exists() bool {
	_, err := os.Stat( R.Path )
	if err != nil {
		if os.IsExist( err ) {
			return true
		}
		return false
	}
	return true
}


// 与IsDir区分
func ( R DMResource ) IsDire() bool {
	s, err := os.Stat( R.Path )
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 递归全部的子项目个数
func ( R DMResource ) LsRecruitCount() int {
	return len( R.LsRecruit( nil ) )
}

// 如果目录下有exe文件，说明这个目录为二进制软件的目录
// 判断方式1
func ( R DMResource ) IsBinaryDirectory1() bool {
	if !R.IsDire() {
		return false
	}
	if rs, e := R.Ls(); e != nil {
		return false
	} else {
		for _, r := range rs {
			if r.IsFile() {
				if ext, _ := r.GetExt(); ext == ".exe" {
					return true
				}
			}
		}
		return false
	}
}

func ( R DMResource ) IsFile() bool {
	return !R.IsDire()
}

// 获取某个文件的md5
func ( R DMResource ) GetMD5() ( md5str string, e error ) {
	e = nil
	tMd5 := md5.New()
	md5str = ""

	if R.IsDire() {
		// e = errors.New("cannot get md5 from a directory")
		return "", nil
	} else {
		var f *os.File
		if f, e = os.Open( R.Path ); e != nil {
			return
		} else {
			io.Copy( tMd5, f )
			md5str = hex.EncodeToString(tMd5.Sum(nil) )
			return
		}
	}
}

// 枚举所有的子文件夹
// ls会自动调用所有子资源的GetBasicFileInfo进行填充
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
		for i, f := range fs {
			rs = append( rs, DMResource{
				Path: R.Path +"/"+ f.Name(),
			} )
			_, e = rs[i].GetBasicFileInfo()
		}
		return
	}
}

// 获取文件扩展名
// 返回形如 ".ext"
// e != nil 时表明路径非法或为一个目录
func ( R DMResource ) GetExt() ( ext string, e error ) {
	if R.IsFile() {
		ext = path.Ext( R.Path )
		return
	} else {
		e = errors.New("cannot get ext from a directory")
		return
	}
}

// dm_1.DMExts
func ( R DMResource ) GetGenre() string {
	if R.IsDire() {
		if R.IsBinaryDirectory1() {
			return "binary"
		}
		return "directory"
	} else {
		ext, _ := R.GetExt()

		for k, v := range DMExts {
			for _, e := range v {
				if ext == e {
					return k
				}
			}
		}
		return "binary"
	}
}

// 包含 R 自身
func ( R DMResource ) LsRecruit( status *DMTaskStatus ) []DMResource {
	if status != nil {
		status.ProgressStage = "walking directories..."
	}
	return append( DMRecruitLs( R, status), R )
}

func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

// 获取文件或文件夹的大小
// 与属性Size不同，当接收者为路径时GetSize可以返回文件夹的递归大小
func ( R DMResource ) GetSize() ( int64, error ) {
	if R.IsDire() {
		return dirSize( R.Path )
	} else {
		return R.Size, nil
	}
}