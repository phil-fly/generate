package powershellHistory

import (
	"io/ioutil"
	"os"
)

//C:\Users\whyh\AppData\Roaming\Microsoft\Windows\PowerShell\PSReadLine
//ConsoleHost_history.txt



func GetPowershellHistory() string {
	dataPath := os.Getenv("USERPROFILE") + "\\AppData\\Roaming\\Microsoft\\Windows\\PowerShell\\PSReadLine\\ConsoleHost_history.txt"
	PowershellHistory,err := ioutil.ReadFile(dataPath)
	if err !=nil {
		return ""
	}
	return string(PowershellHistory)
}