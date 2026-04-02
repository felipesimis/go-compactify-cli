package utils

import "sync/atomic"

type ImageProcessingStats struct {
	InitialSize     atomic.Uint64
	FinalSize       atomic.Uint64
	SkippedImages   atomic.Uint32
	ProcessedImages atomic.Uint32
}
