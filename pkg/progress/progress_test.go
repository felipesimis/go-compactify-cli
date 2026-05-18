package progress

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProgressBar_ShouldInitializeWithCorrectDescription(t *testing.T) {
	var buf bytes.Buffer
	description := "Testing Description"
	pb := NewProgressBar(&buf, 10, 1, description)
	pb.Increment()

	assert.NotNil(t, pb.bar)
	assert.Contains(t, buf.String(), description)
}

func TestProgressBar_ShouldUpdateOutputOnIncrementAndFinish(t *testing.T) {
	var buf bytes.Buffer
	pb := NewProgressBar(&buf, 5, 1, "Progress")

	for range 5 {
		pb.Increment()
	}
	pb.Finish()

	output := buf.String()
	assert.Contains(t, output, "5/5")
	assert.Contains(t, output, "100%")
	assert.Contains(t, output, "[")
	assert.Contains(t, output, "]")
}
