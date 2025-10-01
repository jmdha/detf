package main

import (
	"strings"
)

func SplitStrings(strs []string) ([]string, []string) {
	var left  []string
	var right []string

	for _, str := range strs {
		left  = append(left,  strings.Fields(str)[0])
		right = append(right, strings.Fields(str)[1])
	}

	return left, right
}
