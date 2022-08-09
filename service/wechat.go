package service

import (
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	"github.com/silenceper/wechat/v2/miniprogram/auth"
	"github.com/silenceper/wechat/v2/miniprogram/config"
)

const (
	appId     = "wxc96065c6ab7fe91d"
	appSecret = "e4f2cf519418c225aaf25f71c264ac5a"
)

var wechatMiniProgram *miniprogram.MiniProgram = newWechatMiniProgram()

func GetWechatMiniProgram() *miniprogram.MiniProgram {
	return wechatMiniProgram
}

func newWechatMiniProgram() *miniprogram.MiniProgram {
	memory := cache.NewMemory()
	cfg := &config.Config{
		AppID:     appId,
		AppSecret: appSecret,
		Cache:     memory,
	}
	miniprogram := miniprogram.NewMiniProgram(cfg)
	return miniprogram
}

func FormSessionValue(res auth.ResCode2Session) map[string]string {
	sessionValue := make(map[string]string, 2)
	sessionValue["openid"] = res.OpenID
	sessionValue["session_key"] = res.SessionKey
	return sessionValue
}
