package main

import (
	"io"
	"os/exec"
)

//执行 shell命令
func RunCMDWithShell(cmdStr string, w io.Writer) {
	cmd := exec.Command("bash", "-c", cmdStr)
	//指定输出位置
	cmd.Stderr = w
	cmd.Stdout = w
	err := cmd.Start()
	if err != nil {
		errs(err)
	}
	err = cmd.Wait()
	if err != nil {
		errs(err)
	}
}

func RunCMDWithDos(c string, w io.Writer, args ...string) {
	var cmdArgs []string
	cmdArgs = append(cmdArgs, "/C")
	cmdArgs = append(cmdArgs, c)
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.Command("cmd", cmdArgs...)
	cmd.Stderr = w
	cmd.Stdout = w
	if err := cmd.Run(); err != nil {
		errs("Error: ", err)
	}

}

func RunCMD(c string, w io.Writer, args ...string) {
	cmd := exec.Command(c, args...)
	cmd.Stderr = w
	cmd.Stdout = w
	if err := cmd.Run(); err != nil {
		errs("Error: ", err)
	}

}
