package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

var mainWindow *walk.MainWindow

var mainWindowWidth = 500
var mainWindowHeight = 200

func init() {
	go func() {
		for {
			if mainWindow != nil && mainWindow.Visible() {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		ServerAutoStartup()
	}()
}

func MenuBarInit() []MenuItem {
	return []MenuItem{
		Action{
			Text: "Exit",
			OnTriggered: func() {
				CloseWindows()
			},
		},
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
	}
}

var listenPort, listenTimeout *walk.NumberEdit
var listenAddr *walk.ComboBox
var httpsEnable, authEnable, deleteEnable, uploadEnable, zipEnable, autoRun *walk.CheckBox
var serverFolderBut, accessURL, active *walk.PushButton
var serverFolder *walk.LineEdit
var serverInstance *fileHandler
var mutex sync.Mutex

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

func ServerRunning() bool {
	return serverInstance != nil
}

func ServerStart() error {
	var err error
	serverInstance, err = CreateHttpServer(ConfigGet())
	if err != nil {
		return err
	}
	return nil
}

func ServerShutdown() error {
	err := serverInstance.Shutdown()
	if err != nil {
		return err
	}
	serverInstance = nil
	return nil
}

func ServerStatus(flag bool) {
	accessURL.SetEnabled(flag)

	serverFolder.SetEnabled(!flag)
	serverFolderBut.SetEnabled(!flag)
	listenPort.SetEnabled(!flag)
	listenTimeout.SetEnabled(!flag)
	listenAddr.SetEnabled(!flag)
	httpsEnable.SetEnabled(!flag)
	authEnable.SetEnabled(!flag)
	deleteEnable.SetEnabled(!flag)
	uploadEnable.SetEnabled(!flag)
	zipEnable.SetEnabled(!flag)
	autoRun.SetEnabled(!flag)
}

func ServerAutoStartup() {
	if !ConfigGet().AutoStartup {
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	active.SetEnabled(false)
	defer active.SetEnabled(true)

	if ServerRunning() {
		return
	}

	err := ServerStart()
	if err != nil {
		logs.Warning("http server auto startup failed, %s", err.Error())
		return
	}

	if ServerRunning() {
		active.SetImage(ICON_Stop)
	} else {
		active.SetImage(ICON_Start)
	}
	ServerStatus(ServerRunning())
}

func ServerSwitch() {
	mutex.Lock()
	defer mutex.Unlock()

	time.Sleep(time.Millisecond * 500)

	var err error
	if ServerRunning() {
		err = ServerShutdown()
	} else {
		err = ServerStart()
	}

	if err != nil {
		ErrorBoxAction(mainWindow, err.Error())
	}

	if ServerRunning() {
		active.SetImage(ICON_Stop)
	} else {
		active.SetImage(ICON_Start)
	}
	ServerStatus(ServerRunning())

	active.SetEnabled(true)
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
			Enabled:  false,
			OnClicked: func() {
				OpenBrowserWeb(accessURL.Text())
			},
			OnBoundsChanged: func() {
				BrowseURLUpdate()
			},
		},
		Label{
			Text: "Server Folder: ",
		},
		Composite{
			Layout: HBox{MarginsZero: true},
			Children: []Widget{
				LineEdit{
					AssignTo: &serverFolder,
					Text:     ConfigGet().ServerDir,
					OnEditingFinished: func() {
						dir := serverFolder.Text()
						if dir != "" {
							stat, err := os.Stat(dir)
							if err != nil {
								ErrorBoxAction(mainWindow, "The server folder is not exist")
								serverFolder.SetText(ConfigGet().ServerDir)
								return
							}
							if !stat.IsDir() {
								ErrorBoxAction(mainWindow, "The server folder is not directory")
								serverFolder.SetText(ConfigGet().ServerDir)
								return
							}
							return
						}
						ServerDirSave(dir)
					},
				},
				PushButton{
					AssignTo: &serverFolderBut,
					MaxSize:  Size{Width: 30},
					Text:     " ... ",
					OnClicked: func() {
						dlgDir := new(walk.FileDialog)
						dlgDir.FilePath = ConfigGet().ServerDir
						dlgDir.Flags = win.OFN_EXPLORER
						dlgDir.Title = "Please select a folder as server file directory"

						exist, err := dlgDir.ShowBrowseFolder(mainWindow)
						if err != nil {
							logs.Error(err.Error())
							return
						}
						if exist {
							logs.Info("select %s as server file directory", dlgDir.FilePath)
							serverFolder.SetText(dlgDir.FilePath)
							ServerDirSave(dlgDir.FilePath)
						}
					},
				},
			},
		},
		Label{
			Text: "Listen Address: ",
		},
		Composite{
			Layout: HBox{MarginsZero: true},
			Children: []Widget{
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
					Text: "Port: ",
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
				Label{
					Text: "Timeout: ",
				},
				NumberEdit{
					AssignTo:    &listenTimeout,
					Value:       float64(ConfigGet().Timeout),
					ToolTipText: "0~3600 seconds",
					MaxValue:    3600,
					MinValue:    0,
					OnValueChanged: func() {
						err := ListenTimeoutSave(int64(listenTimeout.Value()))
						if err != nil {
							ErrorBoxAction(mainWindow, err.Error())
						}
					},
				},
				Label{
					Text: " Seconds",
				},
			},
		},
		VSpacer{},
		Composite{
			Layout: Grid{Columns: 3, MarginsZero: true},
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
					AssignTo: &deleteEnable,
					Text:     "Delete Enable",
					Checked:  ConfigGet().DeleteEnable,
					OnCheckedChanged: func() {
						err := DeleteEnableSave(deleteEnable.Checked())
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
				CheckBox{
					AssignTo: &zipEnable,
					Text:     "Zip Enable",
					Checked:  ConfigGet().ZipEnable,
					OnCheckedChanged: func() {
						err := ZipEnableSave(zipEnable.Checked())
						if err != nil {
							ErrorBoxAction(mainWindow, err.Error())
						}
					},
				},
				CheckBox{
					AssignTo: &autoRun,
					Text:     "Auto Startup",
					Checked:  ConfigGet().AutoStartup,
					OnCheckedChanged: func() {
						err := AutoStartupSave(autoRun.Checked())
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
			MinSize:     Size{Height: 64, Width: 64},
			OnClicked: func() {
				active.SetEnabled(false)
				go ServerSwitch()
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
	mutex.Lock()
	if ServerRunning() {
		ServerShutdown()
	}
	mutex.Unlock()

	if mainWindow != nil {
		mainWindow.Close()
		mainWindow = nil
	}
	NotifyExit()
}
