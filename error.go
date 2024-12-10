package main

import (
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func boxAction(from walk.Form, title string, icon *walk.Icon, message string) {
	var dlg *walk.Dialog
	var cancelPB *walk.PushButton

	_, err := Dialog{
		AssignTo:     &dlg,
		Title:        title,
		Icon:         icon,
		MinSize:      Size{Width: 210, Height: 120},
		CancelButton: &cancelPB,
		Layout:       VBox{},
		Children: []Widget{
			TextLabel{
				Text:    message,
				MinSize: Size{Width: 180, Height: 100},
			},
			PushButton{
				AssignTo: &cancelPB,
				Text:     "OK",
				OnClicked: func() {
					dlg.Accept()
				},
			},
		},
	}.Run(from)

	if err != nil {
		logs.Error(err.Error())
	}
}

func ErrorBoxAction(form walk.Form, message string) {
	time.Sleep(200 * time.Millisecond)
	boxAction(form, "Error", walk.IconError(), message)
}

func InfoBoxAction(form walk.Form, message string) {
	time.Sleep(200 * time.Millisecond)
	boxAction(form, "Info", walk.IconInformation(), message)
}

func ConfirmBoxAction(form walk.Form, message string) {
	time.Sleep(200 * time.Millisecond)
	boxAction(form, "Confirm", walk.IconWarning(), message)
}
