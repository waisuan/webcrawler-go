package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

func main() {
	input := "https://monzo.com"

	url, err := url.Parse(input)
	if err != nil {
		log.Fatal(err)
	}
	hostname := strings.TrimPrefix(url.Hostname(), "www.")

	fmt.Println(hostname)

	fmt.Println(url.Host)
}
