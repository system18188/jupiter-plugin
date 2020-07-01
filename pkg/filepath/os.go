package filepath

import "os"

// MkDirAll method creates nested directories with given permission if not exists
func MkDirAll(path string, mode os.FileMode) error {
	if _, err := os.Lstat(path); err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(path, mode); err != nil {
				return err
			}
		}
	}
	return nil
}
