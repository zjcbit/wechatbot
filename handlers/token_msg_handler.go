package handlers

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/qingconglaixueit/wechatbot/config"
	"github.com/qingconglaixueit/wechatbot/pkg/logger"
	"github.com/qingconglaixueit/wechatbot/service"
	"strings"
)

var _ MessageHandlerInterface = (*TokenMessageHandler)(nil)

// TokenMessageHandler 口令消息处理器
type TokenMessageHandler struct {
	// 接收到消息
	msg *openwechat.Message
	// 发送的用户
	sender *openwechat.User
	// 实现的用户业务
	service service.UserServiceInterface
}

func isTokenMessage(message *openwechat.Message) bool {
	c := config.LoadConfig()
	if strings.Contains(message.Content, c.SessionClearToken) {
		return true
	}
	if strings.HasPrefix(message.Content, c.HelpToken) {
		return true
	}
	return false
}

func TokenMessageContextHandler() func(ctx *openwechat.MessageContext) {
	return func(ctx *openwechat.MessageContext) {
		msg := ctx.Message
		// 获取口令消息处理器
		handler, err := NewTokenMessageHandler(msg)
		if err != nil {
			logger.Warning(fmt.Sprintf("init token message handler error: %s", err))
		}

		// 获取口令消息处理器
		err = handler.handle()
		if err != nil {
			logger.Warning(fmt.Sprintf("handle token message error: %s", err))
		}

	}
}

// NewTokenMessageHandler 口令消息处理器
func NewTokenMessageHandler(msg *openwechat.Message) (MessageHandlerInterface, error) {
	sender, err := msg.Sender()
	if err != nil {
		return nil, err
	}
	if msg.IsComeFromGroup() {
		sender, err = msg.SenderInGroup()
	}
	userService := service.NewUserService(c, sender)
	handler := &TokenMessageHandler{
		msg:     msg,
		sender:  sender,
		service: userService,
	}

	return handler, nil
}

// handle 处理口令
func (t *TokenMessageHandler) handle() error {
	return t.ReplyText()
}

// ReplyText 回复清空口令
func (t *TokenMessageHandler) ReplyText() error {
	if strings.Contains(t.msg.Content, config.LoadConfig().SessionClearToken) {
		logger.Info("user clear token")
		t.service.ClearUserSessionContext()
		var err error
		if t.msg.IsComeFromGroup() {
			if !t.msg.IsAt() {
				return err
			}
			atText := "@" + t.sender.NickName + "上下文已经清空，请问下一个问题。"
			_, err = t.msg.ReplyText(atText)
		} else {
			_, err = t.msg.ReplyText("上下文已经清空，请问下一个问题。")
		}
		return err
	} else if strings.Contains(t.msg.Content, config.LoadConfig().HelpToken) {
		logger.Info("提供帮助信息")
		var err error
		msg := fmt.Sprintf("你好, 我是智能机器人\n1.与我对话请包含以下字符:%s\n2.进行下一话题请说: %s", config.LoadConfig().ReplyCondition, config.LoadConfig().SessionClearToken)
		if t.msg.IsComeFromGroup() {
			if !t.msg.IsAt() {
				return err
			}
			atText := "@" + t.sender.NickName + msg
			_, err = t.msg.ReplyText(atText)
		} else {
			_, err = t.msg.ReplyText(msg)
		}
		return err
	}
	return nil
}
