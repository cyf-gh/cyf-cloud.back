package dm_1

import (
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
)

type NotifyFile struct {
		watch *fsnotify.Watcher
}

func NewNotifyFile() *NotifyFile {
	w := new(NotifyFile)
	w.watch, _ = fsnotify.NewWatcher()
	return w
}

// 将会递归监控目录 dir
func (pR *NotifyFile) WatchDir( dir string ) ( e error ) {
	e = filepath.Walk( dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			var path string
			if path, err = filepath.Abs(path); err != nil { return err }
			if err = pR.watch.Add(path); err != nil { return err }
			Info("ADD MONITOR: " + path)
		}
		return nil
	})
	go pR.WatchEvent()
	return nil
}

func (pR *NotifyFile) WatchEvent() {
	for {
		select {
		case ev := <-pR.watch.Events: {
				if ev.Op&fsnotify.Create == fsnotify.Create {
					Info("CREATE: "+ ev.Name)
					file, err := os.Stat(ev.Name)
					if err == nil && file.IsDir() { pR.watch.Add(ev.Name) }
				}

				// TODO: modify md5
				if ev.Op&fsnotify.Write == fsnotify.Write {
				}

				// TODO:
				if ev.Op&fsnotify.Remove == fsnotify.Remove {
					Info("DELETE: ", ev.Name)
					fi, err := os.Stat(ev.Name)
					if err == nil && fi.IsDir() {
						pR.watch.Remove(ev.Name)
						Info("REMOVE MONITOR: ", ev.Name)
					}
				}

				// TODO: modify name
				if ev.Op&fsnotify.Rename == fsnotify.Rename {
					Info("RENAME: ", ev.Name)
					pR.watch.Remove(ev.Name)
				}

				if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
					Info("CHMOD: ", ev.Name)
				}
			}
		case err := <-pR.watch.Errors:
			{
				Fatal("ERROR IN WATCH EVENT: ", err)
				return
			}
		}
	}
}
