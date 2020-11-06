// TODO: 文章是否可以删除？
// TODO: 文章的访问权限问题？
package orm

import (
	err "../../cc/err"
	"fmt"
	"stgogo/comn/convert"
	"time"
)

// 上传的文章
type (
	Post struct {
		Id int64
		Title string
		Text string
		TagIds[] int64
		OwnerId int64
		IsPrivate bool
		Date string
		CreateDate string
	}
	PostInfo struct {
		Id int64
		Title string
		CreateDate string
		IsPrivate bool
		OwnerId int64
		TagIds[] int64
	}
)
// 上传的tag标签结构
// tag同时将作为分区的主要依据
// TODO:tag达到一定数量则升级为分区
type (
	Tag struct {
		Id int64
		Text string
		IsCatalog bool
		Percentage float32
	}
)

func Sync2Post() {
	e := engine_post.Sync2(new(Post))
	err.Check( e )
	e = engine_post.Sync2(new(Tag))
	err.Check( e )
}

// 通过某一个人获取所有他的文章
// 仅通过atk
func GetPostInfosByOwnerAll( OwnerId int64 ) ( []PostInfo, error ) {
	var posts []PostInfo
	e := engine_post.Table("Post").Where( "owner_id = ?", OwnerId).Find(&posts)
	return posts, e
}

func GetAllPublicPostInfosLimited( start, count int ) ( []PostInfo, error ) {
	var posts []PostInfo
	e := engine_post.Table("Post").Where( "is_private = 0" ).Limit( count, start ).Find(&posts)
	return posts, e
}

// 通过某一个人获取所有公开文章
func GetPostInfosByOwnerPublic( OwnerId int64 ) ( []PostInfo, error ) {
	var posts []PostInfo
	e := engine_post.Table("Post").Where( "owner_id = ? and is_private = 0", OwnerId).Find(&posts)
	return posts, e
}

func GetPostInfosAll() ( []PostInfo, error ) {
	var posts []PostInfo
	e := engine_post.Table("Post").Where( "is_private = 0" ).Find(&posts)
	return posts, e
}

func GetPostInfoById( id int64 ) ( PostInfo, error ) {
	post := new(PostInfo)
	_, e := engine_post.Table("Post").ID( id ).Get(post)
	return *post, e
}

func GetPostById( id int64 ) ( Post, error ) {
	post := new(Post)
	_, e := engine_post.Table("Post").ID( id ).Get(post)
	return *post, e
}

func GetPostInfosByIds( ids []int64 ) ( []PostInfo, error ) {
	var ps []PostInfo
	for _, id := range ids {
		if p, e := GetPostInfoById(id); e != nil {
			return nil, e
		} else {
			ps = append( ps, p )
		}
	}
	return ps, nil
}

// 向数据库添加一笔新文章
func NewPost( title, text string, owner int64, tags []string, private bool) (int64, error) {
	tagIds, e  := GetTagIds( tags )
	newPost := &Post{
		Title:     title,
		Text:      text,
		TagIds:    tagIds,
		OwnerId: owner,
		IsPrivate: private,
		Date: time.Now().Format("2006-01-02 15:04:05"),
		CreateDate: time.Now().Format("2006-01-02 15:04:05"),
	}

	_, e = engine_post.Table("Post").Insert( newPost )
	// err.Check( e )
	return newPost.Id, e
}

// 修改文章
func ModifyPost( id int64, title, text string, owner int64, isPrivate bool, tags []string) error {
	tagIds, e := GetTagIds( tags )
	mp := &Post{
		Title:     title,
		Text:      text,
		TagIds:    tagIds,
		OwnerId:   owner,
		IsPrivate: isPrivate,
		Date: time.Now().Format("2006-01-02 15:04:05"),
	}
	_, e = engine_post.Table("Post").AllCols().ID(id).Update(mp)
	return e
}

// 修改文章，不修改内容
// 减轻流量负担
func ModifyPostNoText( id int64, title string, owner int64, tags []string) error {
	tagIds, e := GetTagIds( tags )

	_, e = engine_post.Table("Post").ID(id).Update(&Post{
		Title:     title,
		TagIds:    tagIds,
		OwnerId:   owner,
	})
	return e
}

// 根据一系列tag的名字获取所有的tag的id
// 如果tag不存在则会被创建
func GetTagIds( tags []string ) ( []int64, error ) {
	var tagIds []int64

	for _, tagText := range tags  {
		t := new(Tag)
		GetTag:
		exists, _ := engine_post.Table("Tag").Where( "Text = ?", tagText).Get(t)

		// 如果当前没有该tag，则创建一个新tag
		if !exists {
			_, e := engine_post.Table("Tag").Insert( &Tag {
				Text:tagText,
				IsCatalog:false,
			})
			if e != nil {
				return nil, e
			}
			goto GetTag // 再次获取tag
		}
		// 这里t应该已经被填充
		tagIds = append(tagIds, t.Id)
	}
	return tagIds, nil
}

func GetTagNames( tagIds []int64 ) ( []string, error ) {
	var (
		tags []string
	)

	for _, id := range tagIds {
		tag := new(Tag)
		if exists, e := engine_post.Table("Tag").ID(id).Get(tag); exists && e == nil {
			tags = append(tags, tag.Text)
		} else {
			return nil, e
		}
	}
	return tags, nil
}

func GetPostInfosByTags( tags []string ) ( []PostInfo, error ) {
	var (
		pis []PostInfo
		findEx string
	)
	tagIds, e := GetTagIds( tags )
	findEx = ""
	for i, id := range tagIds {
		sid := convert.I64toa(id)
		// tag交集
		// id = 1
		// tag_ids like '[1,%' or like '%,1,%' or like '%,1]'
		findEx += fmt.Sprintf( "(tag_ids like '[%s,%%' or tag_ids like '%%,%s,%%' or tag_ids like '%%,%s]')", sid, sid, sid )
		if i != len(tagIds) - 1 {
			findEx += "and"
		}
	}
	findEx += " and is_private = 0"
	e = engine_post.Table("Post").Where( findEx ).Find(&pis)
	return pis, e
}

func GetAllTags() ( []Tag, error ) {
	var (
		tags []Tag
		e error
	)
	e = engine_post.Table("Tag").Find(&tags)
	return tags, e
}