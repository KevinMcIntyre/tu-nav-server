package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func IsEmptyString(str string) bool {
	if len(strings.TrimSpace(str)) < 1 {
		return true
	}
	return false
}

func IsValidEmail(email string) bool {
	Re := regexp.MustCompile(`.+\@.+\..+`)
	return Re.MatchString(email)
}

func WritePid() {
	err := ioutil.WriteFile("SERVER_PID", []byte(strconv.Itoa(os.Getpid())), 0644)
	if err != nil {
		log.Fatalf("Error writing SERVER_PID\n")
	}
}

func ConvertTimeObjToTimeString(timeObj time.Time) string {
	return fmt.Sprintf("%02d:%02d", timeObj.Hour(), timeObj.Minute())
}
