package main

import (
	"fmt"
	"sort"
	"sync"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type UserItem struct {
	Index    int
	UserName string
	Password string

	checked bool
}

type UserTable struct {
	sync.RWMutex

	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder

	items []*UserItem
}

func (n *UserTable) RowCount() int {
	return len(n.items)
}

func (n *UserTable) Value(row, col int) interface{} {
	item := n.items[row]
	switch col {
	case 0:
		return item.Index
	case 1:
		return item.UserName
	case 2:
		return item.Password
	}
	panic("unexpected col")
}

func (n *UserTable) Checked(row int) bool {
	return n.items[row].checked
}

func (n *UserTable) SetChecked(row int, checked bool) error {
	n.items[row].checked = checked
	return nil
}

func (m *UserTable) Sort(col int, order walk.SortOrder) error {
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
			return c(a.UserName < b.UserName)
		case 2:
			return c(a.Password < b.Password)
		}
		panic("unreachable")
	})
	return m.SorterBase.Sort(col, order)
}

var userTable *UserTable
var tableView *walk.TableView

func UserTableInit(userList []UserInfo) {
	item := make([]*UserItem, 0)
	for i, user := range userList {
		item = append(item, &UserItem{Index: i, UserName: user.UserName, Password: user.Password})
	}
	userTable.items = item

	userTable.PublishRowsReset()
	userTable.Sort(userTable.sortColumn, userTable.sortOrder)
}

func UserTableAdd(username, password string) error {
	userTable.Lock()
	defer userTable.Unlock()

	userList := ConfigGet().AuthUsers

	var exist bool
	for i, user := range userList {
		if user.UserName == username {
			userList[i].Password = password
			exist = true
		}
	}

	if !exist {
		userList = append(userList, UserInfo{UserName: username, Password: password})
	}

	UserTableInit(userList)
	return UserListSave(userList)
}

func UserTableDelete() error {
	userTable.Lock()
	defer userTable.Unlock()

	userList := make([]UserInfo, 0)
	for _, items := range userTable.items {
		if !items.checked {
			userList = append(userList, UserInfo{UserName: items.UserName, Password: items.Password})
		}
	}

	if len(userList) == len(userTable.items) {
		return fmt.Errorf("please select some user to delete")
	}

	UserTableInit(userList)
	return UserListSave(userList)
}

func UsersAction() {
	var dlg *walk.Dialog
	var addPB, deletePB, acceptPB *walk.PushButton
	var userLine, passwdLine *walk.LineEdit

	userTable = new(UserTable)
	UserTableInit(ConfigGet().AuthUsers)

	_, err := Dialog{
		AssignTo:      &dlg,
		Title:         "User Authentication Edit",
		Icon:          walk.IconInformation(),
		DefaultButton: &acceptPB,
		CancelButton:  &acceptPB,
		Size:          Size{Width: 450, Height: 300},
		MinSize:       Size{Width: 450, Height: 300},
		Layout:        VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 4, MarginsZero: true},
				Children: []Widget{
					Label{
						Text: "Username: ",
					},
					LineEdit{
						AssignTo: &userLine,
						Text:     "",
					},
					PushButton{
						Text: " Random Generation ",
						OnClicked: func() {
							userLine.SetText(GenerateUsername(10))
						},
					},
					PushButton{
						Text: " Paste Clipboard ",
						OnClicked: func() {
							err := PasteClipboard(userLine.Text())
							if err != nil {
								ErrorBoxAction(dlg, err.Error())
							}
						},
					},
					Label{
						Text: "Password: ",
					},
					LineEdit{
						AssignTo: &passwdLine,
						Text:     "",
					},
					PushButton{
						Text: " Random Generation ",
						OnClicked: func() {
							passwdLine.SetText(GenerateUsername(16))
						},
					},
					PushButton{
						Text: " Paste Clipboard ",
						OnClicked: func() {
							err := PasteClipboard(passwdLine.Text())
							if err != nil {
								ErrorBoxAction(dlg, err.Error())
							}
						},
					},
				},
			},
			Label{
				Text: "UserList: ",
			},
			TableView{
				AssignTo:         &tableView,
				AlternatingRowBG: true,
				ColumnsOrderable: true,
				CheckBoxes:       true,
				Columns: []TableViewColumn{
					{Title: "#", Width: 60},
					{Title: "UserName", Width: 150},
					{Title: "Password", Width: 150},
				},
				StyleCell: func(style *walk.CellStyle) {
					if style.Row()%2 == 0 {
						style.BackgroundColor = walk.RGB(248, 248, 255)
					} else {
						style.BackgroundColor = walk.RGB(220, 220, 220)
					}
				},
				Model: userTable,
				OnCurrentIndexChanged: func() {
					index := tableView.CurrentIndex()
					if 0 <= index && index < len(userTable.items) {
						userLine.SetText(userTable.items[index].UserName)
						passwdLine.SetText(userTable.items[index].Password)
					}
				},
			},
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &addPB,
						Text:     "Add",
						OnClicked: func() {
							username := userLine.Text()
							password := passwdLine.Text()
							if username == "" || password == "" {
								ErrorBoxAction(dlg, "Please input username and password!")
								return
							}
							err := UserTableAdd(username, password)
							if err != nil {
								ErrorBoxAction(dlg, err.Error())
								return
							}
							userLine.SetText("")
							passwdLine.SetText("")
						},
					},
					HSpacer{},
					PushButton{
						AssignTo: &deletePB,
						Text:     "Delete",
						OnClicked: func() {
							err := UserTableDelete()
							if err != nil {
								ErrorBoxAction(dlg, err.Error())
								return
							}
						},
					},
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "Cancel",
						OnClicked: func() {
							dlg.Accept()
						},
					},
					HSpacer{},
				},
			},
		},
	}.Run(mainWindow)

	if err != nil {
		logs.Error(err.Error())
	}
}
