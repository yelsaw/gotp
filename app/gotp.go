package main

import (
	"fmt"
	"log"
	"os"

	gotp "github.com/yelsaw/gotp"
)

func main() {
	if len(os.Args) != 2 {
		log.Printf("First arg requires string or path")
		os.Exit(0)
	}

	url := gotp.ArgParser(os.Args[1])

	message, err := gotp.UrlParser(url)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := gotp.Interactive(message).Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(0)
	}
}
