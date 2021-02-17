package orm

import (
	err "../../cc/err"
)

func Sync2VP() {
	e := engine_vp.Sync2(new(VPModel)); err.Assert( e )
}

func VPGetProjectListById( id string ) ( vpis []VPInfoModel, e error) {
	var vps []VPModel
	e = engine_vp.Table("v_p_model").Where( "owner_id = ?", id).Find(&vps)
	for _, vp := range vps {
		vpis = append(vpis, VPInfoModel{
			Id:      vp.Id,
			Title:   vp.Title,
			OwnerId: vp.OwnerId,
		})
	}
	return
}

func VPFindProjectById( vpid string ) ( bool, VPModel, error ) {
	vp := VPModel{}
	exi, e := engine_vp.Table("v_p_model").ID( vpid ).Get(&vp)
	return exi, vp, e
}