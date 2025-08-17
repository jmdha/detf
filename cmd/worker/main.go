package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var ip        string
	var processes int

	// Load commad line arguments
	flag.StringVar(&ip,     "i", "localhost:8080", "which ip to connect to")
	flag.IntVar(&processes, "p", 1,                "how many processes")
	flag.Parse()

	// Initialize
	for err := InitClient(ip); err != nil; {
		log.Printf("%v", err)
		time.Sleep(10 * time.Second)
	}

	// Start processes
	for range processes {
		go run()
	}

	// Run untill exit signal
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}

func run() {
	for {
		match, err := RequestMatch()
		if err != nil {
			log.Printf("Failed to retrieve match with error: %v", err)
			time.Sleep(time.Minute)
			continue
		}
		result, err := Sim(match)
		if err != nil {
			log.Printf("Failed to simulate match with error: %v", err)
			time.Sleep(time.Minute)
			continue
		}
		SendResult(result)
	}
}
