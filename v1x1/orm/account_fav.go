package orm

import (
	"errors"
	"stgogo/comn/convert"
)

func GetAllFavPostInfos( id int64 ) ( pis []PostInfo, e error)  {
	a := &AccountEx{}

	has, e := engine_account.Table("account_ex").Where("account_id = ?",id).Get(a)
	if e != nil {
		return nil, e
	} else if !has {
		return nil, errors.New("account not found")
	}

	if a.FavPosts == nil {
		return []PostInfo{}, nil
	}

	inIds := "("
	for index, i := range a.FavPosts {
		inIds += convert.I64toa( i )
		if index != len(a.FavPosts) - 1 {
			 inIds += ","
		}
	}; inIds += ")"

	e = engine_post.Table("post").Where("id in " + inIds).Find(&pis)
	return
}
