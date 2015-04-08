package main

import (
	"log"
	"strings"
	"time"

	"github.com/astaxie/beego/httplib"
)

var defaultSetting = httplib.BeegoHttpSettings{
	ShowDebug:        false,
	UserAgent:        "bat/" + version,
	ConnectTimeout:   60 * time.Second,
	ReadWriteTimeout: 60 * time.Second,
	Gzip:             true,
}

func getHTTP(method string, url string, args []string) (r *httplib.BeegoHttpRequest) {
	switch method {
	case "GET":
		r = httplib.Get(url)
	case "POST":
		r = httplib.Post(url)
	case "PUT":
		r = httplib.Put(url)
	case "HEAD":
		r = httplib.Head(url)
	case "DELETE":
		r = httplib.Delete(url)
	}
	r.Setting(defaultSetting)
	r.Header("Accept-Encoding", "gzip, deflate")
	if form || method == "GET" {
		r.Header("Accept", "*/*")
	} else {
		r.Header("Accept", "application/json")
	}
	for i := range args {
		// Headers
		strs := strings.Split(args[i], ":")
		if len(strs) == 2 {
			if strs[0] == "Host" {
				r.SetHost(strs[1])
			}
			r.Header(strs[0], strs[1])
			continue
		}
		// Params
		strs = strings.Split(args[i], "=")
		if len(strs) == 2 {
			if form {
				r.Param(strs[0], strs[1])
			} else {
				jsonmap[strs[0]] = strs[1]
			}
			continue
		}
		// files
		strs = strings.Split(args[i], "@")
		if len(strs) == 2 {
			if !form {
				log.Fatal("file upload only support in forms style: -f=true")
			}
			r.PostFile(strs[0], strs[1])
			continue
		}
	}
	if !form && len(jsonmap) > 0 {
		r.JsonBody(jsonmap)
	}
	return
}
