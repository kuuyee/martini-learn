package martini

import (
	"net/http"
)

// ResponseWriter 是http.ResponseWriter的封装，提供一些response的扩展信息
// 推荐其它中间件处理使用这个方式来封装responsewriter
type ResponseWriter interface {
	http.ResponseWriter
	//http.Flusher
	//http.Hijacker
	// Status 返回response的状态码，如果没有写入response则返回0
	Status() int
	// Written 返回一个布尔值表示ResponseWriter是否被写入
	Written() bool
	// Size 返回respose主体内容的大小
	Size() int
	// Before允许在ResponseWriter被写入之前调用。这对于设置header或其它操作
	// 必须在写入ResponseWriter之前发生。
	Before(BeforeFunc)
}

// BeforeFunc 是一个函数，在ResponseWriter被写入前调用
type BeforeFunc func(ResponseWriter)

// NewResponseWriter 创建一个封装了http.ResponseWriter的ResponseWriter
func NewResponseWriter(rw http.ResponseWriter) ResponseWriter {
	return &responseWriter{rw, 0, 0, nil}
}

type responseWriter struct {
	http.ResponseWriter
	status      int
	size        int
	beforeFuncs []BeforeFunc
}

func (rw *responseWriter) WriteHeader(s int) {
	rw.callBefore()
	rw.ResponseWriter.WriteHeader(s)
	rw.status = s
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.Written() {
		// 如果WriteHeader还没有被调用过，那么这里的status应该设置为StatusOK
		rw.WriteHeader(http.StatusOK)
	}
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) Written() bool {
	return rw.status != 0
}

func (rw *responseWriter) Size() int {
	return rw.size
}

func (rw *responseWriter) Before(before BeforeFunc) {
	rw.beforeFuncs = append(rw.beforeFuncs, before)
}

func (rw *responseWriter) callBefore() {
	for i := len(rw.beforeFuncs) - 1; i >= 0; i-- {
		rw.beforeFuncs[i](rw)
	}
}
