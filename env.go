package martini

import (
	"os"
)

// Envs
const (
	Dev  string = "开发环境"
	Test string = "测试环境"
	Prod string = "生产环境"
)

var Env = Dev
var Root string

func setENV(e string) {
	if len(e) > 0 {
		Env = e
	}
}

func init() {
	// 通过环境变量读取EVN值，如果没有这种则采用默认值Dev
	setENV(os.Getenv("MARTINI_ENV"))
	var err error
	Root, err = os.Getwd()
	if err != nil {
		panic(err)
	}
}
