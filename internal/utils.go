package internal

import (
	"strings"

	"github.com/spf13/afero"
)

func Exists(fs afero.Fs, path string) (bool, error) {
	return afero.Exists(fs, path)
}

// NormalizeEnvKey converts a key to an environment variable format
func NormalizeEnvKey(key string) string {
	return strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
}
