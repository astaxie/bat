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
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"runtime"
	"strings"
)

const version = "0.0.1"

var (
	verbose bool
	form    bool
	auth    string
	isjson  = flag.Bool("json", true, "Send the data with json object")
	method  = flag.String("method", "GET", "HTTP Method")
	URL     = flag.String("url", "", "HTTP request URL")
	jsonmap map[string]interface{}
)

func init() {
	flag.BoolVar(&verbose, "verbose", false, "Print the whole HTTP exchange (request and response).")
	flag.BoolVar(&verbose, "v", false, "Print the whole HTTP exchange (request and response).")
	flag.BoolVar(&form, "form", false, "Submitting forms")
	flag.BoolVar(&form, "f", false, "Submitting forms")
	flag.StringVar(&auth, "auth", "", "HTTP auth username:password, USER[:PASS]")
	flag.StringVar(&auth, "a", "", "HTTP auth username:password, USER[:PASS]")
	jsonmap = make(map[string]interface{})
}

func main() {
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		args = filter(args)
	}

	var stdin []byte
	if runtime.GOOS != "windows" {
		fi, err := os.Stdin.Stat()
		if err != nil {
			panic(err)
		}

		if fi.Size() != 0 {
			stdin, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Fatal("Read from Stdin", err)
			}
		}
	}

	if *URL == "" {
		log.Fatalln("bat should has the URL")
	}
	if strings.HasPrefix(*URL, ":") {
		urlb := []byte(*URL)
		if *URL == ":" {
			*URL = "http://localhost/"
		} else if len(*URL) > 1 && urlb[1] != '/' {
			*URL = "http://localhost" + *URL
		} else {
			*URL = "http://localhost" + string(urlb[1:])
		}
	}
	if !strings.HasPrefix(*URL, "http://") && !strings.HasPrefix(*URL, "https://") {
		*URL = "http://" + *URL
	}
	u, err := url.Parse(*URL)
	if err != nil {
		log.Fatal(err)
	}
	*URL = u.String()
	httpreq := getHTTP(*method, *URL, args)

	if len(stdin) > 0 {
		var j interface{}
		err = json.Unmarshal(stdin, &j)
		if err != nil {
			log.Fatal("json.Unmarshal", err)
		}
		httpreq.JsonBody(j)
	}

	res, err := httpreq.Response()
	if err != nil {
		log.Fatalln("can't get the url", err)
	}
	if runtime.GOOS != "windows" {
		fi, err := os.Stdout.Stat()
		if err != nil {
			panic(err)
		}
		if fi.Mode()&os.ModeDevice == os.ModeDevice {
			dump := httpreq.DumpRequest()
			fmt.Println(string(dump))
			fmt.Println("")
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
		} else {
			str, err := httpreq.String()
			if err != nil {
				log.Fatalln("can't get the url", err)
			}
			_, err = os.Stdout.WriteString(str)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		dump := httpreq.DumpRequest()
		fmt.Println(string(dump))
		fmt.Println("")
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
}

var usageinfo string = `bat is a Go implemented CLI cURL-like tool for humans.

Usage:

	bat [flags] [METHOD] URL [ITEM [ITEM]]
	
flags:
  -a, -auth USER[:PASS]   Pass a username:password pair as the argument
  -f, -form=false         Submitting the data as forms
  -j, -json=true          Send the data in json object
  -v, -verbose=false      Print the whole HTTP exchange (request and response)

METHOD:
   bat defaults to either GET (if there is no request data) or POST (with request data).

URL:
  The only information needed to perform a request is a URL. The default scheme is http://,
  which can be omitted from the argument; example.org works just fine.

ITEM:
  Can be any of:
    Query string   key=value
    Header         key:value
    Post data      key=value
    File upload    key@/path/file

Example:
    
	bat beego.me
	
more help information please refer to https://github.com/astaxie/bat	
`

func usage() {
	fmt.Println(usageinfo)
	os.Exit(2)
}
