package validation

import (
	"errors"
	"strings"
)

type Validation interface {
	Validate() error
}

type ValidationComposite struct {
	Validations []Validation
}

func (v ValidationComposite) Validate() error {
	var errs []string
	for _, validate := range v.Validations {
		if err := validate.Validate(); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}
