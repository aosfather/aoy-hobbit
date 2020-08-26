package main

import (
	"os"
	"testing"
)

func TestRunCMD(t *testing.T) {
	RunCMDWithDos("dir", os.Stdout, "c:")
	RunCMD("D:\\tools\\Notepad++\\notepad++.exe", os.Stdout, "help.txt")
}
