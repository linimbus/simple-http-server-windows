package main

import (
	"os"
	"path/filepath"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
)

func IconLoadFromBox(filename string, size walk.Size) *walk.Icon {
	body, err := BoxFile().Bytes(filename)
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	dir := filepath.Join(DEFAULT_HOME, "icon")
	_, err = os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, 0644)
		if err != nil {
			logs.Error(err.Error())
			return nil
		}
	}
	filepath := filepath.Join(dir, filename)
	err = SaveToFile(filepath, body)
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	logs.Info("load icon file, %s", filepath)

	icon, err := walk.NewIconFromFileWithSize(filepath, size)
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	return icon
}

var ICON_Main *walk.Icon
var ICON_Start *walk.Icon
var ICON_Stop *walk.Icon

var ICON_Max_Size = walk.Size{
	Width: 256, Height: 256,
}

var ICON_Min_Size = walk.Size{
	Width: 64, Height: 64,
}

func IconInit() error {
	ICON_Main = IconLoadFromBox("main.ico", ICON_Max_Size)
	ICON_Start = IconLoadFromBox("start.ico", ICON_Min_Size)
	ICON_Stop = IconLoadFromBox("stop.ico", ICON_Min_Size)
	return nil
}
