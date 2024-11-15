package main

import (
	"sort"
	"sync"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type FileItem struct {
	Index      int
	InputFile  string
	OutputFile string
	Status     string

	checked bool
}

type FileModel struct {
	sync.RWMutex

	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder

	items []*FileItem
}

func (n *FileModel) RowCount() int {
	return len(n.items)
}

func (n *FileModel) Value(row, col int) interface{} {
	item := n.items[row]
	switch col {
	case 0:
		return item.Index
	case 1:
		return item.InputFile
	case 2:
		return item.OutputFile
	case 3:
		return item.Status
	}
	panic("unexpected col")
}

func (n *FileModel) Checked(row int) bool {
	return n.items[row].checked
}

func (n *FileModel) SetChecked(row int, checked bool) error {
	n.items[row].checked = checked
	return nil
}

func (m *FileModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order
	sort.SliceStable(m.items, func(i, j int) bool {
		a, b := m.items[i], m.items[j]
		c := func(ls bool) bool {
			if m.sortOrder == walk.SortAscending {
				return ls
			}
			return !ls
		}
		switch m.sortColumn {
		case 0:
			return c(a.Index < b.Index)
		case 1:
			return c(a.InputFile < b.InputFile)
		case 2:
			return c(a.OutputFile < b.OutputFile)
		case 3:
			return c(a.Status < b.Status)
		}
		panic("unreachable")
	})
	return m.SorterBase.Sort(col, order)
}

const (
	STATUS_UNDO = ""
	STATUS_DONE = "done"
	STATUS_FAIL = "failed"
)

var consoleFileTable *FileModel
var tableView *walk.TableView
var progressBar *walk.ProgressBar

func init() {
	consoleFileTable = new(FileModel)
	consoleFileTable.items = make([]*FileItem, 0)
}

func FileTableActive() {
	lt := consoleFileTable

	lt.Lock()
	defer lt.Unlock()

	if len(lt.items) == 0 {
		ErrorBoxAction(mainWindow, "No any file, please scan input directory first!")
		return
	}

	for i, item := range lt.items {
		item.Status = STATUS_DONE

		lt.PublishRowsReset()
		lt.Sort(lt.sortColumn, lt.sortOrder)
		progressBar.SetValue(i / len(lt.items))
	}
}

func FileTableInit(items []*FileItem) {
	lt := consoleFileTable

	lt.Lock()
	defer lt.Unlock()

	tableView.SetCurrentIndex(-1)
	lt.items = items
	lt.PublishRowsReset()
	lt.Sort(lt.sortColumn, lt.sortOrder)
}

func TableWidget() []Widget {
	return []Widget{
		Label{
			Text: "File List:",
		},
		TableView{
			AssignTo:         &tableView,
			AlternatingRowBG: true,
			ColumnsOrderable: true,
			CheckBoxes:       false,
			OnItemActivated: func() {
			},
			Columns: []TableViewColumn{
				{Title: "No", Width: 30},
				{Title: "InputFile", Width: 200},
				{Title: "OutputFile", Width: 200},
				{Title: "Status", Width: 60},
			},
			StyleCell: func(style *walk.CellStyle) {
				if style.Row()%2 == 0 {
					style.BackgroundColor = walk.RGB(248, 248, 255)
				} else {
					style.BackgroundColor = walk.RGB(220, 220, 220)
				}
			},
			Model: consoleFileTable,
		},
		ProgressBar{
			AssignTo: &progressBar,
			MaxValue: 100,
			MinValue: 0,
			MaxSize:  Size{Height: 3},
			MinSize:  Size{Height: 3},
			Value:    0,
		},
	}
}
