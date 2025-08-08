package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"log"
)

func main() {
	var path string
	var port int
	var git  string

	// Load command line arguments
	flag.StringVar(&path, "b", "",   "path to opening book")
	flag.StringVar(&git,  "g", "",   "which git repo to watch")
	flag.IntVar(&port,    "p", 8080, "port to operate on")
	flag.Parse()

	if path == "" {
		log.Fatalf("Missing path argument")
	}
	if git == "" {
		log.Fatalf("Missing git argument")
	}

	// Initialise
	InitBook(path)

	// Start processes
	go SchedulerStart(git)
	go ServerStart(port)

	// Run untill exit signal
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
