package worker

import (
	"fmt"
	"log"
	"time"
	"wxrobot/dingtalk"

	"github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

var (
	Bot           *wechaty.Wechaty
	LogOutChan    chan string
	ReturnChan    chan string
	DingTalkRobot *dingtalk.DingTalkRobot
)

func init() {
	LogOutChan = make(chan string, 1)
	ReturnChan = make(chan string, 1)
	DingTalkRobot = dingtalk.NewDingTalkRobot("", "")
}

func Worker() {
	LogOutChan <- "worker start"
	for {
		select {
		case s := <-LogOutChan:
			log.Printf("start bot: %s\n", s)
			go StartBot()
		}
	}
}

func StartBot() {
	Bot = wechaty.NewWechaty()
	Bot.OnScan(onScan)
	Bot.OnLogin(onLogin)
	Bot.OnMessage(onMessage)
	Bot.OnRoomInvite(onRoomInvite)
	Bot.OnFriendship(onFriendShip)
	Bot.OnLogout(onLogout)
	Bot.OnError(onError)
	err := Bot.Start()
	if err != nil {
		DingTalkRobot.SendText(nil, fmt.Sprintf("Bot 启动失败：%s", err), nil)
		time.Sleep(1 * time.Second)
		ReturnChan <- fmt.Sprintf("启动失败: %s", err)
	}
	select {
	case payload := <-ReturnChan:
		LogOutChan <- payload
		return
	}
}

func onScan(ctx *wechaty.Context, qrCode string, status schemas.ScanStatus, data string) {
	DingTalkRobot.SendText(ctx.Context, fmt.Sprintf("当前状态：%v,登陆扫码：https://wechaty.github.io/qrcode/%s", status, qrCode), nil)
}

func onLogin(ctx *wechaty.Context, user *user.ContactSelf) {
	DingTalkRobot.SendText(ctx.Context, fmt.Sprintf("%s 登陆成功", user.Name()), nil)
}

func onMessage(context *wechaty.Context, message *user.Message) {
	if message.Self() {
		return
	}
	// 消息处理
}

func onRoomInvite(context *wechaty.Context, roomInvitation *user.RoomInvitation) {
	// 自动加入群聊
	roomInvitation.Accept()
}

func onFriendShip(context *wechaty.Context, friendship *user.Friendship) {
	if friendship.Type() == schemas.FriendshipTypeReceive {
		// 自动接受好友邀请
		friendship.Accept()
		// 欢迎语
		friendship.Contact().Say("hello")
	}
}

func onLogout(ctx *wechaty.Context, user *user.ContactSelf, reason string) {
	DingTalkRobot.SendText(ctx.Context, fmt.Sprintf("%s 登出：%s", user.Name(), reason), nil)
	ReturnChan <- fmt.Sprintf("%s 登出：%s", user.Name(), reason)
}

func onError(ctx *wechaty.Context, err error) {
	if err == nil || err.Error() == "" {
		return
	}
	DingTalkRobot.SendText(ctx, fmt.Sprintf("Bot Error:%v", err), nil)
}
