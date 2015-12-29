package martini

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

/* 测试辅助方法*/
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("期望 %v (类型%v) - 为 %v (类型%v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("不期望 %v (类型%v) - 为 %v (类型%v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func Test_New(t *testing.T) {
	m := New()
	fmt.Printf("[KuuYee]====> m=%+v\n", m)
	if m == nil {
		t.Error("martini.New() 不能返回nil")
	}
}

func Test_Martini_RunOnAddr(t *testing.T) {
	go New().RunOnAddr("127.0.0.1:8080")
}

func Test_Martini_ServeHTTP(t *testing.T) {
	result := ""
	response := httptest.NewRecorder()

	m := New()
	m.Use(func(c Context) {
		result += "foo"
		c.Next()
		result += "ban"
	})
	m.Use(func(c Context) {
		result += "bar"
		c.Next()
		result += "baz"
	})
	m.Action(func(res http.ResponseWriter, req *http.Request) {
		result += "bat"
		res.WriteHeader(http.StatusBadRequest)
	})
	m.ServeHTTP(response, (*http.Request)(nil))
	expect(t, result, "foobarbatbazban")
	expect(t, response.Code, http.StatusBadRequest)
}

func Test_Martini_Handlers(t *testing.T) {
	result := ""
	response := httptest.NewRecorder()

	batman := func(c Context) {
		result += "batman!"
	}
	m := New()
	m.Use(func(c Context) {
		result += "foo"
		c.Next()
		result += "ban"
	})
	m.Handlers(
		batman,
		batman,
		batman,
	)
	m.Action(func(res http.ResponseWriter, req *http.Request) {
		result += "bat"
		res.WriteHeader(http.StatusBadRequest)
	})
	m.ServeHTTP(response, (*http.Request)(nil))

	expect(t, result, "batman!batman!batman!bat")
	expect(t, response.Code, http.StatusBadRequest)
}

func Test_Martini_EarlyWrite(t *testing.T) {
	result := ""
	response := httptest.NewRecorder()

	m := New()
	m.Use(func(res http.ResponseWriter) {
		result += "foobar"
		res.Write([]byte("Hello world"))
	})
	m.Use(func() {
		result += "bat"
	})
	m.Action(func(res http.ResponseWriter) {
		result += "baz"
		res.WriteHeader(http.StatusBadRequest)
	})

	m.ServeHTTP(response, (*http.Request)(nil))
	expect(t, result, "foobar")
	expect(t, response.Code, http.StatusOK)
}

func Test_Martini_Written(t *testing.T) {
	response := httptest.NewRecorder()

	m := New()
	m.Handlers(func(res http.ResponseWriter) {
		res.WriteHeader(http.StatusOK)
	})
	ctx := m.createContext(response, (*http.Request)(nil))
	expect(t, ctx.Written(), false)
	ctx.run()
	expect(t, ctx.Written(), true)
}
