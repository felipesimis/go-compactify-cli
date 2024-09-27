package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildOutputPath(t *testing.T) {
	expected := "/path/to/dir/filename"
	result := BuildOutputPath("/path/to/dir", "/path/to/filename")
	assert.Equal(t, expected, result)
}
