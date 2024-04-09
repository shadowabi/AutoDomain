package File

import (
	"errors"
	"os"
)

func DirNonExistCreate(dirpath string) error {
	exist, err := IsDirExist(dirpath)
	if exist {
		return nil
	}
	if err == nil {
		return os.MkdirAll(dirpath, 0755)
	}
	return err
}

func FileNonExistCreate(filepath string) error {
	exist, err := IsFileExist(filepath)
	if exist {
		return nil
	}
	if err == nil {
		f, err := os.Create(filepath)
		if err != nil {
			return err
		}
		return f.Close()
	}
	return err
}

func IsFileExist(filepath string) (bool, error) {
	if s, err := os.Stat(filepath); err == nil {
		// path/to/whatever exists
		if s.IsDir() {
			return false, errors.New(filepath + " is a dir")
		}
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does *not* exist
		return false, nil
	} else {
		return false, err
	}
}

func IsDirExist(dirpath string) (bool, error) {
	if s, err := os.Stat(dirpath); err == nil {
		// path/to/whatever exists
		if s.IsDir() {
			return true, nil
		}
		return false, errors.New(dirpath + " is a file")
	} else if errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does *not* exist
		return false, nil
	} else {
		return false, err
	}
}
