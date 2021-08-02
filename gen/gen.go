package gen

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/simpleforce/simpleforce"
	"github.com/troysellers/go-modifier/config"
	"github.com/troysellers/go-modifier/lorem"
)

func GetValueForType(f map[string]interface{}, c *simpleforce.Client) (interface{}, error) {

	// if can be empty, retun empty on a 10%
	if f["nillable"].(bool) && rand.Intn(10) < 2 {
		return nil, nil
	}
	switch f["type"].(string) {
	case "id":
		return nil, fmt.Errorf("id values are not supported for generation")
	case "boolean":
		return rand.Intn(10) >= 5, nil
	case "string", "encryptedstring":
		l := f["length"].(float64)
		return lorem.Word(1, rand.Intn(int(l))), nil
	case "datetime", "date":
		d := time.Now()
		d = d.AddDate(0, rand.Intn(12), rand.Intn(30))
		return d, nil
	case "reference":
		//TODO : get random related id somehow...
		return nil, fmt.Errorf("reference not implemented yet")
	case "currency", "double":
		p := f["precision"].(float64)
		s := f["scale"].(float64)
		return rand.Intn(int(p)) / int(math.Pow10(int(s))), nil
	case "email":
		return lorem.Email(), nil
	case "location":
		return nil, fmt.Errorf("location value not implemented yet")
	case "percent":
		return float32(rand.Intn(100)), nil
	case "phone":
		return nil, fmt.Errorf("phone value not implemented yet")
	case "picklist", "multipicklist":
		plv := f["picklistValues"].([]interface{})
		selected := plv[rand.Intn(len(plv))]
		return selected.(map[string]interface{})["value"], nil
	case "textarea":
		l := f["length"].(int)
		s := lorem.Sentence(1, l)
		if len(s) < int(l) {
			return s, nil
		} else {
			return s[:l], nil
		}
	case "time":
		return nil, fmt.Errorf("phone value not implemented yet")
	case "url":
		return lorem.Url(), nil
	}

	return nil, nil
}

// pass a mockaroo schema
// receive a path to the downloaded CSV file
func GetDataFromMockaroo(cfg *config.MockarooConfig, s string, r int) (string, error) {

	resp, err := http.Get(fmt.Sprintf("https://api.mockaroo.com/api/generate.csv?key=%v&count=%v&schema=%v", cfg.Key, r, s))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	fName := getFileName(s)
	out, err := os.Create(fName)
	if err != nil {
		fmt.Printf("err creating the file %v\n", err)
		return "", err
	}
	defer out.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		_, copyError := io.Copy(out, resp.Body)
		if copyError != nil {
			return "", err
		}
		return out.Name(), nil
	} else {
		return "", fmt.Errorf("unhandled response code from call to mockaro %v", resp.Status)
	}
}

// builds a file name ~/mockaroo-data/<schema>.csv
// and returns.
// if a file exists at this location it renames it using timestamp
// ~/mockaroo-data/<schema>-<timestamp>.csv
func getFileName(s string) string {
	dataDir, _ := os.UserHomeDir()
	dataDir = fmt.Sprintf("%v/mockaroo-data", dataDir)
	_, err := os.Stat(fmt.Sprintf("%v/%v.csv", dataDir, s))
	if err == nil {
		os.Rename(fmt.Sprintf("%v/%v.csv", dataDir, s), fmt.Sprintf("%v/%v-%v.csv", dataDir, s, time.Now().Unix()))
	}
	return fmt.Sprintf("%v\n", dataDir)
}
