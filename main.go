package main

import (
	"fmt"
	"generate/core/work"
	"os"
)

var (
	RemoteAddr string
	RemotePort string
	GenerateMod string
)

func main(){
	fmt.Println(os.Args[0])
	fmt.Println(RemoteAddr,RemotePort,GenerateMod)
	if GenerateMod == "autotrace" {
		autotrace := &work.Autotrace{}
		autotrace.SetupRemoteAddr(RemoteAddr)
		autotrace.SetupRemotePort(RemotePort)
		autotrace.Work()
	}else{
		work.Generate(RemoteAddr,RemotePort,os.Args[0])
	}
}