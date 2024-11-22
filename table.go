package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type FileItem struct {
	Index     int
	File      string
	FileTime  time.Time
	MatchFile string
	MatchTime time.Time
	HMAC      string
	Status    string

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
		return item.File
	case 2:
		return item.FileTime.Format(time.DateTime)
	case 3:
		return item.MatchFile
	case 4:
		return item.MatchTime.Format(time.DateTime)
	case 5:
		return item.HMAC
	case 6:
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
			return c(a.File < b.File)
		case 2:
			return c(a.FileTime.Compare(b.FileTime) < 0)
		case 3:
			return c(a.MatchFile < b.MatchFile)
		case 4:
			return c(a.MatchTime.Compare(b.MatchTime) < 0)
		case 5:
			return c(a.HMAC < b.HMAC)
		case 6:
			return c(a.Status < b.Status)
		}
		panic("unreachable")
	})
	return m.SorterBase.Sort(col, order)
}

const (
	STATUS_DONE = "found"
	STATUS_FAIL = "read fail"
)

var fileDupTable *FileModel
var tableView *walk.TableView

func init() {
	fileDupTable = new(FileModel)
	fileDupTable.items = make([]*FileItem, 0)
}

func MoveFileActive(isNew bool) {
	lt := fileDupTable

	outputDir := ConfigGet().DestinationDir

	lt.Lock()
	defer lt.Unlock()

	total := len(lt.items)
	if total == 0 {
		return
	}

	stat, err := os.Stat(outputDir)
	if err != nil {
		if err == os.ErrNotExist {
			if err = os.MkdirAll(outputDir, 0664); err != nil {
				ErrorBoxAction(mainWindow, fmt.Sprintf("Create destination directory fail, %s", err.Error()))
				return
			}
		} else {
			ErrorBoxAction(mainWindow, fmt.Sprintf("Access destination directory fail, %s", err.Error()))
			return
		}
	}

	if !stat.IsDir() {
		ErrorBoxAction(mainWindow, "The destination directory is not directory")
		return
	}

	var moveFile string
	for i := 0; i < total; i++ {
		item := lt.items[0]
		if isNew {
			if item.FileTime.Compare(item.MatchTime) > 0 {
				moveFile = item.File
			} else {
				moveFile = item.MatchFile
			}
		} else {
			if item.FileTime.Compare(item.MatchTime) < 0 {
				moveFile = item.File
			} else {
				moveFile = item.MatchFile
			}
		}
		err = os.Rename(moveFile,
			filepath.Join(outputDir,
				fmt.Sprintf("%s%s", time.Now().Format("2006-01-02T15-04-05.000000"), filepath.Ext(moveFile))))
		if err != nil {
			logs.Error(err.Error())
		}

		lt.items = lt.items[1:]
		lt.PublishRowsReset()
		lt.Sort(lt.sortColumn, lt.sortOrder)

		ProcessUpdate(i * 1000 / total)
	}
}

func DeleteFileActive(isNew bool) {
	lt := fileDupTable

	lt.Lock()
	defer lt.Unlock()

	total := len(lt.items)
	if total == 0 {
		return
	}

	var delFile string

	for i := 0; i < total; i++ {
		item := lt.items[0]
		if isNew {
			if item.FileTime.Compare(item.MatchTime) > 0 {
				delFile = item.File
			} else {
				delFile = item.MatchFile
			}
		} else {
			if item.FileTime.Compare(item.MatchTime) < 0 {
				delFile = item.File
			} else {
				delFile = item.MatchFile
			}
		}

		err := os.Remove(delFile)
		if err != nil {
			logs.Error(err.Error())
		}

		lt.items = lt.items[1:]
		lt.PublishRowsReset()
		lt.Sort(lt.sortColumn, lt.sortOrder)

		ProcessUpdate(i * 1000 / total)
	}
}

func SearchFileActive() {
	lt := fileDupTable

	lt.Lock()
	defer lt.Unlock()

	if ConfigGet().SearchDir == "" {
		ErrorBoxAction(mainWindow, "Please set the correct search directory!")
		return
	}

	stat, err := os.Stat(ConfigGet().SearchDir)
	if err != nil {
		ErrorBoxAction(mainWindow, "The search directory not exist!")
		return
	}

	if !stat.IsDir() {
		ErrorBoxAction(mainWindow, "The search directory is not directory!")
		return
	}

	lt.items = make([]*FileItem, 0)
	lt.PublishRowsReset()
	lt.Sort(lt.sortColumn, lt.sortOrder)

	fileList, err := ReadFileList(ConfigGet().SearchDir)
	if err != nil {
		ErrorBoxAction(mainWindow, fmt.Sprintf("Read the %s search directory fail, %s", ConfigGet().SearchDir, err.Error()))
		return
	}

	fileHmacList := make(map[string]FileInfo, 1024)

	i := 0
	for process, file := range fileList {
		hmac, err := ReadFileHMAC(file.file)
		if err != nil {
			logs.Error(err.Error())
			lt.items = append(lt.items, &FileItem{Index: i, File: file.file, FileTime: file.timestamp, Status: STATUS_FAIL})
			i++
		} else {
			matchFile, b := fileHmacList[hmac]
			if b {
				logs.Info("find the duplicate file, %s <-> %s", file.file, matchFile.file)
				lt.items = append(lt.items, &FileItem{
					Index:     i,
					File:      file.file,
					FileTime:  file.timestamp,
					MatchFile: matchFile.file,
					MatchTime: matchFile.timestamp,
					HMAC:      hmac,
					Status:    STATUS_DONE})
				i++
			} else {
				fileHmacList[hmac] = file
			}
		}

		lt.PublishRowsReset()
		lt.Sort(lt.sortColumn, lt.sortOrder)

		ProcessUpdate(process * 1000 / len(fileList))
	}
}

func TableWidget() []Widget {
	return []Widget{
		Label{
			Text: "Duplicate File List:",
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
				{Title: "File", Width: 250},
				{Title: "Time", Width: 120},
				{Title: "Match File", Width: 250},
				{Title: "Time", Width: 120},
				{Title: "SHA256", Width: 150},
				{Title: "Status", Width: 60},
			},
			StyleCell: func(style *walk.CellStyle) {
				if style.Row()%2 == 0 {
					style.BackgroundColor = walk.RGB(248, 248, 255)
				} else {
					style.BackgroundColor = walk.RGB(220, 220, 220)
				}
			},
			Model: fileDupTable,
		},
	}
}
