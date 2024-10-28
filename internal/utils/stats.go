package utils

type ImageProcessingStats struct {
	InitialSize     *uint64
	FinalSize       *uint64
	SkippedImages   *uint32
	ProcessedImages *uint32
}

func NewImageProcessingStats(initialSize, finalSize *uint64, skippedImages, processedImages *uint32) *ImageProcessingStats {
	return &ImageProcessingStats{
		InitialSize:     initialSize,
		FinalSize:       finalSize,
		SkippedImages:   skippedImages,
		ProcessedImages: processedImages,
	}
}
