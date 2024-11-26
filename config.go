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
	TitleName string

	DownloadDir    string
	DownloadEnable bool

	UploadDir    string
	UploadEnable bool

	AuthEnable bool
	AuthUsers  []UserInfo

	ListenAddr string
	ListenPort int64

	HttpsEnable bool
	HttpsInfo   TlsInfo
}

var configCache = Config{
	TitleName: "Simple Http File Server " + VersionGet(),

	DownloadDir:    "",
	DownloadEnable: true,

	UploadDir:    "",
	UploadEnable: true,

	AuthEnable: false,
	AuthUsers:  make([]UserInfo, 0),

	ListenAddr: "0.0.0.0",
	ListenPort: 9000,

	HttpsEnable: false,
	HttpsInfo:   TlsInfo{},
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

func UserListSave(userList []UserInfo) error {
	configCache.AuthUsers = userList
	return configSyncToFile()
}

func UserEnableSave(flag bool) error {
	configCache.AuthEnable = flag
	return configSyncToFile()
}

func DownloadDirSave(dir string) error {
	configCache.DownloadDir = dir
	return configSyncToFile()
}

func DownloadEnableSave(flag bool) error {
	configCache.DownloadEnable = flag
	return configSyncToFile()
}

func UploadDirSave(dir string) error {
	configCache.UploadDir = dir
	return configSyncToFile()
}

func UploadEnableSave(flag bool) error {
	configCache.UploadEnable = flag
	return configSyncToFile()
}

func ListenAddressSave(addr string) error {
	configCache.ListenAddr = addr
	return configSyncToFile()
}

func HttpsEnableSave(flag bool) error {
	configCache.HttpsEnable = flag
	return configSyncToFile()
}

func HttpsInfoSave(info TlsInfo) error {
	configCache.HttpsInfo = info
	return configSyncToFile()
}

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
