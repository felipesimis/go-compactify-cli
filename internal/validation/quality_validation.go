package validation

import "errors"

var (
	ErrInvalidQuality = errors.New("quality must be between 1 and 100")
)

type QualityValidation struct {
	Quality int
}

func (q *QualityValidation) Validate() error {
	if q.Quality < 1 || q.Quality > 100 {
		return ErrInvalidQuality
	}
	return nil
}
