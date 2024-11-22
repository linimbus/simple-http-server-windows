package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/astaxie/beego/logs"
)

type UserInfo struct {
	UserName string
	Password string
}

type TlsInfo struct {
	CA   string
	Cert string
	Key  string
}

type Config struct {
	TitleName    string
	DownloadDirs []string
	UploadDirs   []string
	AuthEnable   bool
	AuthUsers    []UserInfo
	ListenAddr   string
	ListenPort   int64
	HttpsEnable  bool
	HttpsInfo    TlsInfo
}

var configCache = Config{
	TitleName:    "Simple Http File Server " + VersionGet(),
	DownloadDirs: make([]string, 0),
	UploadDirs:   make([]string, 0),
	AuthEnable:   false,
	AuthUsers:    make([]UserInfo, 0),
	ListenAddr:   "0.0.0.0",
	ListenPort:   9000,
	HttpsEnable:  false,
	HttpsInfo:    TlsInfo{},
}

var configFilePath string
var configLock sync.Mutex

func configSyncToFile() error {
	configLock.Lock()
	defer configLock.Unlock()

	value, err := json.MarshalIndent(configCache, "\t", " ")
	if err != nil {
		logs.Error("json marshal config fail, %s", err.Error())
		return err
	}
	return os.WriteFile(configFilePath, value, 0664)
}

func ConfigGet() *Config {
	return &configCache
}

// func SearchDirSave(path string) error {
// 	configCache.SearchDir = path
// 	return configSyncToFile()
// }

// func DestinationDirDirSave(path string) error {
// 	configCache.DestinationDir = path
// 	return configSyncToFile()
// }

func ConfigInit() error {
	configFilePath = fmt.Sprintf("%s%c%s", ConfigDirGet(), os.PathSeparator, "config.json")

	_, err := os.Stat(configFilePath)
	if err != nil {
		err = configSyncToFile()
		if err != nil {
			logs.Error("config sync to file fail, %s", err.Error())
			return err
		}
	}

	value, err := os.ReadFile(configFilePath)
	if err != nil {
		logs.Error("read config file from app data dir fail, %s", err.Error())
		return err
	}

	err = json.Unmarshal(value, &configCache)
	if err != nil {
		logs.Error("json unmarshal config fail, %s", err.Error())
		return err
	}

	return nil
}
