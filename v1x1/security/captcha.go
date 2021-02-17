package security

import (
	"../../cc"
	"../../cc/err"
	"../../cc/err_code"
	captcha "github.com/base64Captcha"
	"github.com/kpango/glg"
	"net/http"
	"strconv"
)

var store = captcha.DefaultMemStore

func newDriver() *captcha.DriverString {
	driver := new(captcha.DriverString)
	driver.Height = 44
	driver.Width = 120
	driver.NoiseCount = 0
	driver.ShowLineOptions = captcha.OptionShowSineLine // | captcha.OptionShowSlimeLine // | captcha.OptionShowHollowLine
	driver.Length = 4
	driver.Source = "1234567890qwertyuipkjhgfdsazxcvbnm"
	driver.Fonts = []string{"wqy-microhei.ttc"}
	return driver
}

func GenerateCaptcha(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r  != nil {
			_ = glg.Error(r)
		}
	}()
	cid := r.FormValue("cid")
	time := r.FormValue("time")
	t, e := strconv.Atoi( time )
	err.Assert( e )

	if t >= 7 {
		cc.HttpReturn( &w, "too much times", err_code.ERR_SECURITY, "", cc.MakeHER200 )
		return
	}

	var driver = newDriver().ConvertFonts()
	c := captcha.NewCaptcha(driver, store)

	_, content, answer := c.Driver.GenerateIdQuestionAnswer()
	// id := "captcha:yufei"
	item, _ := c.Driver.DrawCaptcha(content)
	c.Store.Set(cid, answer)
	glg.Log( "new captcha generated! answer is [" + answer + "]")
	item.WriteTo(w)
}

func CaptchaVerify(code, cid string) bool {
	return store.Verify(cid, code, true)
}