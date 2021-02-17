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
		// \return []DMTag
		a.GET( "/all", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Assert( e )
			tags, e := orm.DMGetAllTags(); err.Assert( e )
			return cc.HerOkWithData( tags )
		} )
		// \brief 添加一个tag
		// \note 与post的tag行为不同，post中的tag通过发送文章时自动添加，dm_1的tag需要手动添加后才可使用
		// \arg[a] tag的名字
		// \return 新添加的tag的id
		a.POST( "/add", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Assert( e )
			var (
				tagNames []string
				tagIds []int64
			)
			e = ap.GetBodyUnmarshal(&tagNames); err.Assert( e ); if len(tagNames) == 0 { panic(errors.New("empty param: a"))}
			for _, tagName := range tagNames {
				if !orm.DMIsTagExist( tagName ) {
					t := orm.DMTag{
						Name: tagName,
					}
					t.Insert()
					tagIds = append(tagIds, orm.DMGetTagIdByName(tagName))
				}
			}
			return cc.HerOkWithData( tagIds )
		} )
		return nil
		// \brief 添加一个tag
		// \note 与post的tag行为不同，post中的tag通过发送文章时自动添加，dm_1的tag需要手动添加后才可使用
		// \arg[a] tag的名字
		// \return tag id
		a.GET( "/adds", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Assert( e )
			tagName := ap.GetFormValue("a"); if tagName == "" { panic(errors.New("empty param: a"))}
			if !orm.DMIsTagExist( tagName ) {
				t := orm.DMTag{
					Name: tagName,
				}
				t.Insert()
				return cc.HerOkWithData( orm.DMGetTagIdByName(t.Name) )
			} else {
				panic( errors.New("tag: " + tagName + " already exists") )
			}
		} )
		return nil
	} )
}
