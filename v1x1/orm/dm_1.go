package orm

import (
	err "../../cc/err"
	cfg "../../config"
)

func Sync2DM1() {
	e := engine_dm.Sync2(new(DMTargetResource)); err.Check( e )
	e = engine_dm.Sync2(new(DMTargetResourceEx)); err.Check( e )
	e = engine_dm.Sync2(new(DMBackupResource)); err.Check( e )
	e = engine_dm.Sync2(new(DMTag)); err.Check( e )
}

const (
	DM_PERMISSION_ALL = 0
	DM_PERMISSION_GUEST = 1
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
// 如果是上帝id则会直接通过
func DMCheckPermission( accountId int64 ) ( hasPermission bool, er error ) {
	hasPermission = false

	// 上帝id，拥有dm访问的所有权限
	if accountId == cfg.DMGodId {
		return true, nil
	}

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
