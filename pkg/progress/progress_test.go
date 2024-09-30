// progress/progress_test.go
package progress

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProgressBar(t *testing.T) {
	var buf bytes.Buffer
	total := 10
	pb := NewProgressBar(&buf, total, 20, "Processing")
	assert.NotNil(t, pb.bar)

	pb.Increment()

	initialOutput := buf.String()
	assert.Contains(t, initialOutput, "\x1b[36mProcessing\x1b[0m")
	assert.Contains(t, initialOutput, "10%")
	assert.Contains(t, initialOutput, "1/10")

	for i := 0; i < total-1; i++ {
		pb.Increment()
	}
	pb.Finish()

	output := buf.String()
	assert.Contains(t, output, "\x1b[36mProcessing\x1b[0m")
	assert.Contains(t, output, "\x1b[32m█\x1b[0m")
	assert.Contains(t, output, "\x1b[32m█\x1b[0m")
	assert.Contains(t, output, "[")
	assert.Contains(t, output, "]")
	assert.Contains(t, output, "100%")
	assert.Contains(t, output, "10/10")
}

func TestCalculateThrottle(t *testing.T) {
	throttle := calculateThrottle(100, 120)
	assert.Equal(t, 40*time.Millisecond, throttle)

	throttle = calculateThrottle(100, 1)
	assert.Equal(t, 1000*time.Millisecond, throttle)

	throttle = calculateThrottle(100, 20)
	assert.Equal(t, 200*time.Millisecond, throttle)
}
