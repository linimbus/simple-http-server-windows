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
		Menu{
			Text: "Setting",
			Items: []MenuItem{
				Action{
					Text: "Download Folders",
					OnTriggered: func() {
					},
				},
				Action{
					Text: "Upload Folders",
					OnTriggered: func() {
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
					Text: "Runlog",
					OnTriggered: func() {
						OpenBrowserWeb(RunlogDirGet())
					},
				},
				Separator{},
				Action{
					Text: "Exit",
					OnTriggered: func() {
						CloseWindows()
					},
				},
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
	}
}

func ConsoleWidget() []Widget {
	var listenPort *walk.NumberEdit
	var listenAddr *walk.ComboBox
	var tlsEnable, authEnable *walk.CheckBox
	var accessURL, active *walk.PushButton

	interfaceList := InterfaceOptions()

	return []Widget{
		Label{
			Text: "Access URL: ",
		},
		PushButton{
			AssignTo: &accessURL,
			Text:     "http://127.0.0.1:8080/",
			OnClicked: func() {
				OpenBrowserWeb(accessURL.Text())
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
					AssignTo: &tlsEnable,
					Text:     "TLS Enable",
					Enabled:  true,
					OnCheckedChanged: func() {

					},
				},
				CheckBox{
					AssignTo: &authEnable,
					Text:     "Auth Enable",
					Enabled:  true,
					OnCheckedChanged: func() {

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
