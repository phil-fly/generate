package wincreds

import (
	"encoding/json"
	"fmt"
	"github.com/danieljoos/wincred"
)

type Credential struct {
	TargetName     string 	`json:"TargetName"`
	LastWritten    string	`json:"LastWritten"`
	UserName       string	`json:"UserName"`
}

func GetWincred() ([]byte,error){

	var Credentials []Credential
	creds, err := wincred.List()
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	for _,v := range creds {
		node:= Credential{
			TargetName:v.TargetName,
			LastWritten: v.LastWritten.Format("2006-01-02 15:04:05"),
			UserName: v.UserName,
		}
		Credentials = append(Credentials, node)
	}
	b, err := json.Marshal(Credentials)
	if err != nil {
		fmt.Println("JSON ERR:", err)
	}
	return b,err
}