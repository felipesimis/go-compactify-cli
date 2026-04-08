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

func TestValidationComposite_Error(t *testing.T) {
	validationStub := ValidationComposite{
		Validations: []Validation{
			MockValidator{Err: ErrMock},
		},
	}

	err := validationStub.Validate()
	assert.NotNil(t, err)
	assert.Equal(t, ErrMock.Error(), err.Error())
}

func TestValidationComposite_MultipleErrors(t *testing.T) {
	validationStub := ValidationComposite{
		Validations: []Validation{
			MockValidator{Err: errors.New("error_one")},
			MockValidator{Err: errors.New("error_two")},
		},
	}

	err := validationStub.Validate()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "error_one")
	assert.Contains(t, err.Error(), "error_two")
	assert.Contains(t, err.Error(), "\n")
}

func TestValidationComposite_Success(t *testing.T) {
	validationStub := ValidationComposite{
		Validations: []Validation{
			MockValidator{Err: nil},
		},
	}

	err := validationStub.Validate()
	assert.Nil(t, err)
}
