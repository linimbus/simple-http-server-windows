package main

import (
	"fmt"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var statusFlow *walk.StatusBarItem
var sessionFlow *walk.StatusBarItem
var requestFlow *walk.StatusBarItem

func StatusUpdate(flow int64, session int64, request int64) {
	if statusFlow != nil {
		statusFlow.SetText(fmt.Sprintf("DataFlow: %s", ByteView(flow)))
		sessionFlow.SetText(fmt.Sprintf("Session: %d", session))
		requestFlow.SetText(fmt.Sprintf("Request: %d", request))
	}
}

func StatusBarInit() []StatusBarItem {
	return []StatusBarItem{
		{
			AssignTo: &statusFlow,
			Text:     "",
			Width:    120,
		},
		{
			AssignTo: &sessionFlow,
			Text:     "",
			Width:    80,
		},
		{
			AssignTo: &requestFlow,
			Text:     "",
			Width:    80,
		},
	}
}
