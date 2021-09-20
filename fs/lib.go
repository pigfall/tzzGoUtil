package fs

import (
	"errors"
	tz_filepath "github.com/pigfall/tzzGoUtil/path/filepath"
	"io/ioutil"
	"os"
	//"path"
)

func CreateThen(filepath string, then func(file *os.File) error) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	return then(file)
}
func ReadAllThen(filepath string, then func(content []byte) error) error {
	bytes, err := ReadFile(filepath)
	if err != nil {
		return err
	}
	return then(bytes)
}

func OpenThen(filepath string, then func(file *os.File) error) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	return then(file)
}

func FileExist(filepath string) (bool, error) {
	_, err := os.Stat(filepath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil

}

func ReadFile(filepath string) ([]byte, error) {
	return ioutil.ReadFile(filepath)
}

func IsNotExist(err error) bool {
	return errors.Is(err, os.ErrNotExist)
}

func WriteFile(filepath string, bytes []byte) error {
	dir := tz_filepath.DirPath(filepath)
	_, err := os.Stat(dir)
	if err != nil {
		if IsNotExist(err) {
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return ioutil.WriteFile(filepath, bytes, os.ModePerm)
}
