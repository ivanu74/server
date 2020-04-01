package main

import (
	"log"
	"os"

	"b.yadro.com/sys/ch-server/infrastructure"
)

//var stdout = log.New(os.Stdout, "[ch-server] ", log.LstdFlags)
var stderr = infrastructure.NewLog(os.Stderr, "[error] [ch-server] ", log.LstdFlags)
var stdout = infrastructure.NewLog(os.Stdout, "[ch-server] ", log.LstdFlags)
