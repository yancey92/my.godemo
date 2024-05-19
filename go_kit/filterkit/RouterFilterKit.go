package filterkit

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"time"
)

var apiStatistics = make(map[string]*apiInfo)

type apiInfo struct {
	url       string
	startTime time.Time
	elapsed   time.Duration
	params    string
}

func BeforeRouterFilter(ctx *context.Context) {
	ctx.Request.ParseForm()
	url := ctx.Input.URL()
	params := fmt.Sprint(ctx.Request.Form)
	now := time.Now()

	info := apiStatistics[url]
	if info == nil {
		apiStatistics[url] = &apiInfo{
			elapsed: 0,
		}
	}

	info = apiStatistics[url]
	info.url = url
	info.startTime = now
	info.params = params
}

func FinishRouterFilter(ctx *context.Context) {
	url := ctx.Input.URL()
	info := apiStatistics[url]
	info.elapsed = time.Since(info.startTime)
}

func InItBeegoFilter(pattern string) {
	beego.InsertFilter(pattern, beego.BeforeRouter, BeforeRouterFilter, false)
	beego.InsertFilter(pattern, beego.FinishRouter, FinishRouterFilter, false)
}
