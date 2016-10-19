package utils

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func IsEmptyString(str string) bool {
	if len(strings.TrimSpace(str)) < 1 {
		return true
	}
	return false
}

func WritePid() {
	err := ioutil.WriteFile("SERVER_PID", []byte(strconv.Itoa(os.Getpid())), 0644)
	if err != nil {
		log.Fatalf("Error writing SERVER_PID\n")
	}
}
