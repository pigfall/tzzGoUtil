package fs

import (
	"errors"
	"io/ioutil"
	"os"
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
	bytes, err := ioutil.ReadFile(filepath)
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
