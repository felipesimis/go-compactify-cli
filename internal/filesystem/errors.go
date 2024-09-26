package filesystem

import "fmt"

type ErrOpenDir struct {
	Path string
	Err  error
}

func (e *ErrOpenDir) Error() string {
	return fmt.Sprintf("failed to open directory %s -> %v", e.Path, e.Err)
}
