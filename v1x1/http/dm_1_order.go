// /order所有操作都会添加数据至数据库
package http

import (
	"../../cc"
	"../../cc/err"
	"../../dm_1"
	"../orm"
	"encoding/json"
	"github.com/dhowden/tag"
	"os"
)

type ID3 struct {
	Album       string `json:"album"`
	AlbumArtist string `json:"albumArtist"`
	Artist      string `json:"artist"`
	Comment     string `json:"comment"`
	Composer    string `json:"composer"`
	DiscNum     int    `json:"disc_num"`
	DiscTotal   int    `json:"disc_total"`
	FileType    string `json:"fileType"`
	Format      string `json:"format"`
	Genre       string `json:"genre"`
	Lyrics      string `json:"lyrics"`
	Picture     interface{} `json:"picture"`
	Raw interface{} `json:"raw"`
	Title      string `json:"title"`
	TrackNum   int    `json:"track_num"`
	TrackTotal int    `json:"track_total"`
	Year       int    `json:"year"`
}

func ComputeAllMD5() {
	rs, _ := orm.GetAllMD5NotComputed()
	if e := dm_1.TaskSharedList.AddTask( "compute_md5", true, 100000, len(rs), 1000 ); e != nil {
		return
	}
	t := &dm_1.TaskSharedList.Lists["compute_md5"][0]

	t.ProgressMax = len(rs)
	t.Progress = 0
	for _, r := range rs {
		t.TryLock()
		t.ProgressStage = "computing md5s..."
		t.CurrentMsg = r.Path
		r.ComputeMD5()
		r.Update()
		t.Progress++
	}
	t.Finished()
}

func GetAllID3() {
	ms, e := orm.DMGetAllMusicResources(); err.Assert( e )
	if e := dm_1.TaskSharedList.AddTask( "load_id3", true, 100000, len(ms), 1000 ); e != nil {
		return
	}
	t := &dm_1.TaskSharedList.Lists["load_id3"][0]

	t.ProgressMax = len(ms)
	t.Progress = 0

	for _, m := range ms {
		t.TryLock()
		t.ProgressStage = "loading id3..."
		t.CurrentMsg = m.Path

		ex, e := m.GetChildResourceEx(); if e != nil { t.Error( e ); continue }
		if ex.Data != "" { // 不重复读取id3
			continue
		}
		f, e := os.Open( m.Path ); if e != nil { t.Error( e ); continue }
		id3, e := tag.ReadFrom(f); if e != nil { t.Error( e ); continue }

		dNum, totalD := id3.Disc()
		tNum, totalT := id3.Track()
		ppic := id3.Picture()
		p := tag.Picture{}
		if ppic != nil {
			p = *ppic
		}
		bs, e := json.Marshal(
			cc.H {
				"format" : id3.Format(),
				"fileType" : id3.FileType(),
				"title" : id3.Title(),
				"album" : id3.Album(),
				"albumArtist" : id3.AlbumArtist(),
				"comment" : id3.Comment(),
				"composer" : id3.Composer(),
				"disc_num" : dNum,
				"disc_total" : totalD,
				"genre" : id3.Genre(),
				"lyrics" : id3.Lyrics(),
				"picture" : p,
				"track_num" : tNum,
				"track_total" : totalT,
				"year" : id3.Year(),
				"artist" : id3.Artist(),
				"raw" : id3.Raw(),
			} ); if e != nil { t.Error( e ); continue }
		ex.Data = string( bs )
		ex.Update()
		t.Progress++
	}
	t.Finished()
}

func init() {
	cc.AddActionGroup( "/v1x1/dm/1/order", func( a cc.ActionGroup ) error {
		// \brief 开始递归所有目录进行资源索引
		// \arg[path] 要递归索引的目录
		// \note 会导致并发 任务名 order_recruit；可暂停；自旋间隔 1s
		// \return ok
		a.GET( "/recruit", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Assert( e )

			rootDir := ap.GetFormValue("d")
			if rootDir == "" {
				rootDir = dm_1.DMRootPath()
			}
			dmRootDir := &dm_1.DMResource{
				Path: rootDir,
			}
			go func() {
				if e := dm_1.TaskSharedList.AddTask( "order_recruit", true, 100000, dmRootDir.LsRecruitCount(), 1000 ); e != nil {
					return
				}
				t := &dm_1.TaskSharedList.Lists["order_recruit"][0]
				lsRootRes := dmRootDir.LsRecruit( t )
				e = orm.DMAddResources( lsRootRes, t); t.Error( e )
				t.Finished()
			} ()
			return cc.HerOk()
		} )
		// \brief 开始将id3信息载入数据库
		// \note 会导致并发 任务名 load_id3；可暂停；自旋间隔 1s
		// \return ok
		a.GET("/compute/md5", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			go ComputeAllMD5()
			return cc.HerOk()
		} )
		// \brief 开始计算数据库中的md5值
		// \note 会导致并发 任务名 compute_md5；可暂停；自旋间隔 1s
		// \return ok
		a.GET("/music/load/id3", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ){
			go GetAllID3()
			return cc.HerOk()
		} )
		// \brief 添加某个目录下的所有资源
		// \return ok
		a.GET( "/ls", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Assert( e )
			dir := ap.GetFormValue( "d" )
			dmDir, e := checkDir( dir ); err.Assert( e )
			lsRes, e := dmDir.Ls(); err.Assert( e )
			e = orm.DMAddResources( lsRes, nil ); err.Assert( e )
			return cc.HerOk()
		} )
		// \brief 添加一个或多个资源
		// \return ok
		a.POST( "", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Assert( e )
			var dmRes []dm_1.DMResource
			e = ap.GetBodyUnmarshal( &dmRes ); err.Assert( e )
			e = orm.DMAddResources( dmRes, nil ); err.Assert( e )
			return cc.HerOk()
		} )
		return nil
	} )
}