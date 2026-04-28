package validation

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockValidator struct {
	Err error
}

func (m MockValidator) Validate() error {
	return m.Err
}

var ErrMock = errors.New("mocked error")

func TestValidationComposite_ShouldReturnError_WhenSingleValidatorFails(t *testing.T) {
	validationStub := ValidationComposite{
		Validations: []Validation{
			MockValidator{Err: ErrMock},
		},
	}

	err := validationStub.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), ErrMock.Error())
}

func TestValidationComposite_ShouldReturnCombinedErrors_WhenMultipleValidatorsFail(t *testing.T) {
	errOne := errors.New("error_one")
	errTwo := errors.New("error_two")
	validationStub := ValidationComposite{
		Validations: []Validation{
			MockValidator{Err: errOne},
			MockValidator{Err: errTwo},
		},
	}

	err := validationStub.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), errOne.Error())
	assert.Contains(t, err.Error(), errTwo.Error())
	assert.Contains(t, err.Error(), "\n")
}

func TestValidationComposite_ShouldSucceed_WhenAllValidatorsPass(t *testing.T) {
	validationStub := ValidationComposite{
		Validations: []Validation{
			MockValidator{Err: nil},
		},
	}

	err := validationStub.Validate()
	assert.NoError(t, err)
}
