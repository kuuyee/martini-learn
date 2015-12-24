package martini

import (
	"net/http"
	"reflect"

	"github.com/codegangsta/inject"
)

// ReturnHandler是一个服务，在路由处理返回后被调用。 ReturnHandler用来写入ResponseWriter
// 基于传递进来的参数。
type ReturnHandler func(Context, []reflect.Value)

func defaultReturnHandler() ReturnHandler {
	return func(ctx Context, vals []reflect.Value) {
		rv := ctx.Get(inject.InterfaceOf((*http.ResponseWriter)(nil)))
		res := rv.Interface().(http.ResponseWriter)
		var responseVal reflect.Value
		if len(vals) > 1 && vals[0].Kind() == reflect.Int {
			res.WriteHeader(int(vals[0].Int()))
			responseVal = vals[1]
		} else if len(vals) > 0 {
			responseVal = vals[0]
		}
		if canDeref(responseVal) {
			responseVal = responseVal.Elem()
		}
		if isByteSlice(responseVal) {
			res.Write(responseVal.Bytes())
		} else {
			res.Write([]byte(responseVal.String()))
		}
	}
}

// 判断是否为字节切片
func isByteSlice(val reflect.Value) bool {
	return val.Kind() == reflect.Slice && val.Type().Elem().Kind() == reflect.Uint8
}

// 判断是否为引用类型 比如接口或指针
func canDeref(val reflect.Value) bool {
	return val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr
}
