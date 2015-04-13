package main

import (
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
