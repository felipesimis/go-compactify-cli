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
