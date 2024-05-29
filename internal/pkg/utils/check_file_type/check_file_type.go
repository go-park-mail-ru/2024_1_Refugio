package check_file_type

import "strings"

// GetFileType returns the file type based on its Content-Type.
func GetFileType(contentType string) string {
	if strings.HasPrefix(contentType, "image") {
		return "Image"
	} else if strings.HasPrefix(contentType, "video") {
		return "Video"
	} else if strings.HasPrefix(contentType, "audio") {
		return "Audio"
	} else if strings.HasPrefix(contentType, "application/pdf") {
		return "PDF Document"
	} else if strings.HasPrefix(contentType, "text/plain") {
		return "Text Document"
	} else {
		return "Unknown"
	}
}
