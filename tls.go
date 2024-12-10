package main

import (
	"fmt"

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
						Text:     ConfigGet().HttpsInfo.CA,
						VScroll:  true,
					},
					Label{
						Text: "TLS Cert*: ",
					},
					TextEdit{
						AssignTo: &tlsCert,
						MinSize:  Size{Height: 100},
						Text:     ConfigGet().HttpsInfo.Cert,
						VScroll:  true,
					},
					Label{
						Text: "TLS Key*: ",
					},
					TextEdit{
						AssignTo: &tlsKey,
						MinSize:  Size{Height: 100},
						Text:     ConfigGet().HttpsInfo.Key,
						VScroll:  true,
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
							if tlsCert.Text() != "" || tlsKey.Text() != "" {
								_, err := CreateTlsConfig(tlsCert.Text(), tlsKey.Text())
								if err != nil {
									ErrorBoxAction(mainWindow, fmt.Sprintf("TLS Cert or Key maybe invalid! %s", err.Error()))
									return
								}
							}

							err := HttpsInfoSave(TlsInfo{
								CA:   tlsCA.Text(),
								Cert: tlsCert.Text(),
								Key:  tlsKey.Text(),
							})
							if err != nil {
								ErrorBoxAction(mainWindow, err.Error())
								return
							}

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
