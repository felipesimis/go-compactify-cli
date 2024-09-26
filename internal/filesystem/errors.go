package filesystem

import (
	"fmt"
)

type ErrOpenDir struct {
	Err error
}

func (e *ErrOpenDir) Error() string {
	return fmt.Sprintf("failed to open directory -> %v", e.Err)
}

type ErrReadDir struct {
	Path string
	Err  error
}

func (e *ErrReadDir) Error() string {
	return fmt.Sprintf("failed to read directory '%s': %v", e.Path, e.Err)
}

type ErrCreateSiblingDir struct {
	Err error
}

func (e *ErrCreateSiblingDir) Error() string {
	return fmt.Sprintf("failed to create sibling directory -> %v", e.Err)
}
