package file

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
)

// reads the entire file into a byte []
func GetCSVBytes(file string) ([]byte, error) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// takes the data and writes is as a CSV to the filename
// returns the fully qualified filename
func WriteCsv(fileNmae string, data [][]string) (string, error) {

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	f, err := os.Create(fmt.Sprintf("%v/tmpdata/%v", home, fileNmae))
	if err != nil {
		return "", err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	for _, row := range data {
		_ = w.Write(row)
	}
	w.Flush()

	return f.Name(), nil
}
