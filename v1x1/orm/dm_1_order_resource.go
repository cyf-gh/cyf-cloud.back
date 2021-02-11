// 处理将磁盘文件记录入数据库的操作
package orm

import (
	"../../dm_1"
	"errors"
	"github.com/kpango/glg"
	"sync"
)

var (
	mutexAddResource sync.Mutex
)

func addDMResGo( r dm_1.DMResource, status *dm_1.DMTaskStatus ) ( e error ) {
	if status != nil {
		status.Msg(r.Path)
	}
	var rrr *DMTargetResource
	md5, e := r.GetMD5() ; if e != nil { return }
	if exist, _ := DMIsTargetResourceExist( r.Path ); exist {
		return
	}
	// mutexAddResource.Lock()
	_, e = engine_dm.Table("d_m_target_resource" ).Insert( DMTargetResource {
		Description:  "",
		MD5:          md5,
		Path:         r.Path,
		TagIds:    	  nil,
		BackupIdList: nil,
		ChildGenre:   r.GetGenre(),
		Rating: 0,
		Dead: false,
	} ); if e != nil { goto END }
	rrr, e = DMGetTargetResourceByPath( r.Path ); if e != nil { goto END }
	if r.GetGenre() != "backup" {
		_, e = engine_dm.Table("d_m_target_resource_ex" ).Insert( DMTargetResourceEx{
			ParentId:                 rrr.Id,
			CompatibilityDescription: "",
			Commentary:               "",
		} ); if e != nil { goto END }
	} else {
		_, e = engine_dm.Table("d_m_backup_resource" ).Insert( DMBackupResource{
			ParentId:        rrr.Id,
			BackupTargetIds: nil,
		});  if e != nil { goto END }
	}
END:
	if status != nil {
		status.Progress++
		if status.Progress == status.ProgressMax {
			status.Finished()
			glg.Log( "finished" )
		}
	}
	if status != nil && e != nil {
		status.Error( e )
	}
	// mutexAddResource.Unlock()
	return
}

func addDMRes( r dm_1.DMResource, status *dm_1.DMTaskStatus ) ( e error ) {
	if status != nil {
		status.Msg(r.Path)
	}
	var rrr *DMTargetResource
	md5, e := r.GetMD5Mask() ; if e != nil { return }
	if exist, _ := DMIsTargetResourceExist( r.Path ); exist {
		return
	}
	// mutexAddResource.Lock()
	_, e = engine_dm.Table("d_m_target_resource" ).Insert( DMTargetResource {
		Description:  "",
		MD5:          md5,
		Path:         r.Path,
		TagIds:    	  nil,
		BackupIdList: nil,
		ChildGenre:   r.GetGenre(),
		Rating: 0,
		Dead: false,
	} ); if e != nil { goto END }
	rrr, e = DMGetTargetResourceByPath( r.Path ); if e != nil { goto END }
	if r.GetGenre() != "backup" {
		_, e = engine_dm.Table("d_m_target_resource_ex" ).Insert( DMTargetResourceEx{
			ParentId:                 rrr.Id,
			CompatibilityDescription: "",
			Commentary:               "",
		} ); if e != nil { goto END }
	} else {
		_, e = engine_dm.Table("d_m_backup_resource" ).Insert( DMBackupResource{
			ParentId:        rrr.Id,
			BackupTargetIds: nil,
		});  if e != nil { goto END }
	}
END:
	if status != nil {
		status.Progress++
		if status.Progress == status.ProgressMax {
			status.Finished()
			glg.Log( "finished" )
		}
	}
	if status != nil && e != nil {
		status.Error( e )
	}
	// mutexAddResource.Unlock()
	return
}


func addDMResGroup( rs []dm_1.DMResource, status *dm_1.DMTaskStatus ) ( e error ) {
	for _, rr := range rs {
		addDMResGo( rr, status )
	}
	return
}

func GetAllMD5NotComputed() ( rs []DMTargetResource, e error ) {
	rs = []DMTargetResource{}
	e = engine_dm.Table("d_m_target_resource").Where( "m_d5 = ?", "md5notcomputed").Find(&rs)
	return
}

// 当添加一笔数据时，系统将会根据dm1.DMExts自动匹配文件的种类
// 只要有一笔数据为重复添加的，就会立刻返回错误
func DMAddResources( rs []dm_1.DMResource, status *dm_1.DMTaskStatus ) ( e error ) {
	if status != nil {
		status.MsgStage( "ordering into database..")
		status.Progress = 0
		status.ProgressMax = len(rs)
	}
	for _, rr := range rs {
		addDMRes( rr, status )
	}
	return
}

// 弃用
// MD5 运算无法使用并发提高速度，因为有读取文件的磁盘速度短板
func DMAddResourcesGo( rs []dm_1.DMResource, status *dm_1.DMTaskStatus, goCount int ) ( e error ) {
	if status != nil {
		status.MsgStage( "ordering into database..")
		status.Progress = 0
		status.ProgressMax = len(rs)
	}
	var delta int = len(rs) / goCount
	start := 0
	end := 0

	for i := 0; i < goCount; i++ {
		if start+delta > len(rs) {
			end = len(rs)
		} else {
			end = start+delta
		}
		go addDMResGroup( rs[start:end], status )
		start = end + 1
	}
	return
}

func DMInsertTag( tag DMTag ) error {
	t := &DMTag{}
	exists, _ := engine_dm.Table("d_m_tag").Where( "name = ?", tag.Name).Get(t)
	if exists {
		return errors.New("tag exists")
	}
	_, e := engine_dm.Table("d_m_tag" ).Insert( tag )
	return e
}

// 指定资源为二进制，或二进制目录
func DMIndicateResourceBinary( id int64 ) ( e error ) {
	tr, e := DMGetTargetResourceById( id )
	tr.ChildGenre = "binary"
	e = tr.Update()
	return
}

// 指定某文件为某资源的备份
func DMIndicateResourceBackupOf( resourceId int64, backupId int64 ) ( e error ) {
	resource, e := DMGetTargetResourceById( resourceId )
	if e != nil { return }
	bk, e := DMGetTargetResourceById( backupId )
	if e != nil { return }
	resource.BackupIdList = append(resource.BackupIdList, backupId)
	backup, e := bk.GetChildResourceBackup()
	if e != nil { return }
	backup.BackupTargetIds = append(backup.BackupTargetIds, resourceId)
	backup.Update()
	return
}
