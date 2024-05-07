package check_image

import (
	"mime/multipart"

	"github.com/disintegration/imaging"
)

// IsImage checks if the provided file is an image
func IsImage(file multipart.File) bool {
	_, err := imaging.Decode(file)
	return err == nil
}
