package core

import (
	"chatgptgo/openai"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
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
	createBtn := widget.NewButton("生成图片", func() {
		if input.Text == "" {
			return
		}
		go func() {
			imageUrl := openai.CreateImgUrl(input.Text)
			if imageUrl == "" {
				return
			}
			//fmt.Println(imageUrl)
			url, _ := storage.ParseURI(imageUrl)
			w := fyne.CurrentApp().NewWindow("图片")
			w.SetContent(canvas.NewImageFromURI(url))
			w.Resize(fyne.NewSize(500, 500))
			//w.SetFixedSize(true)
			w.Show()
		}()
	})
	//布局
	btnBorders := container.NewBorder(nil, nil, clearBtn, createBtn, subBtn)
	content := container.NewVBox(input, btnBorders)
	border := container.NewBorder(content, nil, nil, nil, label)

	myWindow.SetContent(border)
	myWindow.Resize(fyne.NewSize(600, 600))
	myWindow.SetFixedSize(true)
	myWindow.Show()
	myApp.Run()
}
