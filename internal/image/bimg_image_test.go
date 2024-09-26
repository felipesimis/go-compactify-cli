package image

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockedImage() []byte {
	buffer, err := os.ReadFile("../../test/testdata/sample.jpeg")
	if err != nil {
		panic(err)
	}
	return buffer
}

func TestBimgImageWrapper_Size(t *testing.T) {
	img := NewBimgImage(mockedImage())
	size, err := img.Size()

	assert.Nil(t, err)
	assert.Equal(t, 600, size.Width)
	assert.Equal(t, 400, size.Height)
}
