package validation

type Validation interface {
	Validate() error
}

type ValidationComposite struct {
	Validations []Validation
}

func (v ValidationComposite) Validate() error {
	for _, validate := range v.Validations {
		if err := validate.Validate(); err != nil {
			return err
		}
	}
	return nil
}
