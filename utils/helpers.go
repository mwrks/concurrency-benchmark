package utils

import (
	"fmt"
	"strings"
)

// Check filename file format
func FilenameCheck(fileName *string, numRequests *int, concurrency *int) *string {
	if strings.HasSuffix(*fileName, ".xlsx") {
		return fileName
	} else {
		newFileName := fmt.Sprintf("%s-r%d-c%d.xlsx", *fileName, *numRequests, *concurrency)
		return &newFileName
	}
}
