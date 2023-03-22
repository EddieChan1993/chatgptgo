package core

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strings"
)

var msg strings.Builder

func InitGui() {
	myApp := app.New()
	myWindow := myApp.NewWindow("GPT Chat")
	//内容展示
	label := widget.NewMultiLineEntry()
	label.Wrapping = fyne.TextWrapWord //文字自动换行
	//label.Scroll = container.ScrollVerticalOnly
	//输入input
	input := widget.NewEntry()
	input.SetPlaceHolder("输入问题")
	//清空按钮
	clearBtn := widget.NewButton("清空", func() {
		label.SetText("")
		label.Refresh()
		openai.ClearAsk()
		msg.Reset()
	})
	//提交按钮
	subBtn := widget.NewButton("提交", func() {
		go func() {
			ask := "我:" + input.Text
			oldContent := msg.String()
			if oldContent == "" {
				//第一次
				msg.WriteString(ask)
			} else {
				msg.WriteString("\n" + ask)
			}
			label.SetText(msg.String() + "\n   ....正在思考....")
			answer := openai.AskGpt(openai.GetAskContent(ask))
			msg.WriteString("\n" + answer)
			msg.WriteString("\n---------------------------------------------------------")
			label.SetText(msg.String())
		}()
	})
	//布局
	btnBorders := container.NewBorder(nil, nil, clearBtn, nil, subBtn)
	content := container.NewVBox(input, btnBorders)
	border := container.NewBorder(content, nil, nil, nil, label)

	myWindow.SetContent(border)
	myWindow.Resize(fyne.NewSize(600, 600))
	myWindow.Show()
	myApp.Run()
}
