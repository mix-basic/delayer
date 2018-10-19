package utils

import (
	"os"
	"os/exec"
)

// 守护执行
func Daemon() {
	args := os.Args
	var nargs []string
	for _, arg := range args {
		d := arg == "-d"
		daemon := arg == "-daemon"
		if d && daemon {
			nargs = append(nargs, arg)
		}
	}
	cmd := exec.Command(args[0], nargs...)
	cmd.Start()
	os.Exit(0)
}
