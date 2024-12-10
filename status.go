package main

import (
	"fmt"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var dataFlow *walk.StatusBarItem
var sessionFlow *walk.StatusBarItem
var requestFlow *walk.StatusBarItem

func StatusRequestUpdate(cnt int64) {
	if requestFlow != nil {
		requestFlow.SetText(fmt.Sprintf("REQUEST: %d", cnt))
	}
}

func StatusSessionUpdate(cnt int64) {
	if sessionFlow != nil {
		sessionFlow.SetText(fmt.Sprintf("SESSION: %d", cnt))
	}
}

func StatusFlowUpdate(flow int64) {
	if dataFlow != nil {
		dataFlow.SetText(fmt.Sprintf("DATAFLOW: %s", ByteView(flow)))
	}
}

func StatusBarInit() []StatusBarItem {
	return []StatusBarItem{
		{
			AssignTo: &dataFlow,
			Text:     "",
			Width:    120,
		},
		{
			AssignTo: &sessionFlow,
			Text:     "",
			Width:    120,
		},
		{
			AssignTo: &requestFlow,
			Text:     "",
			Width:    120,
		},
	}
}
