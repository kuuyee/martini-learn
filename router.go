package martini

import (
	"net/http"
	"regexp"
	"sync"
)

// Params is a map of name/value pairs for named routes. An instance of martini.
// Params is available to be injected into any route handler.
type Params map[string]string

// Router 是Martini的路由接口,Supports HTTP verbs, stacked handlers, and dependency injection.
type Router interface {
	Routes

	Group(string, func(Router), ...Handler)
	Get(string, ...Handler) Route
	Post(string, ...Handler) Route
	AddRoute(string, string, ...Handler) Route
	NotFound(...Handler)
	Handle(http.ResponseWriter, *http.Request, Context)
}

type router struct {
	routes     []*route
	notFounds  []Handler
	groups     []group
	routesLock sync.RWMutex
}

type group struct {
	pattern  string
	handlers []Handler
}

func NewRouter() *router {
	return &router{notFounds: []Handler{http.NotFound}, groups: make([]group, 0)}
}

// Route 接口用来呈现Martini的路由层
type Route interface {
	// URLWith 返回一个用给定字符串参数渲染的路由URL
	URLWith([]string) string
	// Name 设置路由的名称
	Name(string)
	// GetName返回路由的名字
	GetName() string
	// Pattern 返回路由的pattern
	Pattern() string
	// Method 返回路由的Method
	Method() string
}

type route struct {
	method   string
	regex    *regexp.Regexp
	handlers []Handler
	pattern  string
	name     string
}

// Routes 是Martini路由层的辅助服务
type Routes interface {
	// URLFor 从一个给定的路由返回一个渲染过的URL，
	// Optional params can be passed to fulfill named parameters in the route.
	URLFor(name string, params ...interface{}) string
	// MethodsFor returns an array of methods available for the path
	MethodsFor(path string) []string
	// All returns an array with all the routes in the router.
	All() []Route
}
