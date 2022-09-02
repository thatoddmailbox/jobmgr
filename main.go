package main

import (
	"log"

	"github.com/thatoddmailbox/jobmgr/config"
	"github.com/thatoddmailbox/jobmgr/data"
)

func main() {
	log.Println("jobmgr")

	err := config.Load()
	if err != nil {
		panic(err)
	}

	err = data.Init()
	if err != nil {
		panic(err)
	}
}
