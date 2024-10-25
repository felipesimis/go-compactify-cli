package validation

import (
	"errors"

	"github.com/h2non/bimg"
)

var (
	ErrInvalidGravity = errors.New("invalid gravity")
)

type GravityValidation struct {
	Gravity bimg.Gravity
}

func (g *GravityValidation) Validate() error {
	if g.Gravity < bimg.GravityCentre || g.Gravity > bimg.GravitySmart {
		return ErrInvalidGravity
	}
	return nil
}
