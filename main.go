package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/pquerna/otp/totp"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("First arg requires: %s <full-totp-url>", os.Args[0])
	}

	otpUrl := os.Args[1]

	u, err := url.Parse(otpUrl)
	if err != nil {
		log.Fatalf("Unable to parse 'totp' url: %v", err)
	}

	secret := u.Query().Get("secret")
	if secret == "" {
		log.Fatalf("Missing 'secret' from url")
	}

	otpCode, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		log.Fatalf("Unable to retrieve code: %v", err)
	}

	fmt.Println("Code: ", otpCode)
}
