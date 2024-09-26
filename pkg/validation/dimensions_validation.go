package validation

import (
	"errors"
)

var (
	ErrInvalidDimensions = errors.New("invalid dimensions")
)

type DimensionsValidation struct {
	Width  int
	Height int
}

func (d *DimensionsValidation) Validate() error {
	if d.Width < 1 || d.Height < 1 {
		return ErrInvalidDimensions
	}
	return nil
}
