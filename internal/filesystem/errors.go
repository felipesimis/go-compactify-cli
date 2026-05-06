package filesystem

import "fmt"

type ErrOpenDir struct {
	Path string
	Err  error
}

func (e *ErrOpenDir) Error() string {
	return fmt.Sprintf("failed to open directory '%s': %v", e.Path, e.Err)
}
func (e *ErrOpenDir) Unwrap() error { return e.Err }

type ErrReadDir struct {
	Path string
	Err  error
}

func (e *ErrReadDir) Error() string {
	return fmt.Sprintf("failed to read directory '%s': %v", e.Path, e.Err)
}
func (e *ErrReadDir) Unwrap() error { return e.Err }

type ErrCreateDir struct {
	Path string
	Err  error
}

func (e *ErrCreateDir) Error() string {
	return fmt.Sprintf("failed to create directory '%s': %v", e.Path, e.Err)
}
func (e *ErrCreateDir) Unwrap() error { return e.Err }

type ErrCreateSiblingDir struct {
	Path string
	Err  error
}

func (e *ErrCreateSiblingDir) Error() string {
	return fmt.Sprintf("failed to create sibling directory for '%s': %v", e.Path, e.Err)
}
func (e *ErrCreateSiblingDir) Unwrap() error { return e.Err }

type ErrReadFile struct {
	Path string
	Err  error
}

func (e *ErrReadFile) Error() string {
	return fmt.Sprintf("failed to read file '%s': %v", e.Path, e.Err)
}
func (e *ErrReadFile) Unwrap() error { return e.Err }

type ErrWriteFile struct {
	Path string
	Err  error
}

func (e *ErrWriteFile) Error() string {
	return fmt.Sprintf("failed to write file '%s': %v", e.Path, e.Err)
}
func (e *ErrWriteFile) Unwrap() error { return e.Err }

type ErrWalk struct {
	Path string
	Err  error
}

func (e *ErrWalk) Error() string {
	return fmt.Sprintf("failed to walk directory '%s': %v", e.Path, e.Err)
}
func (e *ErrWalk) Unwrap() error { return e.Err }

type ErrRelPath struct {
	Root   string
	Target string
	Err    error
}

func (e *ErrRelPath) Error() string {
	return fmt.Sprintf("failed to calculate relative path for '%s' against root '%s': %v", e.Target, e.Root, e.Err)
}
func (e *ErrRelPath) Unwrap() error { return e.Err }

type ErrFileInfo struct {
	Path string
	Err  error
}

func (e *ErrFileInfo) Error() string {
	return fmt.Sprintf("failed to get file info for '%s': %v", e.Path, e.Err)
}
func (e *ErrFileInfo) Unwrap() error { return e.Err }
