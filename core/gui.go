package core

import (
	"chatgptgo/gpt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func InitGui() {
	myApp := app.New()
	myWindow := myApp.NewWindow("GPT Chat")
	//内容展示
	label := widget.NewLabel("")
	label.Wrapping = fyne.TextWrapWord //文字自动换行
	//输入input
	input := widget.NewEntry()
	input.SetPlaceHolder("输入问题")
	//清空按钮
	clearBtn := widget.NewButton("清空", func() {
		label.SetText("")
	})
	//提交按钮
	subBtn := widget.NewButton("提交", func() {
		ask := "我:" + input.Text
		if label.Text == "" {
			//第一次
			label.SetText(ask)
		} else {
			label.SetText(label.Text + "\n----------------------\n" + ask)
		}
		answer := gpt.GptApi.AskGpt(gpt.GetAskContent(ask))
		label.SetText(label.Text + "\n" + answer)
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
