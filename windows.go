package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var mainWindow *walk.MainWindow

var mainWindowWidth = 800
var mainWindowHeight = 500

func MenuBarInit() []MenuItem {
	return []MenuItem{
		Action{
			Text: "Runlog",
			OnTriggered: func() {
				OpenBrowserWeb(RunlogDirGet())
			},
		},
		Action{
			Text: "Sponsor",
			OnTriggered: func() {
				AboutAction()
			},
		},
	}
}

func mainWindows() {
	CapSignal(CloseWindows)
	cnt, err := MainWindow{
		Title:     "Duplicate File Cleaner " + VersionGet(),
		Icon:      ICON_Main,
		AssignTo:  &mainWindow,
		MinSize:   Size{Width: mainWindowWidth, Height: mainWindowHeight},
		Size:      Size{Width: mainWindowWidth, Height: mainWindowHeight},
		Layout:    VBox{Margins: Margins{Top: 5, Bottom: 5, Left: 5, Right: 5}},
		MenuItems: MenuBarInit(),
		Children: []Widget{
			Composite{
				Layout:   VBox{},
				Children: ConsoleWidget(),
			},
			Composite{
				Layout:   VBox{Margins: Margins{Top: 0, Bottom: 0, Left: 10, Right: 10}},
				Children: ProcessWidget(),
			},
			Composite{
				Layout:   VBox{},
				Children: TableWidget(),
			},
			Composite{
				Layout:   HBox{},
				Children: ActiveWidget(),
			},
		},
	}.Run()

	if err != nil {
		logs.Error(err.Error())
	} else {
		logs.Info("main windows exit %d", cnt)
	}

	if err := recover(); err != nil {
		logs.Error(err)
	}

	CloseWindows()
}

func CloseWindows() {
	if mainWindow != nil {
		mainWindow.Close()
		mainWindow = nil
	}
}
