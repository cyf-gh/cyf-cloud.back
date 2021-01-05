package orm

import (
	"errors"
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
		ChildId int64
		ChildGenre string
	}
	DMTargetResourceEx struct {
		Id int64
		ParentId int64
		CompatibilityDescription string
		Commentary string
	}
	// backup
	DMBackupResource struct {
		Id int64
		ParentId int64
		// 例如一个备份一个文件夹，则其中的所有文件都会被该备份资源所备份
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
		return DMGetBackupResourceById( PR.ChildId )
	} else {
		return nil, errors.New("resource is not a backup resource, but try to get a backup description")
	}
}

func ( PR *DMTargetResource ) GetChildResourceEx() ( ex *DMTargetResourceEx, e error ) {
	if PR.ChildGenre != "backup" {
		return DMGetTargetResourceExById( PR.ChildId )
	} else {
		return nil, errors.New("resource is a backup resource, but try to get a ex description")
	}
}

// 返回与 R md5相同的资源
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

func ( R DMTag ) Insert() ( e error ) {
	return DMInsertTag( R )
}