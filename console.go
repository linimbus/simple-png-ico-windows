package main

import (
	"os"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

func ConsoleWidget() []Widget {
	var inputDir, outputDir *walk.LineEdit
	// var consoleMode *walk.ComboBox
	var checkPNG, checkICON *walk.CheckBox

	return []Widget{
		Composite{
			Layout: HBox{MarginsZero: true},
			Children: []Widget{
				Label{
					Text: "Input Directory: ",
				},
				LineEdit{
					AssignTo: &inputDir,
					Text:     ConfigGet().InputDir,
					OnEditingFinished: func() {
						dir := inputDir.Text()
						if dir != "" {
							stat, err := os.Stat(dir)
							if err != nil {
								ErrorBoxAction(mainWindow, "Input directory is not exist")
								inputDir.SetText(ConfigGet().InputDir)
								return
							}
							if !stat.IsDir() {
								ErrorBoxAction(mainWindow, "Input directory is not directory")
								inputDir.SetText(ConfigGet().InputDir)
								return
							}
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
				Composite{
					Layout: HBox{Margins: Margins{Left: 10, Right: 10}},
					Children: []Widget{
						CheckBox{
							AssignTo: &checkPNG,
							Checked:  true,
							Text:     "PNG",
							OnCheckedChanged: func() {
							},
						},
					},
				},
			},
		},
		Composite{
			Layout: HBox{MarginsZero: true},
			Children: []Widget{
				Label{
					Text: "Output Directory: ",
				},
				LineEdit{
					AssignTo: &outputDir,
					Text:     ConfigGet().OutputDir,
					OnEditingFinished: func() {
						dir := outputDir.Text()
						if dir != "" {
							stat, err := os.Stat(dir)
							if err != nil {
								ErrorBoxAction(mainWindow, "Output directory is not exist")
								outputDir.SetText(ConfigGet().OutputDir)
								return
							}
							if !stat.IsDir() {
								ErrorBoxAction(mainWindow, "Output directory is not directory")
								inputDir.SetText(ConfigGet().OutputDir)
								return
							}
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
				Composite{
					Layout: HBox{Margins: Margins{Left: 10, Right: 10}},
					Children: []Widget{
						CheckBox{
							AssignTo: &checkICON,
							Checked:  true,
							Text:     "ICON",
							OnCheckedChanged: func() {
							},
						},
					},
				},
			},
		},
	}
}
