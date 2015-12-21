package martini

import (
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
	return m
}

// 如果你想控制你自己的Http Server，那么ServeHTTP是Martini实例的HTTP进入点
func (m *Martini) ServeHTTP(res http.ResponseWriter, req *http.Request) {

}

// 使用提供的host和port运行HTTP Server
func (m *Martini) RunOnAddr(addr string) {
	// TODO:或许应该直接调用http.ListenAndServer，那样就可以保存martini以便后续使用
	// 这样或许能够改善测试，可以接受一个自定义的host和port传入

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

// Handler是任何可以调用的函数。 Martini试图通过处理器的参数列表注入服务
// 如果参数无法注入，Martini将引发panic
type Handler interface{}
