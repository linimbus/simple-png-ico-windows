package main

import (
	"os"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

// func ConsoleRemoteUpdate() {
// 	remote := ConfigGet().RemoteName
// 	remoteList := ConfigGet().RemoteList

// 	var remoteOptions []string
// 	consoleRemoteProxy.SetCurrentIndex(0)
// 	for i, v := range remoteList {
// 		if v.Name == remote {
// 			consoleRemoteProxy.SetCurrentIndex(i)
// 		}
// 		remoteOptions = append(remoteOptions, v.Name)
// 	}
// 	consoleRemoteProxy.SetModel(remoteOptions)
// }

func ConsoleWidget() []Widget {
	var inputDir, outputDir *walk.LineEdit

	var consoleMode *walk.ComboBox
	var consoleSingle *walk.CheckBox

	return []Widget{
		Label{
			Text: "Input Directory: ",
		},
		Composite{
			Layout: HBox{MarginsZero: true},
			Children: []Widget{
				LineEdit{
					AssignTo: &inputDir,
					Text:     ConfigGet().InputDir,
					OnTextChanged: func() {
						dir := inputDir.Text()
						if dir == "" {
							ErrorBoxAction(mainWindow, "input dir is empty")
							inputDir.SetText("")
							return
						}
						_, err := os.Stat(dir)
						if err != nil {
							ErrorBoxAction(mainWindow, "input dir is not exist")
							inputDir.SetText("")
							return
						}
						InputDirSave(dir)
					},
				},
				PushButton{
					MaxSize: Size{Width: 20},
					Text:    "...",
					OnClicked: func() {
						dlgDir := new(walk.FileDialog)
						dlgDir.FilePath = ConfigGet().InputDir
						dlgDir.Flags = win.OFN_EXPLORER
						dlgDir.Title = "Please select a folder as input directory"

						exist, err := dlgDir.ShowBrowseFolder(mainWindow)
						if err != nil {
							logs.Error(err.Error())
							return
						}
						if exist {
							logs.Info("select %s as input directory", dlgDir.FilePath)
							inputDir.SetText(dlgDir.FilePath)
							InputDirSave(dlgDir.FilePath)
						}
					},
				},
			},
		},
		Label{
			Text: "Output Directory: ",
		},
		Composite{
			Layout: HBox{MarginsZero: true},
			Children: []Widget{
				LineEdit{
					AssignTo: &outputDir,
					Text:     ConfigGet().OutputDir,
					OnTextChanged: func() {
						dir := outputDir.Text()
						if dir == "" {
							ErrorBoxAction(mainWindow, "Output directory is empty")
							outputDir.SetText("")
							return
						}
						_, err := os.Stat(dir)
						if err != nil {
							ErrorBoxAction(mainWindow, "Output directory is not exist")
							outputDir.SetText("")
							return
						}
						OutputDirSave(dir)
					},
				},
				PushButton{
					MaxSize: Size{Width: 20},
					Text:    "...",
					OnClicked: func() {
						dlgDir := new(walk.FileDialog)
						dlgDir.FilePath = ConfigGet().OutputDir
						dlgDir.Flags = win.OFN_EXPLORER
						dlgDir.Title = "Please select a folder as output directory"

						exist, err := dlgDir.ShowBrowseFolder(mainWindow)
						if err != nil {
							logs.Error(err.Error())
							return
						}
						if exist {
							logs.Info("select %s as output directory", dlgDir.FilePath)
							outputDir.SetText(dlgDir.FilePath)
							OutputDirSave(dlgDir.FilePath)
						}
					},
				},
			},
		},
		Label{
			Text: "Covert Mode: ",
		},
		ComboBox{
			AssignTo:     &consoleMode,
			CurrentIndex: 0,
			Model:        []string{"PNG TO ICON"},
			OnCurrentIndexChanged: func() {
			},
		},
		Label{
			Text: "Single File: ",
		},
		CheckBox{
			AssignTo: &consoleSingle,
			Checked:  false,
			OnCheckedChanged: func() {
			},
		},
	}
}
