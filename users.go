package main

import (
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

// func DomainTableUpdate(find string) {
// 	item := make([]*DomainItem, 0)
// 	for idx, v := range domainList {
// 		if strings.Index(v, find) == -1 {
// 			continue
// 		}
// 		item = append(item, &DomainItem{
// 			Index: idx, Domain: v,
// 		})
// 	}
// 	domainTable.items = item
// 	domainTable.PublishRowsReset()
// 	domainTable.Sort(domainTable.sortColumn, domainTable.sortOrder)
// }

// func DomainDelete(owner *walk.Dialog) error {
// 	var deleteList []string
// 	for _, v := range domainTable.items {
// 		if v.checked {
// 			deleteList = append(deleteList, v.Domain)
// 		}
// 	}
// 	if len(deleteList) == 0 {
// 		return fmt.Errorf(LangValue("nochoiceobject"))
// 	}

// 	var remanderList []string
// 	for _, v := range domainList {
// 		var exist bool
// 		for _, v2 := range deleteList {
// 			if v == v2 {
// 				exist = true
// 				break
// 			}
// 		}
// 		if !exist {
// 			remanderList = append(remanderList, v)
// 		}
// 	}

// 	domainList = remanderList
// 	DomainSave(remanderList)

// 	InfoBoxAction(owner, fmt.Sprintf("%v %s", deleteList, LangValue("deletesuccess")))

// 	return nil
// }

func UsersAction() {

	var dlg *walk.Dialog
	var addPB, deletePB, acceptPB *walk.PushButton
	var userLine, passwdLine *walk.LineEdit

	userTable = new(UserTable)
	userTable.items = make([]*UserItem, 0)

	for _, user := range ConfigGet().AuthUsers {
		userTable.items = append(userTable.items, &UserItem{
			UserName: user.UserName,
			Password: user.Password,
		})
	}

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
				Layout: Grid{Columns: 2, MarginsZero: true},
				Children: []Widget{
					Label{
						Text: "UserName: ",
					},
					LineEdit{
						AssignTo: &userLine,
						Text:     "",
					},
					Label{
						Text: "Password: ",
					},
					LineEdit{
						AssignTo: &passwdLine,
						Text:     "",
					},
				},
			},
			TableView{
				ToolTipText:      "UserList:",
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
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						AssignTo: &addPB,
						Text:     "Add",
						OnClicked: func() {
							// 	addDomain := addLine.Text()
							// 	if addDomain == "" {
							// 		ErrorBoxAction(dlg, LangValue("inputdomain"))
							// 		return
							// 	}
							// 	err := DomainAdd(addDomain)
							// 	if err != nil {
							// 		ErrorBoxAction(dlg, err.Error())
							// 		return
							// 	}

							// 	go func() {
							// 		InfoBoxAction(dlg, addDomain+" "+LangValue("addsuccess"))
							// 	}()

							// 	addLine.SetText("")
							// 	findLine.SetText("")
							// 	DomainTableUpdate("")
							// 	RouteUpdate()
						},
					},
					PushButton{
						AssignTo: &deletePB,
						Text:     "Delete",
						// OnClicked: func() {
						// 	err := DomainDelete(dlg)
						// 	if err != nil {
						// 		logs.Error(err.Error())
						// 		ErrorBoxAction(dlg, err.Error())
						// 	} else {
						// 		DomainTableUpdate(findLine.Text())
						// 		RouteUpdate()
						// 	}
						// },
					},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "Cancel",
						OnClicked: func() {
							dlg.Accept()
						},
					},
				},
			},
		},
	}.Run(mainWindow)

	if err != nil {
		logs.Error(err.Error())
	}
}
