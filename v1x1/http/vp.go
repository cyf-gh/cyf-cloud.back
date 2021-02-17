package http

import (
	"../../cc"
	"../../cc/err"
	cfg "../../config"
	"../../vp"
	orm "../orm"
	"bytes"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net/http"
	"stgogo/comn/convert"
	"strconv"
	"strings"
	"time"
)

type (
	VPExportModel struct {
		Id int64
		Title string
		Base64 string
	}
)

func init() {
	cc.AddActionGroup( "/v1x1/vp", func( a cc.ActionGroup ) error {
		a.POST( "/update", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			id, e := GetIdByAtk( ap.R ); err.Assert( e )
			vp := &orm.VPModel{}
			e = ap.GetBodyUnmarshal( vp ); err.Assert( e )
			// 无权限
			if vp.OwnerId != id {
				panic( errors.New("you have no permission to modify this visual progress project") )
			}
			// 不存在该vp，添加
			if vp.Id == 0 {
				iid, e := vp.Insert(); err.Assert( e )
				return cc.HerOkWithData( iid )
			} else {
				vp.Update()
				return cc.HerOkWithData( 0 )
			}
		} )
		a.GET( "/projects/list", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			id, e := GetIdByAtk( ap.R ); err.Assert( e )
			// vpid := ap.GetFormValue("id"); err.Assert( e )

			vpis, e := orm.VPGetProjectListById( convert.I64toa( id ) ); err.Assert( e )
			return cc.HerOkWithData( vpis )
		} )
		a.GET( "/project", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			id, e := GetIdByAtk( ap.R ); err.Assert( e )
			vpid := ap.GetFormValue("id"); err.Assert( e )
			exi, vp, e := orm.VPFindProjectById( vpid ); err.Assert( e )
			if !exi {
				panic( "specified project does not exist" )
			}
			if vp.OwnerId != id {
				panic( "you have no permission to this project" )
			}
			return cc.HerOkWithData( vp )
		} )
		a.POST_CONTENT( "/export", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			id, e := GetIdByAtk( ap.R ); err.Assert( e )
			m := &VPExportModel{}
			e = ap.GetBodyUnmarshal( m ); err.Assert( e )

			exi, v, e := orm.VPFindProjectById( convert.I64toa( m.Id ) ); err.Assert( e )

			// process image
			// https://stackoverflow.com/questions/31031589/illegal-base64-data-at-input-byte-4-when-using-base64-stdencoding-decodestrings/49413861
			b64data := m.Base64[strings.IndexByte( m.Base64, ',') + 1 :]
			ib, e := base64.StdEncoding.DecodeString( b64data ); err.Assert( e )

			if !exi {
				panic( "specified project does not exist" )
			}
			if v.OwnerId != id {
				panic( "you have no permission to this project" )
			}
			exportPath := cfg.VPTmpPath + "/" + m.Title +".xlsx"
			e = vp.Export( v.Data, cfg.VPTemplatePath, exportPath, ib )
			err.Assert( e )

			downloadBytes, e := ioutil.ReadFile(exportPath); err.Assert( e )

			w := *ap.W
			r := ap.R
			w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
			w.Header().Set("Content-Disposition", "attachment; filename="+m.Title +".xlsx")
			w.Header().Set("Content-Length", strconv.Itoa(len(string(downloadBytes))))

			http.ServeContent(w, r, m.Title +".xlsx", time.Now(), bytes.NewReader(downloadBytes))
			return cc.HerOk()
		} )

		return nil
	} )
}
