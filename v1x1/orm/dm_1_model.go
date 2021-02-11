package orm

import (
	"errors"
	"../../dm_1"
)
type (
	// target
	DMTargetResource struct {
		Id int64
		Description string
		MD5 string
		Path string `xorm:"unique"`
		TagIds []int64
		BackupIdList []int64
		ChildGenre string
		Rating int
		Dead bool
	}
	DMTargetResourceEx struct {
		Id int64
		ParentId int64 `xorm:"unique"`
		CompatibilityDescription string
		Commentary string
	}
	// backup
	DMBackupResource struct {
		Id int64
		ParentId int64
		Description string
		MD5 string
		Path string `xorm:"unique"`
		// 例如：备份一个文件夹，则其中的所有文件都会被该备份资源所备份
		BackupTargetIds[] int64
	}
	// tag
	DMTag struct {
		Id int64
		Name string `xorm:"unique"`
	}
	// permission
	DMWhiteList struct {
		Id int64
		Level int `xorm:"unique"`
		AccountIds [] int64
	}
)

func ( PR *DMTargetResource ) GetChildResourceBackup() ( *DMBackupResource, error )  {
	if PR.ChildGenre == "backup" {
		b := &DMBackupResource{}
		_, e := engine_dm.Table("d_m_backup_resource").Where("parent_id = ?", PR.Id).Get( b ); if e != nil {
			return b, e
		} else {
			return b, nil
		}
	} else {
		return nil, errors.New("resource is not a backup resource, but try to get a backup description")
	}
}

func ( PR *DMTargetResource ) GetChildResourceEx() ( ex *DMTargetResourceEx, e error ) {
	if PR.ChildGenre != "backup" {
		b := &DMTargetResourceEx{}
		_, e := engine_dm.Table("d_m_target_resource_ex").Where("parent_id = ?", PR.Id).Get( b ); if e != nil {
			return b, e
		} else {
			return b, nil
		}
	} else {
		return nil, errors.New("resource is a backup resource, but try to get a ex description")
	}
}

// 退化为DMResource
// 带文件信息
func ( R DMTargetResource ) Decay() ( *dm_1.DMResource, error ) {
	r := &dm_1.DMResource{
		Path: R.Path,
	}
	_, e := r.GetBasicFileInfo()
	return r, e
}

func ( pR *DMTargetResource ) ComputeMD5() error {
	r, e := pR.Decay()
	md5, e := r.GetMD5()
	pR.MD5 = md5
	return e
}


// 返回与 R md5相同的资源
// 包含自身
func ( R DMTargetResource ) GetClones() ( rs []DMTargetResource, e error ) {
	e = engine_dm.Table("d_m_target_resource").Where("m_d5 = ?", R.MD5 ).Find( &rs )
	return
}

func ( R DMBackupResource ) Update() ( e error ) {
	_, e = engine_dm.Table("d_m_backup_resource").Cols().ID(R.Id).Update(R)
	return
}

func ( R DMTargetResource ) Update() ( e error ) {
	_, e = engine_dm.Table("d_m_target_resource").Cols().ID(R.Id).Update(R)
	return
}

func ( R DMTargetResourceEx ) Update() ( e error ) {
	_, e = engine_dm.Table("d_m_target_resource_ex").Cols().ID(R.Id).Update(R)
	return
}

func DMTags2Strings( tags []DMTag ) ( ts []string ) {
	ts = []string{}
	for _, t := range tags {
		ts = append(ts, t.Name)
	}
	return
}

func DMGetAllTags() ( tags []DMTag, e error ) {
	e = engine_dm.Table("d_m_tag").Find(&tags)
	return
}

func DMIsTagExist( name string ) bool {
	t := &DMTag{}
	exist, _ := engine_dm.Table("d_m_tag").Where("name = ?", name ).Get(t)
	return exist
}

func DMGetTagIdByName( name string ) int64 {
	t := &DMTag{}
	engine_dm.Table("d_m_tag").Where("name = ?", name).Get( t )
	return t.Id
}

func ( R DMTag ) Insert() ( e error ) {
	return DMInsertTag( R )
}