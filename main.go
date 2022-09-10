package main

import (
	"log"

	"github.com/thatoddmailbox/aislack/config"
)

func main() {
	log.Println("aislack")

	err := config.Load()
	if err != nil {
		panic(err)
	}
}
