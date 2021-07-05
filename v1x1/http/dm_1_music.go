package http

import (
	"../../cc"
	err "../../cc/err"
	"../orm"
	"bytes"
	"encoding/json"
	"github.com/dhowden/tag"
	"github.com/kpango/glg"
	lyrics "github.com/rhnvrm/lyric-api-go"
	"io/ioutil"
	"net/http"
	"os"
	"stgogo/comn/convert"
	"strconv"
	"strings"
	"time"
)

func init() {
	cc.AddActionGroup( "/v1x1/dm/1/music", func( a cc.ActionGroup ) error {
		// \brief （弃用）获取所有音乐资源的id3信息
		// \note 一定会有io瓶颈
		a.GET( "/tag/io", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			return cc.HerOk()
			var mds []tag.Metadata
			ms, e := orm.DMGetAllMusicResources(); err.Assert( e )
			for _, m := range ms {
				f, e := os.Open( m.Path ); if e != nil { continue }
				id3, e := tag.ReadFrom(f); if e != nil { continue }
				mds = append( mds, id3 )
			}
			return cc.HerOkWithData(mds)
		} )
		// \brief 获取所有音乐资源的id3信息，通过数据库内查询
		// \note 当音乐资源无id3信息时，将只通过元文件名填充title
		// \arg[head] 开始的位置
		// \arg[end] 结束位置，当为-1时为数组长度
		// \arg[not_default] 不为空时进行缺省数据携带
		// \arg[carry_raw] 不为空时携带raw信息
		// \arg[carry_pic] 不为空时携带picture信息
		// \arg[carry_comment] 不为空时携带comment数据
		// \return {
		// 	"musics" 音乐数据
		// 	"totalCount" 音乐总个数
		// }
		a.GET( "/ex", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Assert( e )

			var exs []orm.DMTargetResourceEx
			ms, e := orm.DMGetAllMusicResources(); err.Assert( e )
			strhead := ap.GetFormValue( "head" )
			strend := ap.GetFormValue( "end" )
			notDef := ap.GetFormValue( "not_default" )
			cRaw := ap.GetFormValue( "carry_raw" )
			cPic := ap.GetFormValue( "carry_pic" )
			cComt := ap.GetFormValue( "carry_comment" )

			var head, end int
			if strhead != "" && strend != "" {
				head, e = strconv.Atoi( strhead ); err.Assert( e )
				end, e = strconv.Atoi( strend ); err.Assert( e )
			}
			count := len(ms)
			if end == -1 {
				end = count
			}
			for _, m := range ms[head:end] {
				ex, _ := m.GetChildResourceEx(); if ex == nil { continue }
				// 如果data为空，则填充文件元数据
				if ex.Data == "" {
					pr, e := ex.GetParentResource(); if pr.Id == 0 { glg.Error( e, "| no such music. ex.GetParentResource()"); continue }
					if pr.Dead {
						glg.Error("try to get a dead music. ex.GetParentResource()")
						continue
					}
					dm, e := pr.Decay(); if e != nil { glg.Error( e ); continue }
					ext, e := dm.GetExt(); if e != nil { glg.Error( e ); continue }
					// 去掉文件后缀
					for i := 0; i < len( dm.Name ); i++ {
						if dm.Name[i] == '.' {
							dm.Name = dm.Name[:i]
							break
						}
					}
					id3 := ID3{ Title: dm.Name , FileType: ext, }
					bs, e := json.Marshal( id3 ); if e != nil { glg.Error( e ); continue }
					ex.Data = string( bs )
					exs = append( exs, *ex )
					continue
				}
				// 如果not_default不为空，则开始筛选数据
				if notDef == "" {
					id3 := ID3{}
					e := json.Unmarshal( []byte(ex.Data), &id3 ); if e != nil { glg.Error( e, ex ); continue }
					if cPic == "" { id3.Picture = nil }
					if cRaw == "" { id3.Raw = nil }
					if cComt == "" { id3.Comment = "" }
					bs, e := json.Marshal( id3 ); if e != nil { glg.Error( e ); continue }
					ex.Data = string( bs )
				}
				exs = append( exs, *ex )
			}

			return cc.HerOkWithData( cc.H{
				"musics": exs,
				"totalCount": count,
			} )
		} )
		// \brief 获取音乐的content
		// \arg[id] ParentId
		// \arg[file_type] 文件类型 例如 mp3 flac wav等 大小写无关 不可携带 .
		a.GET_CONTENT( "/raw", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Assert( e )
			pid := ap.GetFormValue( "id" )
			fileType := ap.GetFormValue("file_type")
			id, e := convert.Atoi64( pid ); err.Assert( e )
			dm, e := orm.DMGetTargetResourceById( id ); err.Assert( e )
			cbs, e := ioutil.ReadFile( dm.Path ); err.Assert( e )

			fileType = strings.ToLower( fileType )
			contentType := ""
			ok := false
			if contentType, ok = cc.ContentType[fileType]; !ok {
				return cc.HerArgInvalid("file_type" )
			}

			w := *ap.W; r := ap.R
			w.Header().Set("Content-Type", contentType )
			w.Header().Set("Content-Length", strconv.Itoa( len( string( cbs ) ) ) )

			http.ServeContent(w, r, "tmp."+fileType, time.Now(), bytes.NewReader( cbs ) )
			return cc.HerOk()
		} )
		// \brief 获取音乐的cover
		// \arg[id] Id 而不是ParentId
		// \return id3的Picture tag json
		a.GET( "/cover", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Assert( e )
			pid := ap.GetFormValue( "id" )
			id, e := convert.Atoi64( pid ); err.Assert( e )
			ex, e := orm.DMGetTargetResourceExById( id ); err.Assert( e )
			id3 := ID3{}
			e = json.Unmarshal( []byte(ex.Data), &id3 ); if e != nil {
				// 无封面
				return cc.HerOk()
			}
			return cc.HerOkWithData( id3.Picture )
		} )
		// \brief 获取音乐的lyrics
		// \arg[title] 歌曲标题（必填）
		// \arg[artist] 歌曲艺术家
		// \return 见
		a.GET_DO( "/lyrics/ex", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			type (
				LyrResult struct {
					Count  int `json:"count"`
					Code   int `json:"code"`
					Result []struct {
					Aid      int    `json:"aid"`
					Lrc      string `json:"lrc"`
					Sid      int    `json:"sid"`
					ArtistID int    `json:"artist_id"`
					Song     string `json:"song"`
					} `json:"result"`
				}
			)
			type NELyr struct {
				Result struct {
					Songs []struct {
						Name    string `json:"name"`
						ID      int    `json:"id"`
						Artists []struct {
							Name      string        `json:"name"`
							ID        int           `json:"id"`
							PicID     int           `json:"picId"`
							Img1V1ID  int           `json:"img1v1Id"`
							BriefDesc string        `json:"briefDesc"`
							PicURL    string        `json:"picUrl"`
							Img1V1URL string        `json:"img1v1Url"`
							AlbumSize int           `json:"albumSize"`
							Alias     []interface{} `json:"alias"`
							Trans     string        `json:"trans"`
							MusicSize int           `json:"musicSize"`
						} `json:"artists"`
						Album struct {
							Name       string `json:"name"`
							ID         int    `json:"id"`
							Type       string `json:"type"`
							Size       int    `json:"size"`
							PicID      int64  `json:"picId"`
							BlurPicURL string `json:"blurPicUrl"`
							CompanyID  int    `json:"companyId"`
							Pic        int64  `json:"pic"`
							PicURL     string `json:"picUrl"`
						} `json:"album"`
						Mvid   int         `json:"mvid"`
						Rtype  int         `json:"rtype"`
						Rurl   interface{} `json:"rurl"`
						Mp3URL string      `json:"mp3Url"`
					} `json:"songs"`
					SongCount int `json:"songCount"`
				} `json:"result"`
				Code int `json:"code"`
			}
			e := DM1CheckPermission( ap.R ); err.Assert( e )
			title := ap.GetFormValue( "title" )
			artist := ap.GetFormValue( "artist" )
			if title == "" {
				return cc.HerArgInvalid( "title" )
			}
			ly := &NELyr{}
			cc.PostJ( "http://music.163.com/api/search/pc", cc.H {
				"s": title,
				"offset": 0,
				"limit": 40,
				"type": 1,
			}, ly )
			glg.Log( ly, artist )
			return cc.HerData( "" )
		} )
		// \brief （弃用）获取音乐的lyrics http://gecimi.com/api/lyric/
		// \arg[title] 歌曲标题（必填）
		// \arg[artist] 歌曲艺术家
		// \return 见 http://doc.gecimi.com/en/latest/
		a.Deprecated("/lyrics/ex").GET_DO( "/lyrics", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			type (
				LyrResult struct {
					Count  int `json:"count"`
					Code   int `json:"code"`
					Result []struct {
						Aid      int    `json:"aid"`
						Lrc      string `json:"lrc"`
						Sid      int    `json:"sid"`
						ArtistID int    `json:"artist_id"`
						Song     string `json:"song"`
					} `json:"result"`
				}
			)
			e := DM1CheckPermission( ap.R ); err.Assert( e )
			title := ap.GetFormValue( "title" )
			artist := ap.GetFormValue( "artist" )
			if title == "" {
				return cc.HerArgInvalid( "title" )
			}


			{
				queryArgs := title
				if artist != "" {
					queryArgs += "/" + title
				}
				lyr := &LyrResult{}

				e = cc.GetJ( "http://gecimi.com/api/lyric/" + queryArgs, lyr ); err.Assert( e )
				if len(lyr.Result) == 0 {
					goto LAG
				}
				lycStr, e := cc.Get( lyr.Result[0].Lrc ); err.Assert( e ); if lycStr != "" {
					glg.Success("found lyric in gecimi")
					return cc.HerData( lycStr )
				}
			}
			LAG:
			// gecimi无结果，启用 https://github.com/rhnvrm/lyric-api-go
			{
				l := lyrics.New()
				lyric, e := l.Search(artist, title); err.Assert( e )
				if lyric != "" {
					glg.Success("found lyric through lyric-api-go")
					return cc.HerData( lyric )
				}
			}
			return cc.HerData( "" )
		} )
		return nil
	} )
}