package main

import rice "github.com/GeertJohan/go.rice"

var box *rice.Box

func BoxInit() error {
	var err error
	conf := rice.Config{
		LocateOrder: []rice.LocateMethod{rice.LocateEmbedded},
	}
	box, err = conf.FindBox("static")
	if err != nil {
		return err
	}
	return nil
}

func BoxFile() *rice.Box {
	return box
}
