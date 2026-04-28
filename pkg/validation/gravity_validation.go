package validation

import (
	"errors"

	"github.com/felipesimis/go-compactify-cli/internal/image"
)

var (
	ErrInvalidGravity = errors.New("invalid gravity")
)

type GravityValidation struct {
	Gravity image.Gravity
}

func (g *GravityValidation) Validate() error {
	if !g.Gravity.IsValid() {
		return ErrInvalidGravity
	}
	return nil
}
