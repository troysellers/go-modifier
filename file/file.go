package file

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
)

// reads the entire file into a byte []
func GetCSVBytes(file string) ([]byte, error) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func BuildFilePath(f string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v/tmpdata/%v", home, f), nil
}

// takes the data and writes is as a CSV to the filename
// returns the fully qualified filename
func WriteCsv(filePath string, data [][]string) (string, error) {

	f, err := os.Create(filePath)
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

// randomly popluate a column with values
// writes these changes to the file
// will append the column if it doesn't exist in the file.
func UpdateColumn(filePath string, col string, ids []string) error {
	if filePath == "" || col == "" || ids == nil {
		return nil
	}
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	data, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return err
	}
	var colIndex int = -1
	var appending bool
	for i, c := range data[0] { // loop through the header
		if strings.EqualFold(c, col) {
			colIndex = i
		}
	}

	if colIndex < 0 {
		appending = true
		data[0] = append(data[0], col)
	}
	// loop through all rows
	for i, _ := range data[1:] {
		if appending {
			data[i+1] = append(data[i+1], ids[rand.Intn(len(ids))])
		} else {
			data[i+1][colIndex] = ids[rand.Intn(len(ids))]
		}
	}
	if _, err := WriteCsv(filePath, data); err != nil {
		return err
	}
	return nil
}
