package main

import (
	"fmt"
	"github.com/stg35/aes/util"
	"github.com/stg35/aes/aes"
)

func main() {
	bytes := []byte(util.GetStringFromFile("./text.txt"))
	var key []byte
	fmt.Scanln(&key)
	s := 0
	for i := range bytes {
		if i % 16 == 0 && i != 0 {
			fmt.Print(aes.Encryption(bytes[s:i], key))
			s = i
		}
	}	
}