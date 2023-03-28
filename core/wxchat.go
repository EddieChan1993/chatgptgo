package core

import (
	"chatgptgo/openai"
	goRuntime "chatgptgo/util"
	"context"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"strings"
)

type wxChat struct {
	bot           *openwechat.Bot
	groupNickName string //群名
	groupUserName string //群代号
}

const defGroupNickName = "CDC GANGSTER" //默认监听群昵称

func InitWxChat() {
	//bot := openwechat.DefaultBot()
	wxChatIns := &wxChat{
		bot: openwechat.DefaultBot(openwechat.Desktop), // 桌面模式，上面登录不上的可以尝试切换这种模式
	}
	reloadStorage := openwechat.NewJsonFileHotReloadStorage("storage.json")
	defer reloadStorage.Close()
	// 注册消息处理函数
	wxChatIns.bot.MessageHandler = wxChatIns.listMsg()
	// 注册登陆二维码回调
	wxChatIns.bot.UUIDCallback = openwechat.PrintlnQrcodeUrl
	// 执行热登录
	if err := wxChatIns.bot.HotLogin(reloadStorage, openwechat.NewRetryLoginOption()); err != nil {
		fmt.Println(err)
		return
	}
	// 获取登陆的用户
	self, err := wxChatIns.bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}
	// 获取所有的好友
	//friends, err := self.Friends()
	//fmt.Println(friends, err)
	// 获取所有的群组
	groups, err := self.Groups()
	for _, gname := range groups {
		if gname.NickName == defGroupNickName {
			wxChatIns.groupUserName = gname.UserName
			wxChatIns.groupNickName = gname.NickName
			break
		}
	}
	fmt.Println(groups, err)
	fmt.Println(wxChatIns.groupNickName, wxChatIns.groupUserName)
	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	wxChatIns.bot.Block()
}

func (w *wxChat) listMsg() func(msg *openwechat.Message) {
	return func(msg *openwechat.Message) {
		if msg.IsText() && msg.FromUserName == w.groupUserName {
			content := strings.TrimSpace(msg.Content)
			fmt.Println("msg", content)
			goRuntime.GoRun(func(ctx context.Context) {
				resp, _ := openai.AskGpt(content)
				//_, err := msg.ReplyText(resp)
				fmt.Println("AIBOT:", resp)
				//if err != nil {
				//log.Fatalf("ReplyText Err %v\n", err)
				//}
			})

		}
	}
}
