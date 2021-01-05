// 用于查询数据库中的资源
package orm

import (
	"fmt"
	"stgogo/comn/convert"
)

func DMGetTargetResourceById( id int64 ) ( tr *DMTargetResource, er error ) {
	return DMGetParentResource( id )
}

// 获取某个子资源信息的父资源信息
func DMGetParentResource( parentId int64 ) (tr *DMTargetResource, er error) {
	tr = &DMTargetResource{}
	if has, e := engine_dm.Table("d_m_target_resource").ID(parentId).Get(tr); !has || e != nil {
		er = e
		return
	}
	return
}

func DMGetTargetResourceByPath( path string ) ( tr *DMTargetResource, er error ) {
	tr = &DMTargetResource{}
	if has, e := engine_dm.Table("d_m_target_resource").Where("path = ?", path).Get(tr); !has || e != nil {
		er = e
		return
	}
	return
}


func DMGetTargetResourcesByTags( tags []string ) ( []DMTargetResource, error ) {
	var (
		pis []DMTargetResource
		findEx string
	)
	tagIds, e := GetTagIds( tags )
	findEx = ""
	for i, id := range tagIds {
		sid := convert.I64toa(id)
		// tag交集
		// id = 1
		// tag_ids like '[1,%' or like '%,1,%' or like '%,1]'
		findEx += fmt.Sprintf( "(tag_ids like '[%s,%%' or tag_ids like '%%,%s,%%' or tag_ids like '%%,%s]')", sid, sid, sid )
		if i != len(tagIds) - 1 {
			findEx += "and"
		}
	}
	findEx += " and is_private = 0"
	e = engine_dm.Table("d_m_target_resource").Where( findEx ).Find(&pis)
	return pis, e
}

func DMGetBackupResourceById( id int64 ) ( bk *DMBackupResource, er error ) {
	bk = &DMBackupResource{}
	if has, e := engine_dm.Table("d_m_backup_resource").ID( id ).Get(bk); !has || e != nil {
		er = e
		return
	}
	return
}

func DMGetTargetResourceExById( id int64 ) ( ex *DMTargetResourceEx, er error ) {
	ex = &DMTargetResourceEx{}
	if has, e := engine_dm.Table("d_m_target_resource_ex").ID( id ).Get(ex); !has || e != nil {
		er = e
		return
	}
	return
}

func DMIsTargetResourceExist( path string ) ( exist bool, e error ) {
	tr := &DMTargetResource{}
	return engine_dm.Table("d_m_target_resource").Where("path = ?", path).Get(tr)
}
