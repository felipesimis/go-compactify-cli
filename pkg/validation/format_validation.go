package validation

import "errors"

var (
	ErrInvalidFormat = errors.New("invalid format")
)

type FormatValidation struct {
	Format string
}

func (f *FormatValidation) Validate() error {
	if f.Format == "" {
		return ErrInvalidFormat
	}
	return nil
}
