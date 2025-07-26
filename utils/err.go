package utils

import "fmt"

type Error struct {
	Source  String
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("error: %s\n%s", e.Message, e.Source.ShowPosition())
}
