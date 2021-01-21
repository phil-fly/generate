package Wifi

import (
	"github.com/phil-fly/generate/utils/cmd"
	"io/ioutil"
)

const GetWifi = "for /f \"skip=9 tokens=1,2 delims=:\" %%i in ('netsh wlan show profiles') do @echo %%j | findstr -i -v echo | netsh wlan show profiles %%j key=clear"

//const GetWifi = `netsh wlan show profiles`
const wifibat = "1.bat"

type WifiCounter struct {
	shellPath     string
	shellFullPath string
}

func (self *WifiCounter) SetShellPath(shellPath string) error {
	self.shellPath = shellPath
	return self.setShellFullPath()
}

func (self *WifiCounter) setShellFullPath() error {
	self.shellFullPath = self.shellPath + "\\" + wifibat
	return self.writeShell()
}

func (self *WifiCounter) writeShell() error {
	return ioutil.WriteFile(self.shellFullPath, []byte(GetWifi), 0744)
}

func (self *WifiCounter) GetWifiInfo() ([]byte, error) {
	info, err := cmd.RunCmdReturnByte(self.shellFullPath)
	if err != nil {
		return nil, err
	}

	//server 进行转换 需要gcc支持
	//var out = make([]byte, len(info))
	//iconv.Convert([]byte(info), out, "gb2312", "utf-8")
//	ioutil.WriteFile("1.txt", info, 0744)
	return info, nil
}
