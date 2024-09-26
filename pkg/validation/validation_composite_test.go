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
	assert.Equal(t, ErrMock, err)
}

func TestValidationComposite_MultipleErrors(t *testing.T) {
	validationStub := ValidationComposite{
		Validations: []Validation{
			MockValidator{Err: errors.New("any_error")},
			MockValidator{Err: ErrMock},
		},
	}

	err := validationStub.Validate()
	assert.NotNil(t, err)
	assert.Equal(t, "any_error", err.Error())
	assert.NotEqual(t, ErrMock, err)
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
