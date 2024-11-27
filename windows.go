package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var mainWindow *walk.MainWindow

var mainWindowWidth = 450
var mainWindowHeight = 200

func MenuBarInit() []MenuItem {
	return []MenuItem{
		Action{
			Text: "Runlog",
			OnTriggered: func() {
				OpenBrowserWeb(RunlogDirGet())
			},
		},
		Action{
			Text: "TLS Edit",
			OnTriggered: func() {
				TlsAction()
			},
		},
		Action{
			Text: "Users Edit",
			OnTriggered: func() {
				UsersAction()
			},
		},
		Action{
			Text: "Mini Windows",
			OnTriggered: func() {
				NotifyAction()
			},
		},
		Action{
			Text: "Sponsor",
			OnTriggered: func() {
				AboutAction()
			},
		},
		Action{
			Text: "Close Windows",
			OnTriggered: func() {
				CloseWindows()
			},
		},
	}
}

func ConsoleWidget() []Widget {
	var listenPort *walk.NumberEdit
	var listenAddr *walk.ComboBox
	var httpsEnable, authEnable, downloadEnable, uploadEnable *walk.CheckBox
	var accessURL, active *walk.PushButton
	var titleName, downloadFolder, uploadFolder *walk.LineEdit

	interfaceList := InterfaceOptions()

	return []Widget{
		Label{
			Text: "Browse URL: ",
		},
		PushButton{
			AssignTo: &accessURL,
			Text:     "http://127.0.0.1:8080/",
			OnClicked: func() {
				OpenBrowserWeb(accessURL.Text())
			},
		},
		Label{
			Text: "Title Name: ",
		},
		LineEdit{
			AssignTo: &titleName,
			Text:     ConfigGet().TitleName,
			OnEditingFinished: func() {

			},
		},
		Label{
			Text: "Download Folder: ",
		},
		LineEdit{
			AssignTo: &downloadFolder,
			OnEditingFinished: func() {

			},
		},
		Label{
			Text: "Upload Folder: ",
		},
		LineEdit{
			AssignTo: &uploadFolder,
			OnEditingFinished: func() {

			},
		},
		Label{
			Text: "Listen Address: ",
		},
		ComboBox{
			AssignTo:           &listenAddr,
			RightToLeftReading: true,
			CurrentIndex:       0,
			Model:              interfaceList,
			OnCurrentIndexChanged: func() {
				// ListenAddressSave(consoleIface.Text())
			},
			OnBoundsChanged: func() {
				// consoleIface.SetCurrentIndex(ConsoleIndex())
			},
		},
		Label{
			Text: "Listen Port: ",
		},
		NumberEdit{
			AssignTo:    &listenPort,
			Value:       float64(ConfigGet().ListenPort),
			ToolTipText: "1~65535",
			MaxValue:    65535,
			MinValue:    1,
			OnValueChanged: func() {
			},
		},
		VSpacer{},
		Composite{
			Layout: HBox{MarginsZero: true},
			Children: []Widget{
				CheckBox{
					AssignTo: &httpsEnable,
					Text:     "Https Enable",
					Checked:  ConfigGet().HttpsEnable,
					OnCheckedChanged: func() {
						err := HttpsEnableSave(httpsEnable.Checked())
						if err != nil {
							ErrorBoxAction(mainWindow, err.Error())
						}
					},
				},
				CheckBox{
					AssignTo: &authEnable,
					Text:     "Auth Enable",
					Checked:  ConfigGet().AuthEnable,
					OnCheckedChanged: func() {
						err := UserEnableSave(authEnable.Checked())
						if err != nil {
							ErrorBoxAction(mainWindow, err.Error())
						}
					},
				},
				CheckBox{
					AssignTo: &downloadEnable,
					Text:     "Download Enable",
					Checked:  ConfigGet().DownloadEnable,
					OnCheckedChanged: func() {
						err := DownloadEnableSave(downloadEnable.Checked())
						if err != nil {
							ErrorBoxAction(mainWindow, err.Error())
						}
					},
				},
				CheckBox{
					AssignTo: &uploadEnable,
					Text:     "Upload Enable",
					Checked:  ConfigGet().UploadEnable,
					OnCheckedChanged: func() {
						err := UploadEnableSave(uploadEnable.Checked())
						if err != nil {
							ErrorBoxAction(mainWindow, err.Error())
						}
					},
				},
			},
		},
		VSpacer{},
		PushButton{
			AssignTo:    &active,
			Image:       ICON_Start,
			Text:        " ",
			ToolTipText: "Startup or Stop",
			MinSize:     Size{Height: 48, Width: 48},
			OnClicked: func() {
				if ServerEnable() {
					active.SetImage(ICON_Stop)
				} else {
					active.SetImage(ICON_Start)
				}
			},
		},
	}
}

func mainWindows() {
	CapSignal(CloseWindows)
	cnt, err := MainWindow{
		Title:          "Simple Http File Server " + VersionGet(),
		Icon:           ICON_Main,
		AssignTo:       &mainWindow,
		MinSize:        Size{Width: mainWindowWidth, Height: mainWindowHeight},
		Size:           Size{Width: mainWindowWidth, Height: mainWindowHeight},
		Layout:         VBox{Margins: Margins{Top: 5, Bottom: 5, Left: 5, Right: 5}},
		Font:           Font{Bold: true},
		MenuItems:      MenuBarInit(),
		StatusBarItems: StatusBarInit(),
		Children: []Widget{
			Composite{
				Layout:   Grid{Columns: 2},
				Children: ConsoleWidget(),
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
	NotifyExit()
}
