package main

import (
	"github.com/astaxie/beego/logs"
)

func main() {
	err := FileInit()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	err = LogInit()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	err = BoxInit()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	err = IconInit()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	err = ConfigInit()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	mainWindows()
}
