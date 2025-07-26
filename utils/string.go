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

func Encompass(strings ...String) String {
	minStart := 0
	maxEnd := 0
	// first make sure they all have the same pointer
	ptr := strings[0].Ptr
	for _, s := range strings {
		if s.Ptr != ptr {
			panic("all strings must have the same pointer")
		}
	}
	for _, s := range strings {
		if s.Start < minStart {
			minStart = s.Start
		}
		if s.Start+s.Length > maxEnd {
			maxEnd = s.Start + s.Length
		}
	}
	return StringFrom(*ptr, minStart, maxEnd-minStart)
}

func (a String) Equal(b String) bool {
	return a.Ptr == b.Ptr && a.Start == b.Start && a.Length == b.Length
}

func (s String) String() string {
	return string(*s.Ptr)[s.Start : s.Start+s.Length]
}

func (s String) ShowPosition() (output string) {
	lines := strings.Split(*s.Ptr, "\n")

	line, column := s.GetLineAndColumn()

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

func (s String) GetLineAndColumn() (int, int) {
	line := 1
	column := 1
	for i := 0; i < s.Start; i++ {
		if (*s.Ptr)[i] == '\n' {
			line++
			column = 1
		} else {
			column++
		}
	}
	return line, column
}
