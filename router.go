package martini

import (
	"net/http"
	"regexp"
	"sync"
)

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
