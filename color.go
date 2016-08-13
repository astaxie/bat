package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/nwidger/jsoncolor"
)

const (
	Gray = uint8(iota + 90)
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White

	EndColor = "\033[0m"
)

func Color(str string, color uint8) string {
	return fmt.Sprintf("%s%s%s", ColorStart(color), str, EndColor)
}

func ColorStart(color uint8) string {
	return fmt.Sprintf("\033[%dm", color)
}

func ColorfulRequest(str string) string {
	lines := strings.Split(str, "\n")
	if printOption&printReqHeader == printReqHeader {
		strs := strings.Split(lines[0], " ")
		strs[0] = Color(strs[0], Magenta)
		strs[1] = Color(strs[1], Cyan)
		strs[2] = Color(strs[2], Magenta)
		lines[0] = strings.Join(strs, " ")
	}
	for i, line := range lines[1:] {
		substr := strings.Split(line, ":")
		if len(substr) < 2 {
			continue
		}
		substr[0] = Color(substr[0], Gray)
		substr[1] = Color(strings.Join(substr[1:], ":"), Cyan)
		lines[i+1] = strings.Join(substr[:2], ":")
	}
	return strings.Join(lines, "\n")
}

func ColorfulResponse(str, contenttype string, pretty bool) string {
	if strings.Contains(contenttype, contentJsonRegex) {
		str = ColorfulJson(str, pretty)
	} else {
		str = ColorfulHTML(str)
	}
	return str
}

func ColorfulJson(str string, pretty bool) string {
	formatter := jsoncolor.NewFormatter()

	formatter.SpaceColor = color.New()
	formatter.CommaColor = color.New()
	formatter.ColonColor = color.New()
	formatter.ObjectColor = color.New()
	formatter.ArrayColor = color.New()
	formatter.FieldQuoteColor = color.New()
	formatter.FieldColor = color.New(color.FgHiMagenta)
	formatter.StringQuoteColor = color.New()
	formatter.StringColor = color.New(color.FgHiCyan)
	formatter.TrueColor = color.New(color.FgHiCyan)
	formatter.FalseColor = color.New(color.FgHiCyan)
	formatter.NumberColor = color.New(color.FgHiCyan)
	formatter.NullColor = color.New(color.FgHiCyan)
	if !pretty {
		formatter.Prefix = ""
		formatter.Indent = ""
	} else {
		formatter.Prefix = ""
		formatter.Indent = "  "
	}

	buf := bytes.NewBuffer(make([]byte, 0, len(str)))
	err := formatter.Format(buf, []byte(str))
	if err != nil {
		return str
	}

	return buf.String()
}

func ColorfulHTML(str string) string {
	return Color(str, Green)
}
