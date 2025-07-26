package utils

import "fmt"

type Error struct {
	Source  string
	Message string
	Line    int
	Column  int
}

func (e Error) Error() string {
	return fmt.Sprintf("error: %s\n%s", e.Message, showPosition(e.Source, e.Line, e.Column))
}
