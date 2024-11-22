package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/astaxie/beego/logs"
)

func VersionGet() string {
	return "v0.1.0"
}

func SaveToFile(name string, body []byte) error {
	return os.WriteFile(name, body, 0664)
}

func CapSignal(proc func()) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signalChan
		proc()
		logs.Error("recv signcal %s, ready to exit", sig.String())
		os.Exit(-1)
	}()
}

type FileInfo struct {
	file      string
	timestamp time.Time
}

func ReadFileList(dir string) ([]FileInfo, error) {
	output := make([]FileInfo, 0)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			output = append(output, FileInfo{file: absPath, timestamp: info.ModTime()})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return output, nil
}

func ReadFileHMAC(filePath string) (string, error) {
	fd, err := os.Open(filePath)
	if err != nil {
		logs.Error("%s open fail, %s", filePath, err.Error())
		return "", err
	}
	defer fd.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, fd); err != nil {
		logs.Error("%s read fail, %s", filePath, err.Error())
		return "", err
	}

	hashInBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)

	return hashString, nil
}
