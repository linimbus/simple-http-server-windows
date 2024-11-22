package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var processBar *walk.ProgressBar

func ProcessWidget() []Widget {
	return []Widget{
		ProgressBar{
			AssignTo: &processBar,
			MaxValue: 1000,
			MinValue: 0,
			MinSize:  Size{Height: 5},
			MaxSize:  Size{Height: 5},
		},
	}
}

func ProcessUpdate(value int) {
	processBar.SetValue(value)
}
