package orm

import (
	"errors"
)

func AddFav( id, pid int64 ) (*AccountEx, error) {
	a := &AccountEx{}
	has, e := engine.Table("account_ex").ID(id).Get(a)
	if e != nil {
		return nil, e
	} else if !has {
		return nil, errors.New("account not found")
	}
	a.FavPosts = append(a.FavPosts, pid)
	engine.Table("account_ex").ID(id).Update( a.FavPosts )
	return a, nil
}

func UpdateFav( id int64, mdFavList []int64 ) (*AccountEx, error) {
	a := &AccountEx{}
	has, e := engine.Table("account_ex").ID(id).Get(a)
	if e != nil {
		return nil, e
	} else if !has {
		return nil, errors.New("account not found")
	}
	a.FavPosts = mdFavList
	engine.Table("account_ex").ID(id).Update( a.FavPosts )
	return a, nil
}