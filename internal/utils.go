package internal

import "github.com/spf13/afero"

func Exists(fs afero.Fs, path string) (bool, error) {
	return afero.Exists(fs, path)
}
