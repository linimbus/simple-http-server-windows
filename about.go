package main

import (
	"os/exec"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func OpenBrowserWeb(url string) {
	cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	err := cmd.Run()
	if err != nil {
		logs.Error("run cmd fail, %s", err.Error())
	}
}

func AboutAction() {
	var ok *walk.PushButton
	var about *walk.Dialog
	var err error

	_, err = Dialog{
		AssignTo:      &about,
		Title:         "Sponsor",
		Icon:          walk.IconInformation(),
		MinSize:       Size{Width: 200, Height: 200},
		DefaultButton: &ok,
		Layout:        VBox{},
		Children: []Widget{
			Composite{
				Layout: HBox{},
				Children: []Widget{
					TextLabel{
						MinSize: Size{Width: 180, Height: 180},
						Text:    "Thank you for using my software. If you are pleased with it, you can donate through the link below. Thank you very much for your support.",
					},
				},
			},
			PushButton{
				Text: "paypal.me",
				OnClicked: func() {
					OpenBrowserWeb("https://paypal.me/lixiangyun")
				},
			},
			PushButton{
				Text:      "OK",
				OnClicked: func() { about.Cancel() },
			},
		},
	}.Run(mainWindow)

	if err != nil {
		logs.Error(err.Error())
	}
}
