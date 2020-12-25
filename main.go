package main

import (
	"fmt"
	"github.com/phil-fly/generate/core/work"
	"os"
)

var (
	RemoteAddr  string
	RemotePort  string
	GenerateMod string
	Rid			string
)

func main() {
	fmt.Println(os.Args[0])
	fmt.Println(RemoteAddr, RemotePort, GenerateMod)
	if GenerateMod == "autotrace" {
		autotrace := &work.Autotrace{}
		autotrace.SetupRemoteAddr(RemoteAddr)
		autotrace.SetupRemotePort(RemotePort)
		go autotrace.Work()
	} else {
		autotrace := &work.Autotrace{}
		autotrace.SetupRemoteAddr(RemoteAddr)
		autotrace.SetupRemotePort("8080")
		autotrace.SetRid(Rid)
		go autotrace.Work()
		work.Generate(RemoteAddr, RemotePort, os.Args[0])
	}
}
