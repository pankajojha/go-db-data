package dbutil

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)


func GetEc2(fileName string) string {
	cmd := exec.Command("bash", fileName)
	cmd.Stdin = strings.NewReader("")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return out.String()
}

func GetUrls(fileName string) string {
	//cmd := exec.Command("cat", fileName)
	cmd := exec.Command("cat", "urls")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	return string(out)
}

func GetFileData(fileName string) string {
	fmt.Println("getFileData inside " + fileName)
	cmd := exec.Command("cat", fileName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	} else {
		return string(out)
	}
	return "hello"
}
