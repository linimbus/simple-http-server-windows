package main

import (
	"fmt"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var requestFlow *walk.StatusBarItem

func StatusRequestUpdate(cnt int64) {
	if requestFlow != nil {
		requestFlow.SetText(fmt.Sprintf(" REQUEST: %d", cnt))
	}
}

func StatusBarInit() []StatusBarItem {
	return []StatusBarItem{
		{
			AssignTo: &requestFlow,
			Text:     " REQUEST: 0",
			Width:    120,
		},
	}
}
