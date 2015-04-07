// Copyright 2015 bat authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// Bat is a go implement CLI, cURL-like tool for humans
// bat [flags] [METHOD] URL [ITEM [ITEM]]
package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
)

const version = "0.0.1"

var (
	verbose bool
	form    bool
	auth    string
	json    = flag.Bool("json", true, "Send the data with json object")
	method  = flag.String("method", "GET", "HTTP Method")
	URL     = flag.String("url", "", "HTTP request URL")
)

func init() {
	flag.BoolVar(&verbose, "verbose", false, "Print the whole HTTP exchange (request and response).")
	flag.BoolVar(&verbose, "v", false, "Print the whole HTTP exchange (request and response).")
	flag.BoolVar(&form, "form", false, "Submitting forms")
	flag.BoolVar(&form, "f", false, "Submitting forms")
	flag.StringVar(&auth, "auth", "", "HTTP auth username:password, USER[:PASS]")
	flag.StringVar(&auth, "a", "", "HTTP auth username:password, USER[:PASS]")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		args = filter(args)
	}
	if *URL == "" {
		log.Fatalln("bat should has the URL")
	}
	httpreq := getHTTP(*method, *URL, args)

	res, err := httpreq.Response()
	if err != nil {
		log.Fatalln("can't get the url", err)
	}
	fmt.Println(res.Proto, res.Status)
	for k, v := range res.Header {
		fmt.Println(k, ":", strings.Join(v, " "))
	}
	str, err := httpreq.String()
	if err != nil {
		log.Fatalln("can't get the url", err)
	}
	fmt.Println("")
	fmt.Println(str)
}
