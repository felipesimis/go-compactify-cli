package validation

import (
	"errors"
)

var (
	ErrInvalidWidth  = errors.New("invalid width")
	ErrWidthTooSmall = errors.New("width is below the minimum allowed value")
	ErrWidthTooLarge = errors.New("width exceeds the maximum allowed value")
)

type WidthValidation struct {
	Width    int
	MinWidth int
	MaxWidth int
}

func (w *WidthValidation) Validate() error {
	if w.Width < w.MinWidth {
		return ErrWidthTooSmall
	}
	return nil
}
