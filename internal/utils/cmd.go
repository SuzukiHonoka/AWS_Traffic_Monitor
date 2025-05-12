package utils

import (
	"log"
	"os/exec"
	"strings"
)

func CheckAwsCli() {
	_, err := exec.LookPath("aws")
	if err != nil {
		log.Fatalf("please install and login aws-cli first")
	}
}

func Execute(s string) ([]byte, error) {
	args := strings.Split(s, " ")
	cmd := exec.Command(args[0], args[1:]...)
	return cmd.CombinedOutput()
}
