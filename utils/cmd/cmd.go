package cmd

import (
	"bytes"
	"log"
	"os/exec"
	"syscall"
)

func RunInWindows(cmdstr string) (string, error) {

	cmd := exec.Command("cmd", "/c", cmdstr)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := cmd.Run()
	if err != nil {
		log.Println(err.Error(), stderr.String())
		return "", err
	} else {
		return out.String(), err
	}
}

func RunCmdReturnByte(cmd string) ([]byte, error) {
	cmdExec := exec.Command("cmd", "/C", cmd)
	cmdExec.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	c, err := cmdExec.Output()
	return c, err
}

func RunCmdReturnString(cmd string) (string, error) {
	cmdExec := exec.Command("cmd", "/C", cmd)
	cmdExec.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	c, err := cmdExec.Output()
	return string(c), err
}

func RunCmd(cmd string) error {
	cmdExec := exec.Command("cmd", "/C", cmd)
	cmdExec.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	c, err := cmdExec.Output()
	log.Println(c)
	return err
}
