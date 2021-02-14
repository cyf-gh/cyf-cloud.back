package http

import (
	"../../cc"
	err "../../cc/err"
	"../orm"
	"encoding/json"
	"github.com/dhowden/tag"
	"os"
	"strconv"
)

func init() {
	cc.AddActionGroup( "/v1x1/dm/1/music", func( a cc.ActionGroup ) error {
		// \brief （弃用）获取所有音乐资源的id3信息
		// \note 一定会有io瓶颈
		a.GET( "/tag/io", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			return cc.HerOk()
			var mds []tag.Metadata
			ms, e := orm.DMGetAllMusicResources(); err.Check( e )
			for _, m := range ms {
				f, e := os.Open( m.Path ); if e != nil { continue }
				id3, e := tag.ReadFrom(f); if e != nil { continue }
				mds = append( mds, id3 )
			}
			return cc.HerOkWithData(mds)
		} )
		// \brief 获取所有音乐资源的id3信息，通过数据库内查询
		// \arg[head] 开始的位置
		// \arg[end] 结束位置，当为-1时为数组长度
		// \arg[with_pic_data] 不为空时读取携带图片data的数据
		// \return {
		// 	"musics" 音乐数据
		// 	"totalCount" 音乐总个数
		// }
		a.GET( "/ex", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			var exs []orm.DMTargetResourceEx
			ms, e := orm.DMGetAllMusicResources(); err.Check( e )
			strhead := ap.GetFormValue( "head" )
			strend := ap.GetFormValue( "end" )
			withPicData := ap.GetFormValue( "with_pic_data" )

			var head, end int
			if strhead != "" && strend != "" {
				head, e = strconv.Atoi( strhead ); err.Check( e )
				end, e = strconv.Atoi( strend ); err.Check( e )
			}
			count := len(ms)
			for _, m := range ms[head:end] {
				ex, _ := m.GetChildResourceEx(); if ex == nil { continue }
				// 如果with_pic_data为空，则不传输picture data数据
				id3 := ID3{}
				if withPicData != "" {
					e := json.Unmarshal( []byte(ex.Data), &id3 ); if e != nil { continue }
					id3.Picture.Data = ""
				}
				bs, e := json.Marshal( id3 ); if e != nil { continue }
				ex.Data = string( bs )
				//
				exs = append( exs, *ex )
			}


			return cc.HerOkWithData( cc.H{
				"musics": exs,
				"totalCount": count,
			} )
		} )
		return nil
	} )
}