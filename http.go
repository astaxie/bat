package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/astaxie/bat/httplib"
)

var defaultSetting = httplib.BeegoHttpSettings{
	ShowDebug:        true,
	UserAgent:        "bat/" + version,
	ConnectTimeout:   60 * time.Second,
	ReadWriteTimeout: 60 * time.Second,
	Gzip:             true,
	DumpBody:         true,
}

func getHTTP(method string, url string, args []string) (r *httplib.BeegoHttpRequest) {
	r = httplib.NewBeegoRequest(url, method)
	r.Setting(defaultSetting)
	r.Header("Accept-Encoding", "gzip, deflate")
	if *isjson {
		r.Header("Accept", "application/json")
	} else if form || method == "GET" {
		r.Header("Accept", "*/*")
	} else {
		r.Header("Accept", "application/json")
	}
	for i := range args {
		// Json raws
		strs := strings.SplitN(args[i], ":=", 2)
		if len(strs) == 2 {
			if strings.HasPrefix(strs[1], "@") {
				f, err := os.Open(strings.TrimLeft(strs[1], "@"))
				if err != nil {
					log.Fatal("Read File", strings.TrimLeft(strs[1], "@"), err)
				}
				content, err := ioutil.ReadAll(f)
				if err != nil {
					log.Fatal("ReadAll from File", strings.TrimLeft(strs[1], "@"), err)
				}
				var j interface{}
				err = json.Unmarshal(content, &j)
				if err != nil {
					log.Fatal("Read from File", strings.TrimLeft(strs[1], "@"), "Unmarshal", err)
				}
				jsonmap[strs[0]] = j
				continue
			}
			jsonmap[strs[0]] = json.RawMessage(strs[1])
			continue
		}
		// Headers
		strs = strings.Split(args[i], ":")
		if len(strs) >= 2 {
			if strs[0] == "Host" {
				r.SetHost(strings.Join(strs[1:], ":"))
			}
			r.Header(strs[0], strings.Join(strs[1:], ":"))
			continue
		}
		// files
		strs = strings.SplitN(args[i], "@", 2)
		if !*isjson && len(strs) == 2 {
			if !form {
				log.Fatal("file upload only support in forms style: -f=true")
			}
			r.PostFile(strs[0], strs[1])
			continue
		}
		// Params
		strs = strings.SplitN(args[i], "=", 2)
		if len(strs) == 2 {
			if strings.HasPrefix(strs[1], "@") {
				f, err := os.Open(strings.TrimLeft(strs[1], "@"))
				if err != nil {
					log.Fatal("Read File", strings.TrimLeft(strs[1], "@"), err)
				}
				content, err := ioutil.ReadAll(f)
				if err != nil {
					log.Fatal("ReadAll from File", strings.TrimLeft(strs[1], "@"), err)
				}
				strs[1] = string(content)
			}
			if form || method == "GET" {
				r.Param(strs[0], strs[1])
			} else {
				jsonmap[strs[0]] = strs[1]
			}
			continue
		}
	}
	if !form && len(jsonmap) > 0 {
		_, err := r.JsonBody(jsonmap)
		if err != nil {
			log.Fatal("fail to marshal json: ", err)
		}
	}
	return
}

func formatResponseBody(res *http.Response, httpreq *httplib.BeegoHttpRequest, pretty bool) string {
	body, err := httpreq.Bytes()
	if err != nil {
		log.Fatalln("can't get the url", err)
	}
	fmt.Println("")
	match, err := regexp.MatchString(contentJsonRegex, res.Header.Get("Content-Type"))
	if err != nil {
		log.Fatalln("failed to compile regex", err)
	}
	if pretty && match {
		var output bytes.Buffer
		err := json.Indent(&output, body, "", "  ")
		if err != nil {
			log.Fatal("Response Json Indent: ", err)
		}

		return output.String()
	}

	return string(body)
}
