package http

import (
	"../../cc"
	"../../cc/err"
	"../orm"
	"github.com/gorilla/feeds"
	"stgogo/comn/convert"
	"time"
)

func posts2FeedItems( ps []orm.Post, pfeed **feeds.Feed ) {
	f := "2006-01-02 15:04:05"
	for _, p := range ps  {
		created, _ := time.Parse(f, p.CreateDate)
		updated, _ := time.Parse(f, p.Date)
		tl := len(p.Text)

		if tl > 30 {
			tl = 30
		}

		a, e := orm.GetAccount( p.OwnerId ); err.Check( e )
		(*pfeed).Items = append((*pfeed).Items, &feeds.Item{
			Title:       p.Title,
			Link:        &feeds.Link{ Href:"https://se.cyf-cloud.cn:8888/post/reader?id="+ convert.I64toa( p.Id ) },
			Author:      &feeds.Author{ a.Name, a.Email },
			Id:          convert.I64toa( p.Id ),
			Description: p.Text[1:tl],
			Updated:     updated,
			Created:     created,
		})
	}
}

func init() {
	cc.AddActionGroup( "/v1x1/feeds", func( a cc.ActionGroup ) error {
		// \brief 返回订阅post的结果
		// \type GET
		// \arg[uid] 订阅的作者id，如果为空则订阅所有的文章
		// \arg[a] 订阅的类型，可为 rss/atom/json
		a.GET_DO( "/post",  func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			var (
				uid, pfeedData string
				pfeed *feeds.Feed
				e error
			)

			if uid = ap.R.FormValue("uid"); uid == "" {
				// 订阅所有文章
				pfeed = &feeds.Feed {
					Title:       "cyf-cloud blog",
					Link:        &feeds.Link{Href: "https://se.cyf-cloud.cn:8888/posts"},
					Description: "cyf-cloud blog 2020 - Writing Everything",
					Author:      &feeds.Author{Name: "everyone", Email: "nil"},
					Created:     time.Now(),
				}
				ps, e := orm.GetPostsPublicAll(); err.Check( e )
				posts2FeedItems( ps, &pfeed )
			} else {
				// 订阅某人的文章
				id, e := convert.Atoi64( uid ); err.Check( e )
				a, e := orm.GetAccount( id ); err.Check( e )
				pfeed = &feeds.Feed {
					Title:       "["+ a.Name + "] cyf-cloud blog",
					Link:        &feeds.Link{Href: "https://se.cyf-cloud.cn:8888/user/home?id="+ convert.I64toa( a.Id ) },
					Description: "cyf-cloud blog 2020 - Writing Everything",
					Author:      &feeds.Author{Name: a.Name, Email: a.Email},
					Created:     time.Now(),
				}
				ps, e := orm.GetPostsByOwnerPublic( id ); err.Check( e )
				posts2FeedItems( ps, &pfeed )
			}
			switch ap.R.FormValue("a") {
			case "rss":
				pfeedData, e = pfeed.ToRss(); err.Check( e )
				break
			case "atom":
				pfeedData, e = pfeed.ToAtom(); err.Check( e )
				break
			case "json":
				pfeedData, e = pfeed.ToJSON(); err.Check( e )
				break
			default:
				pfeedData = "please choose rss/atom/json"
				break
			}
			return cc.HerData( pfeedData )
		} )

	    return nil
	} )
}
