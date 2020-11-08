package rate_limiter

import (
	"goproxywork/middleware"
	"goproxywork/proxy"
	"log"
	"net/http"
	"net/url"
)

var addr = "127.0.0.1:8002"

// 熔断方案

func RateLimitServer() {
	coreFunc := func(c *middleware.SliceRouterContext) http.Handler {
		rs1 := "http://127.0.0.1:2003/base"
		url1, err := url.Parse(rs1)
		if err != nil {
			log.Println(err)
		}
		rs2 := "hhtp://127.0.0.1:2004/base"
		url2, err := url.Parse(rs2)
		if err != nil {
			log.Println(err)
		}
		urls := []*url.URL{url1, url2}
		return proxy.NewMultipleHostsReverseProxy(c, urls)
	}
	sliceRoute := middleware.NewSliceRouter()
	sliceRoute.Group("/").Use(middleware.RateLimiter())
	routerHandler := middleware.NewSliceRouterHandler(coreFunc, sliceRoute)
	_ = http.ListenAndServe(addr, routerHandler)

}
