package validation

import (
	"errors"
)

var (
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
	if w.MaxWidth > 0 && w.Width > w.MaxWidth {
		return ErrWidthTooLarge
	}
	return nil
}
