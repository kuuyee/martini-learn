package martini

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Logger(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()

	m := New()
	// 取代testing的log
	m.Map(log.New(buff, "[KuuYee] ", 0))
	m.Use(Logger())
	m.Use(func(res http.ResponseWriter) {
		res.WriteHeader(http.StatusFound)
	})

	req, err := http.NewRequest("GET", "http://localhost:3001/foobar", nil)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("[KuuYee]====> req.Code")
	m.ServeHTTP(recorder, req)
	//expect(t, recorder.Code, http.StatusNotFound)
	//expect(t, len(buff.String()), 0)
}
