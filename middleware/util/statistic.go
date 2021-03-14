// 记录访问中间件辅助函数
// TODO: 进行ip地址和访问路由的的统计
package middlewareUtil

import (
	"encoding/json"
	"github.com/kpango/glg"
	"io/ioutil"
	"net/http"
	"strings"
)

// GetIP gets a requests IP address by reading off the forwarded-for
// header (for proxies) and falls back to use the remote address.
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

type IPInfo struct {
	Data struct {
		Area      string      `json:"area"`
		Country   string      `json:"country"`
		IspID     string      `json:"isp_id"`
		QueryIP   string      `json:"queryIp"`
		City      string      `json:"city"`
		IP        string      `json:"ip"`
		Isp       string      `json:"isp"`
		County    string      `json:"county"`
		RegionID  string      `json:"region_id"`
		AreaID    string      `json:"area_id"`
		CountyID  interface{} `json:"county_id"`
		Region    string      `json:"region"`
		CountryID string      `json:"country_id"`
		CityID    string      `json:"city_id"`
	} `json:"data"`
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

func CheckIPInfo( r *http.Request ) {
	AddrStr := GetIP( r )
	ipp := strings.Split(AddrStr,":")
	ipStr := ipp[0]
	glg.Info( "[AccessRecord] "+ ipStr )
	if resp, e := http.Get("http://ip.taobao.com/outGetIpInfo?ip="+ipStr+"&&accessKey=alibaba-inc"); e!=nil {
		glg.Error("in CheckIPInfo http.Get: ", e )
	} else {
		ipi := IPInfo{}
		if s, e := ioutil.ReadAll(resp.Body); e != nil {
			glg.Error("in CheckIPInfo ioutil.ReadAll: ", e )
		} else {
			json.Unmarshal( s, &ipi )
			// glg.Info( ipi )
			// TODO: 使用缓存记录访问ip的情况。
		}
	}
}
