package wechat

import (
	qrcode "github.com/skip2/go-qrcode"
	"golang.org/x/sys/windows/registry"
	"io/ioutil"
	"regexp"
)

type Wechat interface {
}

type WechatCounter struct {
	resourcePath string
}

func (self *WechatCounter) SetResourcePath() error {
	personal, err := getPersonal()
	if err != nil {
		return err
	}
	self.resourcePath = personal + "\\WeChat Files\\All Users\\config\\config.data"
	return nil
}

func (self *WechatCounter) GetResourcePath() string {
	return self.resourcePath
}

func (self *WechatCounter) Wxid2Qrcode(wxid string) error {
	var err = qrcode.WriteFile("weixin://contacts/profile/"+wxid, qrcode.Medium, 256, wxid+".png")
	if err != nil {
		return err
	}
	return nil
}

func (self *WechatCounter) GetWxID() string {
	resourceContent, err := ioutil.ReadFile(self.resourcePath)
	if err != nil {

		return ""
	}

	var re = regexp.MustCompile(`(?m)Documents\\WeChat Files\\(.*)\\config\\AccInfo\.dat`)
	for _, match := range re.FindAllStringSubmatch(string(resourceContent), -1) {
		return match[1]
	}
	return ""
}

// 1、获取我得文档路径
// 2、获取微信id wxid = "C:\\Users\\{{.username}}\\Documents\\WeChat Files\\All Users\\config\\config.data"

//获取我的文档路径
func getPersonal() (string, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Explorer\Shell Folders`, registry.ALL_ACCESS)

	if err != nil {
		return "", err
	}
	defer k.Close()

	s, _, err := k.GetStringValue("Personal")
	if err != nil {
		return "", err
	}
	return s, nil
}
