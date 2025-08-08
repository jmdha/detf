package main

import (
	"sync"
	"errors"
	"flag"
	"os"
	"os/signal"
	"syscall"
)

type test struct {
	active bool
	id     uint64
	wins   uint64
	losses uint64
	draws  uint64
	book   uint64
}

type match struct {
	id  uint64
	pos string
}

type result struct {
	id   uint64
	win  bool
	draw bool
}

var mtx   sync.Mutex
var tests []test

func HandleResult(res result) {
	mtx.Lock()
	defer mtx.Unlock()
	if res.draw {
		tests[res.id].draws += 1
	} else {
		if res.win {
			tests[res.id].wins += 1
		} else {
			tests[res.id].losses += 1
		}
	}
}

func NextMatch() (match, error) {
	mtx.Lock()
	defer mtx.Unlock()
	for _, test := range tests {
		if !test.active {
			continue
		}
		test.book += 1
		return match {
			id:  test.id,
			pos: book[test.book],
		}, nil
	}
	return match {}, 
	       errors.New("no tests currently active")
}

func main() {
	var path string
	var port int

	// Load command line arguments
	flag.StringVar(&path, "b", "", "path to opening book")
	flag.IntVar(&port, "p", 8080, "port to operate on")
	flag.Parse()

	// Initialise
	InitBook(path)

	// Start processes
	go ServerStart(port)

	// Run untill exit signal
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
