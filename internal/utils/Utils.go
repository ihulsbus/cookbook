package utils

import "os"

func InitFolder(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		err := os.Mkdir(path, 0750)
		if err != nil {
			return err
		}
		return nil
	}
	return err
}
