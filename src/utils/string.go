package utils

import (
	"fmt"
	"strings"
)

// Wrapper around string pointer which enables reading strings on a file,
// acts like a read-only view of a string.
type String struct {
	Ptr    *string
	Start  int
	Length int
}

func StringFrom(s string, start int, length int) String {
	return String{
		Ptr:    &s,
		Start:  start,
		Length: length,
	}
}

func (a String) Equal(b String) bool {
	return a.Ptr == b.Ptr && a.Start == b.Start && a.Length == b.Length
}

func (s String) String() string {
	return string(*s.Ptr)[s.Start : s.Start+s.Length]
}

func showPosition(s string, line int, column int) (output string) {
	lines := strings.Split(s, "\n")

	if line > 1 {
		output += "...\n"
	}

	if line > 2 {
		output += fmt.Sprintf("%d:%s\n", line-1, lines[line-2])
	}

	output += fmt.Sprintf("%d:%s\n", line, lines[line-1])

	for i := 0; i < column; i++ {
		output += " "
	}

	output += "^^^\n"

	if line < len(lines)-1 {
		output += fmt.Sprintf("%d:%s\n", line+1, lines[line])
	}

	if line < len(lines) {
		output += "...\n"
	}

	return
}
