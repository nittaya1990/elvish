//go:build !windows && !plan9 && !js
// +build !windows,!plan9,!js

package client

import (
	"os"
	"syscall"
)

func procAttrForSpawn(files []*os.File) *os.ProcAttr {
	return &os.ProcAttr{
		Dir:   "/",
		Env:   []string{},
		Files: files,
		Sys: &syscall.SysProcAttr{
			Setsid: true, // detach from current terminal
		},
	}
}
