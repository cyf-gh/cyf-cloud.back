// 编写频繁操作的api
// @see http/README.md
package http

type (
	LikeInfoModel struct {
		Count int
		Liked bool
	}
)

const (
	_postViewPrefix = "$post_view$"
	_postLikeItPrefix = "$post_like_it$"
)



