package main

import (
	"crypto/sha256"
	"fmt"
	"io"

	valid "github.com/asaskevich/govalidator"
	"github.com/tjgq/clipboard"
)

type CBType struct {
	Type string
	Data string
}

var (
	old = ""
)

func ClipBoardStart() {

	clipboard.Notify(CBChan)

	for str := range CBChan {
		fmt.Println("GOT DATA!!!", str)
		if len(str) > 0 {
			h256 := sha256.New()
			io.WriteString(h256, str)
			sha := fmt.Sprintf("%x", h256.Sum(nil))

			cb := CBType{
				Data: str,
				Type: "",
			}

			if sha != old {
				old = sha

				if valid.IsEmail(str) {
					cb.Type = "email"
				} else if valid.IsURL(str) {
					cb.Type = "url"
				} else if valid.IsJSON(str) {
					cb.Type = "json"
				} else if valid.IsIP(str) {
					if valid.IsIPv4(str) {
						cb.Type = "ipv4"
					} else if valid.IsIPv6(str) {
						cb.Type = "ipv6"
					}
				} else {
					cb.Type = "none"
				}

				OutputChan <- cb
			}
		}
	}

}
