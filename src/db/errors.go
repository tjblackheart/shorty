package db

import "fmt"

type (
	ErrNotImplemented struct{}
	ErrNotFound       struct{}
	ErrUnique         struct{}
)

func (e ErrNotImplemented) Error() string {
	return fmt.Sprintf("Not implemented.")
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("Entry not found.")
}

func (e ErrUnique) Error() string {
	return fmt.Sprintf("Entry already exists")
}
