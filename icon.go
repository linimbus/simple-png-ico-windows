package main

import (
	"os"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
)

func IconLoadFromBox(filename string, size walk.Size) *walk.Icon {
	body, err := Asset(filename)
	if err != nil {
		logs.Error(err.Error())
		return walk.IconApplication()
	}
	dir := DEFAULT_HOME + "\\icon\\"
	_, err = os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, 0644)
		if err != nil {
			logs.Error(err.Error())
			return walk.IconApplication()
		}
	}
	filepath := dir + filename
	err = SaveToFile(filepath, body)
	if err != nil {
		logs.Error(err.Error())
		return walk.IconApplication()
	}
	icon, err := walk.NewIconFromFileWithSize(filepath, size)
	if err != nil {
		logs.Error(err.Error())
		return walk.IconApplication()
	}
	return icon
}

var ICON_Main *walk.Icon

func IconInit() {
	ICON_Main = IconLoadFromBox("main.ico", walk.Size{
		Width: 512, Height: 512,
	})
}
