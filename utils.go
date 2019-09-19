package main

import (
	"strings"
)

func removeLastLine(str string) string {
	arr := strings.Split(str, "\n")
	return strings.Join(arr[:len(arr)-2], "\n")
}
