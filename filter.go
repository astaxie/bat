package main

import (
	"log"
	"strings"
)

var methodList = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

func filter(args []string) []string {
	var i int
	if inSlice(strings.ToUpper(args[i]), methodList) {
		*method = strings.ToUpper(args[i])
		i++
	} else if len(args) > 0 {
		for _, v := range args[1:] {
			// defaults to either GET (with no request data) or POST (with request data).
			// Params
			strs := strings.Split(v, "=")
			if len(strs) == 2 {
				*method = "POST"
				break
			}
			// files
			strs = strings.Split(v, "@")
			if len(strs) == 2 {
				*method = "POST"
				break
			}
		}
	}
	if len(args) <= i {
		log.Fatal("Miss the URL")
	}
	*URL = args[i]
	i++
	return args[i:]
}
