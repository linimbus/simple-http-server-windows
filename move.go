package main

import (
	"os"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

func MoveAction(from walk.Form, activeFunc func(isNew bool), cancelFunc func()) {
	time.Sleep(200 * time.Millisecond)

	var dlg *walk.Dialog
	var destinationDir *walk.LineEdit
	var acceptPB, cancelPB *walk.PushButton
	var olderCheck, newerCheck *walk.CheckBox

	_, err := Dialog{
		AssignTo:      &dlg,
		Title:         "Move Options",
		Icon:          walk.IconInformation(),
		MinSize:       Size{Width: 450, Height: 150},
		Size:          Size{Width: 450, Height: 150},
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		Layout:        VBox{},
		Children: []Widget{
			VSpacer{},
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					Label{
						Text: "Destination Directory: ",
					},
					LineEdit{
						AssignTo: &destinationDir,
						Text:     ConfigGet().DestinationDir,
						OnEditingFinished: func() {
							dir := destinationDir.Text()
							if dir != "" {
								stat, err := os.Stat(dir)
								if err != nil {
									ErrorBoxAction(mainWindow, "The move to destination directory is not exist")
									destinationDir.SetText(ConfigGet().DestinationDir)
									return
								}
								if !stat.IsDir() {
									ErrorBoxAction(mainWindow, "The move to destination directory is not directory")
									destinationDir.SetText(ConfigGet().DestinationDir)
									return
								}
							}
							DestinationDirDirSave(dir)
						},
					},
					PushButton{
						MaxSize: Size{Width: 20},
						Text:    "...",
						OnClicked: func() {
							dlgDir := new(walk.FileDialog)
							dlgDir.FilePath = ConfigGet().DestinationDir
							dlgDir.Flags = win.OFN_EXPLORER
							dlgDir.Title = "Please select a folder as destination directory"

							exist, err := dlgDir.ShowBrowseFolder(mainWindow)
							if err != nil {
								logs.Error(err.Error())
								return
							}
							if exist {
								logs.Info("select %s as destination directory", dlgDir.FilePath)
								destinationDir.SetText(dlgDir.FilePath)
								DestinationDirDirSave(dlgDir.FilePath)
							}
						},
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					CheckBox{
						AssignTo: &olderCheck,
						Text:     "Move Older File",
						Checked:  true,
						OnCheckedChanged: func() {
							newerCheck.SetChecked(!olderCheck.Checked())
						},
					},
					CheckBox{
						AssignTo: &newerCheck,
						Text:     "Move Newer File",
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
							if ConfigGet().DestinationDir == "" {
								ErrorBoxAction(mainWindow, "The move to destination directory is empty")
								return
							}
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
