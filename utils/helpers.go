package utils

import (
	"fmt"
	"strings"
)

// Check filename file format
func FilenameCheck(fileName *string) *string {
	if strings.HasSuffix(*fileName, ".xlsx") {
		return fileName
	} else {
		newFileName := fmt.Sprintf("%s.xlsx", *fileName)
		return &newFileName
	}
}
