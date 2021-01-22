package dm_1

import (
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
)

type (
	DMNotifyFile struct {
		Watch *fsnotify.Watcher
		Create, Write, Remove, Rename, Chmod DMWatchEventFunc
	}
	DMWatchEventFunc func( event fsnotify.Event )
)
func DMNewNotifyFileEmpty() *DMNotifyFile {
	w := new(DMNotifyFile)
	w.Watch, _ = fsnotify.NewWatcher()
	w.Create = func(event fsnotify.Event) {}
	w.Remove = func(event fsnotify.Event) {}
	w.Rename = func(event fsnotify.Event) {}
	w.Write = func(event fsnotify.Event) {}
	w.Chmod = func(event fsnotify.Event) {}

	return w
}

// 递归监控目录 dir
func (pR *DMNotifyFile) AddWatchDirRecruit( dir string ) ( e error ) {
	e = filepath.Walk( dir, func( path string, info os.FileInfo, err error ) error {
		if info.IsDir() {
			if path, err = filepath.Abs( path ); err != nil { return err }
				if err = pR.Watch.Add( path ); err != nil { return err }
			Info( "ADD MONITOR: " + path )
		}
		return nil
	})
	return e
}

func ( pR *DMNotifyFile ) RemoveWatchDirRecruit( dir string ) ( e error ) {
	e = filepath.Walk( dir, func( path string, info os.FileInfo, err error ) error {
		if info.IsDir() {
			if path, err = filepath.Abs( path ); err != nil { return err }
			if err = pR.Watch.Remove( path ); err != nil { return err }
			Info( "REMOVE MONITOR: " + path )
		}
		return nil
	})
	return e
}

func (pR *DMNotifyFile) WatchEvent() {
	defer pR.Watch.Close()
	for {
		select {
		case ev := <- pR.Watch.Events: {
				if ev.Op&fsnotify.Create == fsnotify.Create {
					Info("CREATE: "+ ev.Name)
					file, err := os.Stat(ev.Name)
					if err == nil { pR.Create( ev ) }
					if err == nil && file.IsDir() { pR.AddWatchDirRecruit( ev.Name ) }
				}

				if ev.Op&fsnotify.Write == fsnotify.Write {
					Info("WRITE: "+ ev.Name)
					pR.Write( ev )
				}

				if ev.Op&fsnotify.Remove == fsnotify.Remove {
					Info("DELETE: ", ev.Name)
					fi, err := os.Stat(ev.Name)
					if err == nil { pR.Remove( ev ) }
					if err == nil && fi.IsDir() {
						pR.RemoveWatchDirRecruit( ev.Name )
						Info("REMOVE MONITOR: ", ev.Name)
					}
				}

				if ev.Op&fsnotify.Rename == fsnotify.Rename {
					Info("RENAME: ", ev.Name)
					pR.Rename( ev )
					// 名字更换后的文件夹会被认为是CREATE的
					pR.Watch.Remove( ev.Name )
				}

				if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
					{ pR.Chmod( ev ) }
					Info("CHMOD: ", ev.Name)
				}
			}
		case err := <- pR.Watch.Errors:
			{
				Fatal("ERROR IN WATCH EVENT: ", err)
				return
			}
		}
	}
}
