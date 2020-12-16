package wechat

import "testing"

func TestWechatCounter_SetResourcePath(t *testing.T) {
	Wechat := &WechatCounter{}
	err := Wechat.SetResourcePath()
	t.Log(err)
	t.Log(Wechat.GetResourcePath())
}

func TestWechatCounter_GetWxID(t *testing.T) {
	Wechat := &WechatCounter{}
	err := Wechat.SetResourcePath()
	t.Log(err)
	t.Log(Wechat.GetWxID())
}