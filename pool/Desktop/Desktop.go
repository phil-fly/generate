package Desktop

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"os"
	"path/filepath"
)

//获取桌面路径
func getDesktopPath() (string, error) {
	//net use H: \\${remote_ip}\${remote_folder} "password" /user:"username" /persistent:yes
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Explorer\Shell Folders`, registry.ALL_ACCESS)

	if err != nil {
		return "", err
	}
	defer k.Close()

	s, _, err := k.GetStringValue("Desktop")
	if err != nil {
		return "", err
	}
	return s, nil
}

func GetDesktopFilelist() string {
	desktopPath, err := getDesktopPath()
	if err != nil {
		return ""
	}
	files, _ := filepath.Glob(desktopPath + string(os.PathSeparator) + "*")
	var Info string
	Info = "桌面文件列表:\n"
	for _, v := range files {
		Info += fmt.Sprintf("%s\n", v)
	}
	return Info
}
