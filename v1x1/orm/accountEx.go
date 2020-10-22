package orm

import (
	"errors"
)

func AddFav( id, pid int64 ) (*AccountEx, error) {
	a := &AccountEx{}
	has, e := engine_account.Table("account_ex").Where("account_id = ?",id).Get(a)
	if e != nil {
		return nil, e
	} else if !has {
		return nil, errors.New("account not found")
	}
	if a.FavPosts == nil {
		a.FavPosts = []int64{}
	}
	a.FavPosts = append(a.FavPosts, pid)
	if _, e =  engine_account.Table("account_ex").Where("account_id = ?",id).AllCols().Update( a ); e != nil {
		return nil, e
	}
	return a, nil
}

func RemoveFav( id, pid int64 ) (*AccountEx, error) {
	a := &AccountEx{}
	has, e := engine_account.Table("account_ex").Where("account_id = ?",id).Get(a)
	if e != nil {
		return nil, e
	} else if !has {
		return nil, errors.New("account not found")
	}

	// 从列表中移除
	for i, f := range a.FavPosts {
		if f == pid {
			a.FavPosts = append( a.FavPosts[:i], a.FavPosts[i+1:]... )
			break
		}
	}

	if _, e = engine_account.Table("account_ex").Where("account_id = ?",id).AllCols().Update( a ); e != nil {
		return nil, e
	}
	return a, nil
}

func IsPostFav( id, pid int64 ) ( isFav bool, e error ) {
	a := &AccountEx{}
	has, e := engine_account.Table("account_ex").Where("account_id = ?",id).Get(a)
	if e != nil {
		return false, e
	} else if !has {
		return false, errors.New("account not found")
	}

	// 从列表中移除
	for _, f := range a.FavPosts {
		if f == pid {
			return true, nil
		}
	}
	return false, nil
}

func UpdateFav( id int64, mdFavList []int64 ) (*AccountEx, error) {
	a := &AccountEx{}
	has, e := engine_account.Table("account_ex").Where("account_id = ?",id).Get(a)
	if e != nil {
		return nil, e
	} else if !has {
		return nil, errors.New("account not found")
	}
	a.FavPosts = mdFavList
	if _, e = engine_account.Table("account_ex").Where("account_id = ?",id).AllCols().Update( a ); e != nil {
		return nil, e
	}
	return a, nil
}