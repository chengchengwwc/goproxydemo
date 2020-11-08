package middleware

import (
	"context"
	"math"
	"net/http"
	"strings"
)

// 通用中间件

const abortIndex int8 = math.MaxInt8 / 2 // 最多63个中间件

type HandleFunc func(*SliceRouterContext)

type SliceRouter struct {
	groups []*SliceGroup
}

type SliceGroup struct {
	*SliceRouter
	path    string
	handles []HandleFunc
}

// rounter 上下文
type SliceRouterContext struct {
	Rw  http.ResponseWriter
	Req *http.Request
	Ctx context.Context
	*SliceGroup
	index int8
}

func newSliceRouterContext(rw http.ResponseWriter, req *http.Request, r *SliceRouter) *SliceRouterContext {
	newSliceGroup := &SliceGroup{}
	matchUrlLen := 0
	for _, group := range r.groups {
		if strings.HasPrefix(req.RequestURI, group.path) {
			pathLen := len(group.path)
			if pathLen > matchUrlLen {
				matchUrlLen = pathLen
				*newSliceGroup = *group //浅拷贝数组指针
			}
		}
	}
	c := &SliceRouterContext{Rw: rw, Req: req, SliceGroup: newSliceGroup, Ctx: req.Context()}
	c.Reset()
	return c
}

func (c *SliceRouterContext) Get(key interface{}) interface{} {
	return c.Ctx.Value(key)
}

func (c *SliceRouterContext) Set(key, val interface{}) {
	c.Ctx = context.WithValue(c.Ctx, key, val)
}

type SliceRouterHandler struct {
	coreFunc func(routerContext *SliceRouterContext) http.Handler
	router   *SliceRouter
}

func (c *SliceRouterContext) Reset() {
	c.index = -1
}

func (c *SliceRouterContext) Next() {
	c.index++
	for c.index < int8(len(c.handles)) {
		c.handles[c.index](c)
		c.index++
	}
}

func (c *SliceRouterContext) Abort() {
	c.index = abortIndex
}

func (c *SliceRouterContext) IsAborted() bool {
	return c.index >= abortIndex
}

func (w *SliceRouterHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	c := newSliceRouterContext(rw, req, w.router)
	if w.coreFunc != nil {
		c.handles = append(c.handles, func(c *SliceRouterContext) {
			w.coreFunc(c).ServeHTTP(rw, req)
		})
	}
	c.Reset()
	c.Next()
}

func NewSliceRouterHandler(coreFunc func(routerContext *SliceRouterContext) http.Handler, router *SliceRouter) *SliceRouterHandler {
	return &SliceRouterHandler{
		coreFunc: coreFunc,
		router:   router,
	}
}

// 构造router
func NewSliceRouter() *SliceRouter {
	return &SliceRouter{}
}

// 创建Group
func (g *SliceRouter) Group(path string) *SliceGroup {
	return &SliceGroup{
		SliceRouter: g,
		path:        path,
	}
}

// 构建回调方法
func (g *SliceGroup) Use(middlewares ...HandleFunc) *SliceGroup {
	// 用以注册方法
	g.handles = append(g.handles, middlewares...)
	existsFlag := false
	for _, oldGroup := range g.SliceRouter.groups {
		if oldGroup == g {
			existsFlag = true
		}
	}
	if !existsFlag {
		// 添加
		g.SliceRouter.groups = append(g.SliceRouter.groups, g)
	}
	return g
}
