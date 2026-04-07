package progress

import (
	"fmt"
	"io"
	"time"

	"github.com/schollz/progressbar/v3"
)

type ProgressBar struct {
	bar *progressbar.ProgressBar
}

const (
	defaultBarWidth = 40
	defaultThrottle = 40 * time.Millisecond
	minThrottle     = 40 * time.Millisecond
	maxThrottle     = 1000 * time.Millisecond
)

func NewProgressBar(writer io.Writer, total, concurrency int, description string) *ProgressBar {
	throttle := calculateThrottle(total, concurrency)

	bar := progressbar.NewOptions(total,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSetDescription(fmt.Sprintf("[cyan]%s[reset]", description)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]█[reset]",
			SaucerHead:    "[green]█[reset]",
			SaucerPadding: "░",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionSetWidth(defaultBarWidth),
		progressbar.OptionSetWriter(writer),
		progressbar.OptionThrottle(throttle),
		progressbar.OptionUseANSICodes(true),
		progressbar.OptionOnCompletion(func() {
			fmt.Println()
		}),
	)
	return &ProgressBar{bar: bar}
}

func (p *ProgressBar) Increment() {
	p.bar.Add(1)
}

func (p *ProgressBar) Finish() {
	p.bar.Finish()
}

func calculateThrottle(total, concurrency int) time.Duration {
	adjustmentFactor := float64(total) / float64(concurrency)
	throttle := defaultThrottle * time.Duration(adjustmentFactor)

	if throttle < minThrottle {
		return minThrottle
	}
	if throttle > maxThrottle {
		return maxThrottle
	}

	return throttle
}
