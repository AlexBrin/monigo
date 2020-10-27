package monigo

import (
	"reflect"
	"strings"
)

type colors struct {
	Reset string
	Black string
	Red string
	Green string
	Yellow string
	Blue string
	Purple string
	Sky string
	Gray string
	Corral string
	Indigo string
	Pink string
	Cyan string
}

var (
	reset = "\033[0m"

	Text = colors {
		reset,
		"\033[30m",
		"\033[31m",
		"\033[32m",
		"\033[33m",
		"\033[34m",
		"\033[35m",
		"\033[36m",
		"\033[37m",
		"\033[90m",
		"\033[93m",
		"\033[94m",
		"\033[95m",
	}

	Background = colors {
		reset,
		"\033[40m",
		"\033[41m",
		"\033[42m",
		"\033[43m",
		"\033[44m",
		"\033[45m",
		"\033[46m",
		"\033[100m",
		"\033[101m",
		"\033[104m",
		"\033[105m",
		"\033[106m",
	}
)


func ClearColors(row string) string {
	textReflect := reflect.ValueOf(Text)
	for i := 0; i < textReflect.NumField(); i++ {
		fieldReflect := textReflect.Field(i)
		row = strings.ReplaceAll(row, fieldReflect.String(), "")
	}

	bgReflect := reflect.ValueOf(Background)
	for i := 0; i < bgReflect.NumField(); i++ {
		row = strings.ReplaceAll(row, bgReflect.Field(i).String(), "")
	}

	return row
}
