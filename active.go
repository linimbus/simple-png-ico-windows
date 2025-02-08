package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var scan, cancel *walk.PushButton

func ScanFileList() {
	fileList, err := ReadFileList(ConfigGet().InputDir)
	if err != nil {
		ErrorBoxAction(mainWindow, err.Error())
		return
	}

	fileItem := make([]*FileItem, 0)
	for i, file := range fileList {
		item := new(FileItem)
		item.Index = i
		item.InputFile = file
		item.Status = STATUS_UNDO

		fileItem = append(fileItem, item)
	}
	FileTableInit(fileItem)

	FileTableActive()

	scan.SetEnabled(true)
	cancel.SetEnabled(true)
}

func ActiveWidget() []Widget {
	return []Widget{
		PushButton{
			AssignTo: &scan,
			Text:     "Scan Covert",
			OnClicked: func() {
				scan.SetEnabled(false)
				cancel.SetEnabled(false)

				go ScanFileList()
			},
		},
		PushButton{
			AssignTo: &cancel,
			Text:     "Cancel",
			OnClicked: func() {
				CloseWindows()
			},
		},
		HSpacer{},
	}
}
