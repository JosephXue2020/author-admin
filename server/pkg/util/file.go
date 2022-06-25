package util

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
)

func LoadGob(p string, x interface{}) error {
	if !PathExist(p) {
		err := fmt.Errorf("File not exist: %s\n", p)
		return err
	}

	data, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	err = decoder.Decode(x)
	if err != nil {
		err := fmt.Errorf("Filed to decode to type any.")
		return err
	}

	return nil
}

func SaveGob(p string, x interface{}) (int, error) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(x)
	if err != nil {
		return 0, err
	}

	byteData := buf.Bytes()
	err = ioutil.WriteFile(p, byteData, 0666)
	if err != nil {
		return 0, err
	}

	return len(byteData), nil
}
