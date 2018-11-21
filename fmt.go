package common

import (
	"fmt"
	"strconv"
	"strings"
)

// add whitespace to the beginning of each line
func Tab(input string) string {
	return NewLineSpace(input, 4)
}

// add whitespace to the beginning of each line
func NewLineSpace(input string, spaces int) string {
	code := `%` + strconv.Itoa(spaces) + `v`
	space := fmt.Sprintf(code, "")
	return space + strings.Replace(input, "\n", "\n"+space, -1)
}

// return wrapped string by the line width
func WordWrap(input string, lineWidth int) string {
	words := strings.Fields(strings.TrimSpace(input))
	if len(words) == 0 {
		return input
	}

	wrapped := words[0]
	remainingRoom := lineWidth - len(wrapped)

	for _, word := range words[1:] {

		if len(word)+1 > remainingRoom {
			wrapped += "\n" + word
			remainingRoom = lineWidth - len(word)
			continue
		}

		wrapped += " " + word
		remainingRoom -= 1 + len(word)
	}
	return wrapped
}

// wrap like this:
//Description: yada yada yada yada yada yada yada yada
//             yada yada yada yada yada yada yada yada
func DescriptionWrap(input string, lineWidth int) string {
	split := strings.SplitN(input, " ", 1)
	if len(split) != 2 {
		return input
	}
	firstWord, remainder := split[0], split[1]
	spaces := len(firstWord) + 1
	wrapWidth := lineWidth - spaces
	wrappedDesc := WordWrap(remainder, wrapWidth)

	space := fmt.Sprintf(`%`+strconv.Itoa(spaces)+`v`, "")
	return firstWord + " " +
		strings.Replace(wrappedDesc, "\n", "\n"+space, -1)
}
