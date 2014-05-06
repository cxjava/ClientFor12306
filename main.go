package main

import (
	"log"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {

	var mw *walk.MainWindow
	var acceptPB *walk.PushButton

	if _, err := (MainWindow{
		AssignTo: &mw,
		Title:    "Animal Details",
		MinSize:  Size{100, 50},
		Layout:   VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "用户名:",
					},
					LineEdit{
						Name: "username",
					},

					Label{
						Text: "密　码:",
					},
					LineEdit{
						Name:         "password",
						PasswordMode: true,
					},

					Label{
						Text: "验证码:",
					},
					LineEdit{
						Name: "captcha",
					},
					VSpacer{
						ColumnSpan: 1,
						Size:       8,
					},
					LineEdit{
						// ColumnSpan: 2,
						Name: "captcha",
					},
					VSpacer{
						ColumnSpan: 1,
						Size:       8,
					},
					PushButton{
						// ColumnSpan: 2,
						AssignTo: &acceptPB,
						Text:     "登陆",
						OnClicked: func() {

						},
					},
				},
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}
}
