package json

import (
	stdjson "encoding/json"
	"os"
)

func EncodeJsonToFile(filename string, jsons ...interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close() // Ignore error.
	encoder := stdjson.NewEncoder(file)
	for _, j := range jsons {
		err = encoder.Encode(j)
		if err != nil {
			return err
		}
	}
	return nil
}

func EncodeJsonToFileIfNotExist(filename string, jsons ...interface{}) (hasWritten bool, err error) {
	if _, err = os.Stat(filename); err == nil {
		return false, nil
	}
	if !os.IsNotExist(err) {
		return false, err
	}
	err = EncodeJsonToFile(filename, jsons...)
	return err == nil, err
}

func DecodeJsonFromFile(filename string, out ...interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close() // Ignore error.
	decoder := stdjson.NewDecoder(file)
	for _, o := range out {
		err = decoder.Decode(o)
		if err != nil {
			return err
		}
	}
	return nil
}

func DecodeJsonFromFileIfExist(filename string, out ...interface{}) (hasRead bool, err error) {
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		// File does NOT exist.
		return false, nil
	}
	err = DecodeJsonFromFile(filename, out...)
	return err == nil, err
}
