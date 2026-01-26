package utils

import (
	"path/filepath"
)

// buildPathFromChecksum: slash 12 first bytes of checksum with every 4 characters
func BuildPathFromChecksum(checksum, fileMimeType string) string {
	if len(checksum) < 12 {
		return checksum + fileMimeType
	}

	parts := []string{
		checksum[0:4],
		checksum[4:8],
		checksum[8:12],
		checksum + fileMimeType,
	}
	return filepath.Join(parts...)
}
