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

type ReadDirError struct {
	Path string
	Err  error
}

func (e *ReadDirError) Error() string {
	return fmt.Sprintf("failed to read directory '%s': %v", e.Path, e.Err)
}
