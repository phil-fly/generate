package persistence

import (
	"golang.org/x/sys/windows/registry"
	"testing"
)

func TestPersistence_Enable(t *testing.T) {
	self := &Persistence{
		folderPath: "",
		filename:   "generate.exe",
	}
	Personal,_:= getPersonal()

	self.SetfolderPath(Personal,"\\windows")
	err := self.Enable()
	t.Log(err)
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

func TestPersistence_Disable(t *testing.T) {
	self := &Persistence{
		folderPath: "",
		filename:   "generate.exe",
	}
	Personal,_:= getPersonal()

	self.SetfolderPath(Personal,"\\windows")
	err := self.Disable()
	t.Log(err)
}