package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"github.com/astaxie/beego/logs"
)

func walkFunc(w *zip.Writer, basePath, path string, stat os.FileInfo) error {
	if stat.IsDir() {
		return nil
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	path, err = filepath.Rel(basePath, path)
	if err != nil {
		return err
	}
	zw, err := w.Create(path)
	if err != nil {
		return err
	}
	if _, err := io.Copy(zw, file); err != nil {
		return err
	}
	return w.Flush()
}

func FileZip(w io.Writer, path string) error {
	basePath := path

	wZip := zip.NewWriter(w)
	defer func() {
		if err := wZip.Close(); err != nil {
			logs.Error(err.Error())
		}
	}()

	return filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return walkFunc(wZip, basePath, path, info)
	})
}
