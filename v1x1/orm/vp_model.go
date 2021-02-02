package orm

type (
	// target
	VPModel struct {
		Id int64
		Title string `xorm:"unique"`
		OwnerId int64	// 前端必须返回正确的ownerid，不得为0
		Data string
	}
	VPInfoModel struct {
		Id int64
		Title string
		OwnerId int64
	}
)

func ( R VPModel ) Update() ( e error ) {
	_, e = engine_vp.Table("v_p_model").Cols().ID(R.Id).Update(R)
	return
}

func ( R VPModel ) Insert() ( id int64, e error ) {
	_, e = engine_vp.Table("v_p_model" ).Insert( R )
	in := &VPModel{}
	_, e = engine_vp.Table("v_p_model").Where("title = ?", R.Title).Get(in)
	id = in.Id
	return
}