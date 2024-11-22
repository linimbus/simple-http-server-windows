package main

import (
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func DeleteAction(from walk.Form, activeFunc func(isNew bool), cancelFunc func()) {
	time.Sleep(200 * time.Millisecond)

	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton
	var olderCheck, newerCheck *walk.CheckBox

	_, err := Dialog{
		AssignTo:      &dlg,
		Title:         "Delete Options",
		Icon:          walk.IconInformation(),
		MinSize:       Size{Width: 210, Height: 150},
		Size:          Size{Width: 210, Height: 150},
		MaxSize:       Size{Width: 310, Height: 210},
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		Layout:        VBox{},
		Children: []Widget{
			VSpacer{},

			Composite{
				Layout: HBox{},
				Children: []Widget{
					CheckBox{
						AssignTo: &olderCheck,
						Text:     "Delete Older File",
						Checked:  true,
						OnCheckedChanged: func() {
							newerCheck.SetChecked(!olderCheck.Checked())
						},
					},
					CheckBox{
						AssignTo: &newerCheck,
						Text:     "Delete Newer File",
						Checked:  false,
						OnCheckedChanged: func() {
							olderCheck.SetChecked(!newerCheck.Checked())
						},
					},
				},
			},
			VSpacer{},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						AssignTo: &acceptPB,
						Text:     "OK",
						OnClicked: func() {
							dlg.Accept()
							go activeFunc(newerCheck.Checked())
						},
					},
					PushButton{
						AssignTo: &cancelPB,
						Text:     "Cancel",
						OnClicked: func() {
							dlg.Accept()
							go cancelFunc()
						},
					},
				},
			},
		},
	}.Run(from)

	if err != nil {
		logs.Error(err.Error())
	}
}
