package martini

import (
	"testing"
)

func Test_Root(t *testing.T) {
	if len(Root) == 0 {
		t.Errorf("root路径必须被设置")
	}
}
