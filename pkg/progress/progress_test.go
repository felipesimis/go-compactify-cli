package progress

import (
	"bytes"
	"testing"
	"time"

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

func TestCalculateThrottle_ShouldRespectBoundaries(t *testing.T) {
	tests := []struct {
		name        string
		total       int
		concurrency int
		expected    time.Duration
	}{
		{
			name:        "should use min throttle when adjustment is very small",
			total:       10,
			concurrency: 100,
			expected:    40 * time.Millisecond,
		},
		{
			name:        "should use max throttle when adjustment is very large",
			total:       1000,
			concurrency: 1,
			expected:    1000 * time.Millisecond,
		},
		{
			name:        "should calculate intermediate throttle correctly",
			total:       100,
			concurrency: 20,
			expected:    200 * time.Millisecond,
		},
		{
			name:        "should use default throttle when perfectly balanced",
			total:       1,
			concurrency: 1,
			expected:    40 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := calculateThrottle(tt.total, tt.concurrency)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
