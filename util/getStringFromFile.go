package util

import (
	"bytes"
	"io"
	"log"
	"os"
)

func GetStringFromFile(route string) string{
	file, err := os.Open(route)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	buff := bytes.NewBuffer(nil)
	io.Copy(buff, file)
	return buff.String()
}