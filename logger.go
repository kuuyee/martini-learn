package martini

import (
	"log"
	"net/http"
	"time"
)

// Logger 返回一个中间件处理器，用来记录请入输入和响应输出。
func Logger() Handler {
	return func(res http.ResponseWriter, req *http.Request, c Context, log *log.Logger) {
		start := time.Now()

		// 获取访问的IP地址
		addr := req.Header.Get("X-Real-IP")
		if addr == "" {
			// 获取被转发了的访问IP地址
			addr = req.Header.Get("X-Forwarded-For")
			if addr == "" {
				addr = req.RemoteAddr
			}
		}

		log.Printf("[KuuYee] 开始于 %s %s for %s", req.Method, req.URL.Path, addr)

		rw := res.(ResponseWriter)
		c.Next()

		log.Printf("[KuuYee] 完成与 %v %s in %v\n", rw.Status(), http.StatusText(rw.Status()), time.Since(start))
	}
}
