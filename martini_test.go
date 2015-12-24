package martini

import (
	"fmt"
	"testing"
)

func Test_New(t *testing.T) {
	m := New()
	fmt.Printf("[KuuYee]====> m=%+v", m)
	if m == nil {
		t.Error("martini.New() 不能返回nil")
	}
}

func Test_Martini_RunOnAddr(t *testing.T) {
	m := New()
	go m.RunOnAddr("127.0.0.1:8080")
}
