package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"strings"
)

func inSlice(str string, l []string) bool {
	for i := range l {
		if l[i] == str {
			return true
		}
	}
	return false
}

func toRealType(str string) interface{} {
	if i, err := isint(str); err == nil {
		return i
	}
	if b, err := isbool(str); err == nil {
		return b
	}
	if f, err := isfloat(str); err == nil {
		return f
	}
	if strings.HasPrefix(str, "[") && strings.HasSuffix(str, "]") {
		bstr := []byte(str)
		strs := strings.Split(string(bstr[1:len(bstr)-1]), ",")
		var r []interface{}
		for _, s := range strs {
			if i, err := isint(s); err == nil {
				r = append(r, i)
				continue
			}
			if i, err := isbool(s); err == nil {
				r = append(r, i)
				continue
			}
			if i, err := isfloat(s); err == nil {
				r = append(r, i)
				continue
			}
			r = append(r, strings.Trim(s, "\"' "))
		}
		return r
	}
	return str
}

func isint(v string) (i int, err error) {
	return strconv.Atoi(v)
}

func isbool(v string) (bool, error) {
	return strconv.ParseBool(v)
}

func isfloat(v string) (float64, error) {
	return strconv.ParseFloat(v, 64)
}

// Convert bytes to human readable string. Like a 2 MB, 64.2 KB, 52 B
func FormatBytes(i int64) (result string) {
	switch {
	case i > (1024 * 1024 * 1024 * 1024):
		result = fmt.Sprintf("%#.02f TB", float64(i)/1024/1024/1024/1024)
	case i > (1024 * 1024 * 1024):
		result = fmt.Sprintf("%#.02f GB", float64(i)/1024/1024/1024)
	case i > (1024 * 1024):
		result = fmt.Sprintf("%#.02f MB", float64(i)/1024/1024)
	case i > 1024:
		result = fmt.Sprintf("%#.02f KB", float64(i)/1024)
	default:
		result = fmt.Sprintf("%d B", i)
	}
	result = strings.Trim(result, " ")
	return
}

func getSessionDir(batDir string) (sessDir string) {
	usr, homeErr := user.Current()
	if homeErr != nil {
		fmt.Println(homeErr)
	}
	sessDir = usr.HomeDir + `/` + batDir
	if _, err := os.Stat(sessDir); os.IsNotExist(err) {
		os.Mkdir(sessDir, 0755)
	}
	return sessDir
}

func readSessionFile(directory string, file string) http.Cookie {
	jsonb, _ := ioutil.ReadFile(directory + `/` + file)
	var cookieContent http.Cookie
	jsonErr := json.Unmarshal(jsonb, &cookieContent)
	if jsonErr != nil {
		if _, err := os.Stat(directory + `/` + file); !os.IsNotExist(err) {
			fmt.Println(`Unable to parse cookie JSON file`)
			fmt.Println("")
		}
	}
	return cookieContent
}
