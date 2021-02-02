package http

import (
	"../../cc"
	"../../cc/err"
	orm "../orm"
	"errors"
	"stgogo/comn/convert"
)

func init() {
	cc.AddActionGroup( "/v1x1/vp", func( a cc.ActionGroup ) error {
		a.POST( "/update", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			id, e := GetIdByAtk( ap.R ); err.Check( e )
			vp := &orm.VPModel{}
			e = ap.GetBodyUnmarshal( vp ); err.Check( e )
			// 无权限
			if vp.OwnerId != id {
				panic( errors.New("you have no permission to modify this visual progress project") )
			}
			// 不存在该vp，添加
			if vp.Id == 0 {
				iid, e := vp.Insert(); err.Check( e )
				return cc.HerOkWithData( iid )
			} else {
				vp.Update()
				return cc.HerOkWithData( 0 )
			}
		} )
		a.GET( "/projects/list", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			id, e := GetIdByAtk( ap.R ); err.Check( e )
			// vpid := ap.GetFormValue("id"); err.Check( e )

			vpis, e := orm.VPGetProjectListById( convert.I64toa( id ) ); err.Check( e )
			return cc.HerOkWithData( vpis )
		} )
		a.GET( "/project", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			id, e := GetIdByAtk( ap.R ); err.Check( e )
			vpid := ap.GetFormValue("id"); err.Check( e )
			exi, vp, e := orm.VPFindProjectById( vpid ); err.Check( e )
			if !exi {
				panic( "specified project does not exist" )
			}
			if vp.OwnerId != id {
				panic( "you have no permission to this project" )
			}
			return cc.HerOkWithData( vp )
		} )
		return nil
	} )
}
