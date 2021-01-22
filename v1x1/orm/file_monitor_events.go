// 通过文件监控使数据库自适应
package orm

import (
	"../../dm_1"
	"github.com/fsnotify/fsnotify"
	"github.com/kpango/glg"
)

func createDb( event fsnotify.Event ) {
	DMAddResources( []dm_1.DMResource{ { Path: event.Name } }, nil )
}

func removeDb( event fsnotify.Event ) {
	if tr, e := DMGetTargetResourceByPath(event.Name); e != nil {
		glg.Error( "[WATCHER] REMOVE: ", e )
	} else {
		tr.Dead = true
		tr.Update()
	}
}

func writeDb( event fsnotify.Event ) {
	if tr, e := DMGetTargetResourceByPath(event.Name); e != nil {
		glg.Error( "[WATCHER] WRITE: ", e )
	} else {
		if r, e := tr.Decay(); e != nil { glg.Error( "[WATCHER] WRITE while resource decaying: ", e ) } else {
			if tr.MD5, e = r.GetMD5(); e != nil { glg.Error( "[WATCHER] WRITE while getting new md5: ", e ) } else {
				tr.Update()
			}
		}
	}
}

func renameDb( event fsnotify.Event ) {
	if tr, e := DMGetTargetResourceByPath(event.Name); e != nil {
		glg.Error( "[WATCHER] RENAME: ", e )
	} else {
		tr.Dead = true
		tr.Update()
	}
}

func DMNewNotifyFileDb() *dm_1.DMNotifyFile {
	w := new(dm_1.DMNotifyFile)
	w.Watch, _ = fsnotify.NewWatcher()

	w.Create = createDb
	w.Remove = removeDb
	w.Rename = renameDb
	w.Write = writeDb
	w.Chmod = func(event fsnotify.Event) {}

	return w
}