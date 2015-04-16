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

// Bat is a Go implemented CLI cURL-like tool for humans
// bat [flags] [METHOD] URL [ITEM [ITEM]]
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"

	"github.com/astaxie/bat/httplib"
)

const version = "0.0.1"

var (
	verbose          bool
	form             bool
	pretty           bool
	auth             string
	proxy            string
	isjson           = flag.Bool("json", true, "Send the data as a JSON object")
	method           = flag.String("method", "GET", "HTTP method")
	URL              = flag.String("url", "", "HTTP request URL")
	jsonmap          map[string]interface{}
	contentJsonRegex = `application/json`
)

func init() {
	flag.BoolVar(&pretty, "pretty", true, "Print Json Pretty Fomat")
	flag.BoolVar(&pretty, "p", true, "Print Json Pretty Fomat")
	flag.BoolVar(&verbose, "verbose", false, "Print the whole HTTP exchange (request and response)")
	flag.BoolVar(&verbose, "v", false, "Print the whole HTTP exchange (request and response)")
	flag.BoolVar(&form, "form", false, "Submitting as a form")
	flag.BoolVar(&form, "f", false, "Submitting as a form")
	flag.StringVar(&auth, "auth", "", "HTTP authentication username:password, USER[:PASS]")
	flag.StringVar(&auth, "a", "", "HTTP authentication username:password, USER[:PASS]")
	flag.StringVar(&proxy, "proxy", "", "Proxy host and port, PROXY_URL")
	jsonmap = make(map[string]interface{})
}

func formatResponseBody(res *http.Response, httpreq *httplib.BeegoHttpRequest, pretty bool) string {
	str, err := httpreq.String()
	if err != nil {
		log.Fatalln("can't get the url", err)
	}
	fmt.Println("")
	if pretty && strings.Contains(res.Header.Get("Content-Type"), contentJsonRegex) {
		var output interface{}
		err = json.Unmarshal([]byte(str), &output)
		if err != nil {
			log.Fatal("Response Json Unmarshal", err)
		}
		bstr, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			log.Fatal("Response Json MarshalIndent", err)
		}
		str = string(bstr)
	}

	return str
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
		usage()
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
	if auth != "" {
		userpass := strings.Split(auth, ":")
		if len(userpass) == 2 {
			u.User = url.UserPassword(userpass[0], userpass[1])
		} else {
			u.User = url.User(auth)
		}
	}
	*URL = u.String()
	httpreq := getHTTP(*method, *URL, args)
	if u.User != nil {
		password, _ := u.User.Password()
		httpreq.GetRequest().SetBasicAuth(u.User.Username(), password)
	}
	// Proxy Support
	if proxy != "" {
		purl, err := url.Parse(proxy)
		if err != nil {
			log.Fatal("Proxy Url parse err", err)
		}
		httpreq.SetProxy(http.ProxyURL(purl))
	} else {
		eurl, err := http.ProxyFromEnvironment(httpreq.GetRequest())
		if err != nil {
			log.Fatal("Environment Proxy Url parse err", err)
		}
		httpreq.SetProxy(http.ProxyURL(eurl))
	}

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
			fmt.Println(ColorfulRequest(string(dump)))
			fmt.Println("")
			fmt.Println(Color(res.Proto, Magenta), Color(res.Status, Green))
			for k, v := range res.Header {
				fmt.Println(Color(k, Gray), ":", Color(strings.Join(v, " "), Cyan))
			}
			fmt.Println("")

			body := formatResponseBody(res, httpreq, pretty)
			fmt.Println(ColorfulResponse(body))
		} else {
			body := formatResponseBody(res, httpreq, pretty)
			_, err = os.Stdout.WriteString(body)
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
		fmt.Println("")

		body := formatResponseBody(res, httpreq, pretty)
		fmt.Println(body)
	}
}

var usageinfo string = `bat is a Go implemented CLI cURL-like tool for humans.

Usage:

	bat [flags] [METHOD] URL [ITEM [ITEM]]
	
flags:
  -a, -auth=USER[:PASS]       Pass a username:password pair as the argument
  -f, -form=false             Submitting the data as a form
  -j, -json=true              Send the data in a JSON object
  -v, -verbose=false          Print the whole HTTP exchange (request and response)
  -p, -pretty=true            Print Json Pretty Fomat
  -proxy=PROXY_URL            Proxy with host and port

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
