package main

import (
	"fmt"
	"os"
	"strings"

	srv "github.com/linimbus/simple-http-server-windows/server"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

var mainWindow *walk.MainWindow

var mainWindowWidth = 500
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

var listenPort *walk.NumberEdit
var listenAddr *walk.ComboBox
var httpsEnable, authEnable, downloadEnable, uploadEnable *walk.CheckBox
var accessURL, active *walk.PushButton
var titleName, downloadFolder, uploadFolder *walk.LineEdit

func BrowseURLUpdate() {
	addr := ConfigGet().ListenAddr
	if addr == "::" || addr == "0.0.0.0" {
		addr = "localhost"
	}
	schema := "http"
	if ConfigGet().HttpsEnable {
		schema = "https"
	}

	if strings.Contains(addr, ":") {
		accessURL.SetText(fmt.Sprintf("%s://[%s]:%d/", schema, addr, ConfigGet().ListenPort))
	} else {
		accessURL.SetText(fmt.Sprintf("%s://%s:%d/", schema, addr, ConfigGet().ListenPort))
	}
}

func HttpServerStartup() error {
	addr := ConfigGet().ListenAddr

	listen := ""
	if strings.Contains(addr, ":") {
		listen = fmt.Sprintf("[%s]:%d", addr, ConfigGet().ListenPort)
	} else {
		listen = fmt.Sprintf("%s:%d", addr, ConfigGet().ListenPort)
	}

	err := srv.HttpServer(listen, srv.Routes{}, "", "")
	return err
}

func ConsoleWidget() []Widget {
	interfaceList := InterfaceOptions()

	return []Widget{
		Label{
			Text: "Browse URL: ",
		},
		PushButton{
			AssignTo: &accessURL,
			Text:     "",
			OnClicked: func() {
				OpenBrowserWeb(accessURL.Text())
			},
			OnBoundsChanged: func() {
				BrowseURLUpdate()
			},
		},
		Label{
			Text: "Title Name: ",
		},
		LineEdit{
			AssignTo: &titleName,
			Text:     ConfigGet().TitleName,
			OnEditingFinished: func() {
				err := TitleNameSave(titleName.Text())
				if err != nil {
					ErrorBoxAction(mainWindow, err.Error())
				}
			},
		},
		Label{
			Text: "Download Folder: ",
		},
		Composite{
			Layout: HBox{MarginsZero: true},
			Children: []Widget{
				LineEdit{
					AssignTo: &downloadFolder,
					Text:     ConfigGet().DownloadDir,
					OnEditingFinished: func() {
						dir := downloadFolder.Text()
						if dir != "" {
							stat, err := os.Stat(dir)
							if err != nil {
								ErrorBoxAction(mainWindow, "The download folder is not exist")
								downloadFolder.SetText(ConfigGet().DownloadDir)
								return
							}
							if !stat.IsDir() {
								ErrorBoxAction(mainWindow, "The download folder is not directory")
								downloadFolder.SetText(ConfigGet().DownloadDir)
								return
							}
							return
						}
						DownloadDirSave(dir)
					},
				},
				PushButton{
					MaxSize: Size{Width: 30},
					Text:    " ... ",
					OnClicked: func() {
						dlgDir := new(walk.FileDialog)
						dlgDir.FilePath = ConfigGet().DownloadDir
						dlgDir.Flags = win.OFN_EXPLORER
						dlgDir.Title = "Please select a folder as download file directory"

						exist, err := dlgDir.ShowBrowseFolder(mainWindow)
						if err != nil {
							logs.Error(err.Error())
							return
						}
						if exist {
							logs.Info("select %s as download file directory", dlgDir.FilePath)
							downloadFolder.SetText(dlgDir.FilePath)
							DownloadDirSave(dlgDir.FilePath)
						}
					},
				},
			},
		},
		Label{
			Text: "Upload Folder: ",
		},
		Composite{
			Layout: HBox{MarginsZero: true},
			Children: []Widget{
				LineEdit{
					AssignTo: &uploadFolder,
					Text:     ConfigGet().UploadDir,
					OnEditingFinished: func() {
						dir := uploadFolder.Text()
						if dir != "" {
							stat, err := os.Stat(dir)
							if err != nil {
								ErrorBoxAction(mainWindow, "The upload folder is not exist")
								uploadFolder.SetText(ConfigGet().UploadDir)
								return
							}
							if !stat.IsDir() {
								ErrorBoxAction(mainWindow, "The upload folder is not directory")
								uploadFolder.SetText(ConfigGet().UploadDir)
								return
							}
							return
						}
						UploadDirSave(dir)
					},
				},
				PushButton{
					MaxSize: Size{Width: 30},
					Text:    " ... ",
					OnClicked: func() {
						dlgDir := new(walk.FileDialog)
						dlgDir.FilePath = ConfigGet().UploadDir
						dlgDir.Flags = win.OFN_EXPLORER
						dlgDir.Title = "Please select a folder as upload file directory"

						exist, err := dlgDir.ShowBrowseFolder(mainWindow)
						if err != nil {
							logs.Error(err.Error())
							return
						}
						if exist {
							logs.Info("select %s as upload file directory", dlgDir.FilePath)
							uploadFolder.SetText(dlgDir.FilePath)
							UploadDirSave(dlgDir.FilePath)
						}
					},
				},
			},
		},
		Label{
			Text: "Listen Address: ",
		},
		ComboBox{
			AssignTo: &listenAddr,
			CurrentIndex: func() int {
				addr := ConfigGet().ListenAddr
				for i, item := range interfaceList {
					if addr == item {
						return i
					}
				}
				return 0
			},
			Model: interfaceList,
			OnCurrentIndexChanged: func() {
				err := ListenAddressSave(listenAddr.Text())
				if err != nil {
					ErrorBoxAction(mainWindow, err.Error())
				} else {
					BrowseURLUpdate()
				}
			},
			OnBoundsChanged: func() {
				addr := ConfigGet().ListenAddr
				for i, item := range interfaceList {
					if addr == item {
						listenAddr.SetCurrentIndex(i)
						return
					}
				}
				listenAddr.SetCurrentIndex(0)
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
				err := ListenPortSave(int64(listenPort.Value()))
				if err != nil {
					ErrorBoxAction(mainWindow, err.Error())
				} else {
					BrowseURLUpdate()
				}
			},
		},
		VSpacer{},
		Composite{
			Layout: HBox{Margins: Margins{Top: 0, Bottom: 0, Left: 5, Right: 5}},
			Children: []Widget{
				CheckBox{
					AssignTo: &httpsEnable,
					Text:     "Https Enable",
					Checked:  ConfigGet().HttpsEnable,
					OnCheckedChanged: func() {
						err := HttpsEnableSave(httpsEnable.Checked())
						if err != nil {
							ErrorBoxAction(mainWindow, err.Error())
						} else {
							BrowseURLUpdate()
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

				active.SetImage(ICON_Stop)

				go func() {
					err := HttpServerStartup()
					if err != nil {
						ErrorBoxAction(mainWindow, err.Error())
					}
				}()

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
