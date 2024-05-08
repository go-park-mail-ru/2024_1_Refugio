package generate_filename

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateUniqueFileName generates a unique file name based on the current time, random number, and specified format.
func GenerateUniqueFileName(format string) string {
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(1000)

	currentTime := time.Now().Format("20060102_150405")
	uniqueFileName := fmt.Sprintf("%s_%d%s", currentTime, randomNum, format)

	return uniqueFileName
}
