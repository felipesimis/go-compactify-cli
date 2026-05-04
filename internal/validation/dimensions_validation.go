package validation

import (
	"errors"
)

var (
	ErrInvalidDimensions = errors.New("invalid dimensions")
)

const minDimension = 10

type DimensionsValidation struct {
	Width  int
	Height int
}

func (d *DimensionsValidation) Validate() error {
	if d.Width < minDimension || d.Height < minDimension {
		return ErrInvalidDimensions
	}
	return nil
}
