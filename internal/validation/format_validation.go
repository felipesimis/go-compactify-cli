package validation

import "errors"

var (
	ErrFormatRequired = errors.New("format is required")
	ErrInvalidFormat  = errors.New("invalid format")
)

type FormatValidation struct {
	Format string
}

var supportedFormats = map[string]bool{
	"jpg":  true,
	"jpeg": true,
	"png":  true,
	"webp": true,
}

func (f *FormatValidation) Validate() error {
	if f.Format == "" {
		return ErrFormatRequired
	}
	if _, ok := supportedFormats[f.Format]; !ok {
		return ErrInvalidFormat
	}
	return nil
}
