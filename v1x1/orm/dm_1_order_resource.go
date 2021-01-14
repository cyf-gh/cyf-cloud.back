// 处理将磁盘文件记录入数据库的操作
package orm

import (
	"../../dm_1"
	"errors"
)

// 当添加一笔数据时，系统将会根据dm1.DMExts自动匹配文件的种类
// 只要有一笔数据为重复添加的，就会立刻返回错误
func DMAddResource( rs []dm_1.DMResource ) ( e error ) {
	var md5 string
	exist := false
	for _, r := range rs {
		md5, e = r.GetMD5() ; if e != nil { return }
		if exist, e = DMIsTargetResourceExist( r.Path ); exist {
			continue
		} else if e != nil {
			return e
		}
		_, e = engine_dm.Table("d_m_target_resource" ).Insert( DMTargetResource {
			Description:  "",
			MD5:          md5,
			Path:         r.Path,
			TagIds:    	  nil,
			BackupIdList: nil,
			ChildGenre:   r.GetGenre(),
			Rating: 0,
		} ); if e != nil { return e }
		var rr *DMTargetResource
		rr, e = DMGetTargetResourceByPath( r.Path ); if e != nil {
			return e
		}
		if r.GetGenre() != "backup" {
			_, e = engine_dm.Table("d_m_target_resource_ex" ).Insert( DMTargetResourceEx{
				ParentId:                 rr.Id,
				CompatibilityDescription: "",
				Commentary:               "",
			} ); if e != nil { return e }
		} else {
			_, e = engine_dm.Table("d_m_backup_resource" ).Insert( DMBackupResource{
				ParentId:        rr.Id,
				BackupTargetIds: nil,
			});  if e != nil { return e }
		}
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
