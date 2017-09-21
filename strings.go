package common

import "strings"

// Str2Multiline - split string to string array per line
func Str2Multiline(s string) []string {
	return strings.Split(s, "\n")
}
