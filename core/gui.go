package core

import (
	"chatgptgo/openai"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"strings"
)

type gui struct {
	infinite                *widget.ProgressBarInfinite
	createImgBtn, submitBtn *widget.Button
	msg                     *strings.Builder
	app                     fyne.App
}

func InitGui() {
	//图片生成条
	myApp := app.New()
	infinite := widget.NewProgressBarInfinite()
	infinite.Hidden = true
	//提交按钮
	subBtn := widget.NewButton("提交问题", nil)
	//生成图片按钮
	createBtn := widget.NewButton("AI图片生成", nil)
	ins := &gui{
		infinite:     infinite,
		createImgBtn: createBtn,
		submitBtn:    subBtn,
		msg:          &strings.Builder{},
		app:          myApp,
	}
	ins.isLockBtn()
	ins.show()
}

func (this_ *gui) show() {
	myWindow := this_.app.NewWindow("GPT Chat")
	//内容展示
	label := widget.NewMultiLineEntry()
	label.Wrapping = fyne.TextWrapBreak //文字自动换行
	//输入input
	input := widget.NewEntry()
	input.SetPlaceHolder("输入问题/图片内容描述")
	//清空按钮
	clearBtn := widget.NewButton("清空", func() {
		label.SetText("")
		label.Refresh()
		openai.ClearAsk()
		this_.msg.Reset()
	})
	this_.submitBtn.OnTapped = func() {
		go func() {
			ask := "我:" + input.Text
			oldContent := this_.msg.String()
			if oldContent == "" {
				//第一次
				this_.msg.WriteString(ask)
			} else {
				this_.msg.WriteString("\n" + ask)
			}
			label.SetText(this_.msg.String())
			this_.infinite.Show()
			this_.msg.WriteString("\n" + openai.Ai)
			openai.AskGptStream(openai.GetAskContent(ask), func(answer string, err error) {
				if err != nil {
					answer = fmt.Sprintf("GPT ERROR %v", err)
					return
				}
				this_.msg.WriteString(answer)
				label.SetText(this_.msg.String())
			})
			this_.infinite.Hide()
			this_.msg.WriteString("\n---------------------------------------------------------")
			label.SetText(this_.msg.String())
		}()
	}
	this_.createImgBtn.OnTapped = func() {
		if input.Text == "" {
			return
		}
		go func() {
			this_.infinite.Show()
			imageUrl, err := openai.CreateImgUrl(input.Text)
			if imageUrl != "" {
				//fmt.Println(imageUrl)
				url, _ := storage.ParseURI(imageUrl)
				w := fyne.CurrentApp().NewWindow("图片")
				w.SetContent(canvas.NewImageFromURI(url))
				w.Resize(fyne.NewSize(500, 500))
				//w.SetFixedSize(true)
				w.Show()
			}
			if err != nil {
				answer := fmt.Sprintf("CreateIMG ERROR %v", err)
				this_.msg.WriteString("\n" + answer)
				this_.msg.WriteString("\n---------------------------------------------------------")
				label.SetText(this_.msg.String())
			}
			this_.infinite.Hide()
		}()
	}
	//布局
	btnBorders := container.NewBorder(nil, nil, clearBtn, this_.createImgBtn, this_.submitBtn)
	content := container.NewVBox(input, btnBorders)
	border := container.NewBorder(content, this_.infinite, nil, nil, label)

	myWindow.SetContent(border)
	myWindow.Resize(fyne.NewSize(600, 600))
	myWindow.SetFixedSize(true)
	myWindow.Show()
	this_.app.Run()
}

func (this_ *gui) isLockBtn() {
	go func() {
		for {
			if this_.infinite.Visible() == false {
				this_.createImgBtn.Enable()
				this_.submitBtn.Enable()
			} else {
				this_.createImgBtn.Disable()
				this_.submitBtn.Disable()
			}
		}
	}()
}
