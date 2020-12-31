package orm

import (
	err "../../cc/err"
)

func Sync2DM1() {
	e := engine_dm.Sync2(new(DMTargetResource)); err.Check( e )
	e = engine_dm.Sync2(new(DMBinaryResource)); err.Check( e )
	e = engine_dm.Sync2(new(DMImageResource)); err.Check( e )
	e = engine_dm.Sync2(new(DMMusicResource)); err.Check( e )
	e = engine_dm.Sync2(new(DMVideoResource)); err.Check( e )

	e = engine_dm.Sync2(new(DMBackupResource)); err.Check( e )
	e = engine_dm.Sync2(new(DMTag)); err.Check( e )
}

const (
	DM_PERMISSION_ALL = 0
	DM_PERMISSION_GUEST = 0
	)

// 获取文件管理白名单，返回某一等级的管理权限组
// 当 wl == {} && e == nil 时表明不存在该权限组
func DMGetWhiteList( level int ) ( wl DMWhiteList, er error ){
	wl = DMWhiteList{}
	if has, e := engine_dm.Table("d_m_target_resource").Where("level = ?", level ).Get(&wl); !has || e != nil {
		er = e
		return
	}
	return
}

// 用于在API时第一时间进行权限检查
func DMCheckPermission( accountId int64 ) ( hasPermission bool, er error ) {
	hasPermission = false

	if wl, e := DMGetWhiteList( DM_PERMISSION_ALL ); e != nil {
		er = e
		return
	} else {
		for _, a := range wl.AccountIds {
			if accountId == a {
				hasPermission = true
			}
		}
	}
	return
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
