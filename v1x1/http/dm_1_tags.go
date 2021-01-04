package http

import (
	"../../cc"
	"../../cc/err"
	"../orm"
	"errors"
)

func init() {
	cc.AddActionGroup( "/v1x1/dm/1/tags", func( a cc.ActionGroup ) error {
		// \brief 获取所有的tag
		a.GET( "/all", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Check( e )
			tags, e := orm.DMGetAllTags(); err.Check( e )
			return cc.HerOkWithData( tags )
		} )
		// \brief 与post的tag行为不同，post中的tag通过发送文章时自动添加，dm_1的tag需要手动添加后才可使用
		// \arg[a] tag的名字
		a.GET( "/add", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Check( e )
			tagName := ap.GetFormValue("a"); if tagName == "" { panic(errors.New("empty param: a"))}
			if !orm.DMIsTagExist( tagName ) {
				t := orm.DMTag{
					Name: tagName,
				}
				t.Insert()
				return cc.HerOk()
			} else {
				panic( errors.New("tag: " + tagName + " already exists") )
			}
		} )
		return nil
	} )
}
