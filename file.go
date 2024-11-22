package main

import (
	"fmt"
	"os"
)

var DEFAULT_HOME string

func RunlogDirGet() string {
	dir := fmt.Sprintf("%s\\runlog", DEFAULT_HOME)
	_, err := os.Stat(dir)
	if err != nil {
		os.MkdirAll(dir, 644)
	}
	return dir
}

func ConfigDirGet() string {
	dir := fmt.Sprintf("%s\\config", DEFAULT_HOME)
	_, err := os.Stat(dir)
	if err != nil {
		os.MkdirAll(dir, 644)
	}
	return dir
}

func appDataDir() string {
	datadir := os.Getenv("APPDATA")
	if datadir == "" {
		datadir = os.Getenv("CD")
	}
	if datadir == "" {
		datadir = ".\\"
	} else {
		datadir = fmt.Sprintf("%s\\SimpleHttpFileServer", datadir)
	}
	return datadir
}

func appDataDirInit() error {
	dir := appDataDir()
	_, err := os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, 644)
		if err != nil {
			return err
		}
	}
	DEFAULT_HOME = dir
	return nil
}

func FileInit() error {
	err := appDataDirInit()
	if err != nil {
		return err
	}
	return nil
}
