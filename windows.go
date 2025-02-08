package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var mainWindow *walk.MainWindow

var mainWindowWidth = 700
var mainWindowHeight = 400

func OptionWidget() []Widget {
	var checkBox *walk.CheckBox

	return []Widget{
		CheckBox{
			AssignTo: &checkBox,
			Text:     "Timestamp",
			Checked:  ConfigGet().TimeStamp,
			OnCheckedChanged: func() {
				TimeStampSave(checkBox.Checked())
			},
		},
	}
}

func MenuBarInit() []MenuItem {
	return []MenuItem{
		Menu{
			Text: "Setting",
			Items: []MenuItem{
				Action{
					Text: "Runlog",
					OnTriggered: func() {
						OpenBrowserWeb(RunlogDirGet())
					},
				},
				Separator{},
				Action{
					Text: "Exit",
					OnTriggered: func() {
						CloseWindows()
					},
				},
			},
		},
		Action{
			Text: "Sponsor",
			OnTriggered: func() {
				AboutAction()
			},
		},
	}
}

func mainWindows() {
	CapSignal(CloseWindows)
	cnt, err := MainWindow{
		Title:     "Simple PNG ICO convert " + VersionGet(),
		Icon:      ICON_Main,
		AssignTo:  &mainWindow,
		MinSize:   Size{Width: mainWindowWidth, Height: mainWindowHeight},
		Size:      Size{Width: mainWindowWidth, Height: mainWindowHeight},
		Layout:    VBox{Margins: Margins{Top: 10, Bottom: 10, Left: 10, Right: 10}},
		MenuItems: MenuBarInit(),
		Children: []Widget{
			Composite{
				Layout:   VBox{MarginsZero: true},
				Children: ConsoleWidget(),
			},
			Composite{
				Layout:   HBox{MarginsZero: true},
				Children: OptionWidget(),
			},
			Composite{
				Layout:   VBox{MarginsZero: true},
				Children: TableWidget(),
			},
			Composite{
				Layout:   HBox{MarginsZero: true},
				Children: ActiveWidget(),
			},
		},
	}.Run()

	if err != nil {
		logs.Error(err.Error())
	} else {
		logs.Info("main windows exit %d", cnt)
	}

	if err := recover(); err != nil {
		logs.Error(err)
	}

	CloseWindows()
}

func CloseWindows() {
	if mainWindow != nil {
		mainWindow.Close()
		mainWindow = nil
	}
}
