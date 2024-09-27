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

func NewProgressBar(writer io.Writer, total int, description string) *ProgressBar {
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
		progressbar.OptionSetWidth(40),
		progressbar.OptionSetWriter(writer),
		progressbar.OptionThrottle(50*time.Millisecond),
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
