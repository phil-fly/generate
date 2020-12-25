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

func TestWechatCounter_Wxid2Qrcode(t *testing.T) {
	Wechat := &WechatCounter{}
	Wechat.Wxid2Qrcode("wxid_u6ucefs3btb222")
}