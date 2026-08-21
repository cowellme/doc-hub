package service

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

func buildStoredFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	base := strings.TrimSuffix(originalName, ext)
	base = strings.ReplaceAll(strings.ToLower(base), " ", "_")

	return fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), base, ext)
}
