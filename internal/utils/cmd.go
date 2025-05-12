package utils

import (
	"os/exec"
	"strings"
)

func Execute(s string) ([]byte, error) {
	args := strings.Split(s, " ")
	cmd := exec.Command(args[0], args[1:]...)
	return cmd.CombinedOutput()
}
