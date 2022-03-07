package main

import (
	"fmt"
	"log"
	"regexp"
)

func main() {
	str := "12313113/infra/file-access:v1.1.1.1"
	exp := `^(hub.fastonetech.com)(.*file-access)`
	//match1, err := regexp.Match(`^(hub.fastonetech.com).*`, []byte(str))
	//if err != nil {
	//	return
	//}
	reg, err := regexp.Compile(exp)
	if err != nil {
		log.Fatal(err)
	}
	match := reg.ReplaceAllString(str, "r.fastonetech.com:5000$2")

	fmt.Println(match)
}
