package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func TlsAction() {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton
	var tlsCA, tlsCert, tlsKey *walk.TextEdit

	_, err := Dialog{
		AssignTo:      &dlg,
		Title:         "TLS Certificate Edit",
		Icon:          walk.IconInformation(),
		DefaultButton: &acceptPB,
		CancelButton:  &acceptPB,
		Size:          Size{Width: 450, Height: 300},
		MinSize:       Size{Width: 450, Height: 300},
		Layout:        VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "TLS CA: ",
					},
					TextEdit{
						AssignTo: &tlsCA,
						MinSize:  Size{Height: 100},

						Text: "",
					},
					Label{
						Text: "TLS Cert: ",
					},
					TextEdit{
						AssignTo: &tlsCert,
						MinSize:  Size{Height: 100},

						Text: "",
					},
					Label{
						Text: "TLS Key: ",
					},
					TextEdit{
						AssignTo: &tlsKey,
						MinSize:  Size{Height: 100},
						Text:     "",
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "OK",
						OnClicked: func() {
							dlg.Accept()
						},
					},
					HSpacer{},
					PushButton{
						AssignTo: &cancelPB,
						Text:     "Cancel",
						OnClicked: func() {
							dlg.Cancel()
						},
					},
					HSpacer{},
				},
			},
		},
	}.Run(mainWindow)

	if err != nil {
		logs.Error(err.Error())
	}
}
