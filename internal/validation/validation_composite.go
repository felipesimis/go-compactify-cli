package validation

import (
	"errors"
)

type Validation interface {
	Validate() error
}

type ValidationComposite struct {
	Validations []Validation
}

func (v ValidationComposite) Validate() error {
	var errs []error
	for _, validate := range v.Validations {
		if err := validate.Validate(); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
