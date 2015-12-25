package martini

import (
	_ "fmt"
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/codegangsta/inject"
)

// Martini是一个顶层Web框架，inject.Injector方法可以调用全局能够映射的服务
type Martini struct {
	inject.Injector             // 依赖注入实例
	handlers        []Handler   //处理器列表
	action          Handler     //处理器执行某种Action
	logger          *log.Logger //全局日志器
}

// 创建一个最直接的Martini实例，这可以让你完全控制使用的中间件
func New() *Martini {
	m := &Martini{
		Injector: inject.New(),
		action:   func() {},
		logger:   log.New(os.Stdout, "[martini-learn] ", 0),
	}
	m.Map(m.logger)
	m.Map(defaultReturnHandler())
	return m
}

// 如果你想控制你自己的Http Server，那么ServeHTTP是Martini实例的HTTP进入点
func (m *Martini) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	m.createContext(res, req).run()
}

// 使用提供的host和port运行HTTP Server
func (m *Martini) RunOnAddr(addr string) {
	// TODO:或许应该直接调用http.ListenAndServer，那样就可以保存martini以便后续使用
	// 这样或许能够改善测试，可以接受一个自定义的host和port传入

	//logger := m.Injector.Get(reflect.TypeOf(m.logger)).Interface().(*log.Logger)
	logger := m.Injector.Get(reflect.TypeOf(m.logger)).Interface().(*log.Logger)
	logger.Printf("服务器监听在 %s (%s)\n", addr, Env)
	logger.Fatalln(http.ListenAndServe(addr, m))
}

// 运行http Server，监听端口读取系统环境变量 os.GetEvn("PORT") 默认监听3000端口
func (m *Martini) Run() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}
	host := os.Getenv("HOST")
	m.RunOnAddr(host + ":" + port)
}

func (m *Martini) createContext(res http.ResponseWriter, req *http.Request) *context {
	c := &context{
		inject.New(),
		m.handlers,
		m.action,
		NewResponseWriter(res),
		0,
	}
	c.SetParent(m)
	c.MapTo(c, (*Context)(nil))
	c.MapTo(c.rw, (*http.ResponseWriter)(nil))
	c.Map(req)
	return c
}

// Handler是任何可以调用的函数。 Martini试图通过处理器的参数列表注入服务
// 如果参数无法注入，Martini将引发panic
type Handler interface{}

func validateHandler(handler Handler) {
	if reflect.TypeOf(handler).Kind() != reflect.Func {
		panic("[KuuYee]===> Martini处理器必须是可调用的函数！")
	}
}

// Context 呈现一个请求上下文。服务可以通过这个借口映射请求层
type Context interface {
	inject.Injector
	// Next 是可选函数，可以让中间件处理器在其它处理器之后执行。
	// 这非常有用，对于一些必须在http请求之后完成的操作
	Next()
	// 写入返回，无论这个上下文是否有Response
	Written() bool
}

type context struct {
	inject.Injector
	handlers []Handler
	action   Handler
	rw       ResponseWriter
	index    int
}

func (c *context) Next() {
	c.index += 1
	c.run()
}

func (c *context) Written() bool {
	return c.rw.Written()
}

func (c *context) handler() Handler {
	// 如果c.index 小于handlers数组的长度，表示索引还没执行指向最后一个handler
	if c.index < len(c.handlers) {
		return c.handlers[c.index]
	}

	// 如果c.index等于handlers数组的长度，表示已经执行完所有的handler,下面开始aciton
	if c.index == len(c.handlers) {
		return c.action
	}
	panic("错误的上下文处理器index")
}

// Action设置一个处理器，在所有的中间件处理执行完后执行，这个被用来在martini.Classic()
// 中设置martini.Router
func (m *Martini) Action(handler Handler) {
	validateHandler(handler)
	m.action = handler
}

// Use 在处理中心中加入一个中间件处理器，如果处理不能被调用则引发一个panic，处理器的
// 调用顺序是按照加入的顺序执行的。
func (m *Martini) Use(handler Handler) {
	validateHandler(handler)

	m.handlers = append(m.handlers, handler)
}

func (c *context) run() {
	for c.index <= len(c.handlers) {
		_, err := c.Invoke(c.handler())
		if err != nil {
			panic(err)
		}
		c.index += 1

		// 检查是否某个处理器已经写入了Response，如果有那么直接返回
		if c.Written() {
			return
		}
	}
}
