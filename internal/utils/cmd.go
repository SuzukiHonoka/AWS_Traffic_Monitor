package utils

import (
	"bytes"
	"os/exec"
	"strings"
)

func CheckAwsCli() {
	_, err := exec.LookPath("aws")
	if err != nil {
		panic("please install and login aws-cli first")
	}
}

func Execute(s string) ([]byte, error) {
	args := strings.Split(s, " ")
	cmd := exec.Command(args[0], args[1:]...)

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
