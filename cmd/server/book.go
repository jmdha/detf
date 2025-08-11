package main

import (
	"os"
	"log"
	"bufio"
)

var book []string

func InitBook(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
	        book = append(book, scanner.Text())
	}
}
