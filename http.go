package main

import (
	"strings"
	"time"

	"github.com/astaxie/beego/httplib"
)

var defaultSetting = httplib.BeegoHttpSettings{
	ShowDebug:        true,
	UserAgent:        "bat/" + version,
	ConnectTimeout:   60 * time.Second,
	ReadWriteTimeout: 60 * time.Second,
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
	r.Header("Accept", "*/*")
	r.Header("Accept-Encoding", "gzip, deflate")
	for i := range args {
		// Headers
		strs := strings.Split(args[i], ":")
		if len(strs) == 2 {
			r.Header(strs[0], strs[1])
			continue
		}
		// Params
		strs = strings.Split(args[i], "=")
		if len(strs) == 2 {
			r.Param(strs[0], strs[1])
			continue
		}
		// files
		strs = strings.Split(args[i], "@")
		if len(strs) == 2 {
			r.PostFile(strs[0], strs[1])
			continue
		}
	}
	return
}
