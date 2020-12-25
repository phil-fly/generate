package Wifi

import (
	"golang.org/x/sys/windows/registry"
	"testing"
)

func TestWifiCounter_GetWifiInfo(t *testing.T) {
	wifi := &WifiCounter{}
	Personal, _ := getPersonal()
	t.Log(Personal)
	wifi.SetShellPath(Personal)
	wifi.GetWifiInfo()
}

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
