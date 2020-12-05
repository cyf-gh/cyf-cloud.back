package orm

type (
	InfoComponent struct {
		Markdown string
		LastUpdate string
	}
	HomeInfoModel struct {
		Name, Avatar string
		PostCount int
		// 通过特殊的标题名来进行匹配
		Info, Projects, Gears InfoComponent

		Id int64
		Level string
		Exp int64
	}
)

func post2InfoComponent( posts []Post ) ( ic InfoComponent ) {
	if len(posts) == 0 {
		ic = InfoComponent{
			Markdown:   "",
			LastUpdate: "",
		}
		return
	}
	ic = InfoComponent{
		Markdown:   posts[0].Text,
		LastUpdate: posts[0].Date,
	}
	if posts[0].IsPrivate {
		ic.Markdown = "PRIVATE"
	}
	return
}

// --------- 账户信息模块 ---------
func GetUserHomeInfo( id int64 ) ( him HomeInfoModel, e error ) {
	var (
		a *Account
		ex *AccountEx
		info, gears, projects []Post )

	if a, e = GetAccount( id ); e != nil { return }
	if ex, e = GetAccountEx( id ); e != nil { return }

	if info, e = GetInfoComponentPost( id, "MyInfo" ); e != nil { return }
	if gears, e = GetInfoComponentPost( id, "MyGears" ); e != nil { return }
	if projects, e = GetInfoComponentPost( id, "MyProjects" ); e != nil { return }

	him = HomeInfoModel{
		Name:      a.Name,
		Avatar:    ex.Avatar,
		PostCount: 0,
		Info:      post2InfoComponent(info),
		Projects:  post2InfoComponent(projects),
		Gears:     post2InfoComponent(gears),
		Id:        a.Id,
		Level:     ex.Level,
		Exp:       ex.Exp,
	}
	return
}