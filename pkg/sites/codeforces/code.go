package codeforces

import (
	"codeforces_cli/pkg/runner"
	"fmt"
	"os"
	"path/filepath"
)

// Language detection by extension
var languageExtensions = map[string]string{
	".cpp": GNUGPP17,
	".go":  GO,
	".py":  PYTHON3,
}

func NewCode(filePath string) (*runner.Code, error) {
	// Expand environment variables
	filePath = os.ExpandEnv(filePath)
	ext := filepath.Ext(filePath)
	language, exists := languageExtensions[ext]

	if !exists {
		return nil, fmt.Errorf("Unknown extension %s", ext)
	}
	return &runner.Code{
		FilePath: filePath,
		Language: language,
	}, nil
}
