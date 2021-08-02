package gen

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/troysellers/go-modifier/config"
	"github.com/troysellers/go-modifier/lorem"
	"github.com/tzmfreedom/go-soapforce"
)

func GetValueForType(f *soapforce.Field, c *soapforce.Client) (interface{}, error) {

	// if can be empty, retun empty on a 10%
	if f.Nillable && rand.Intn(10) < 2 {
		return nil, nil
	}
	switch *f.Type_ {
	case "id":
		return nil, fmt.Errorf("id values are not supported for generation")
	case "boolean":
		return rand.Intn(10) >= 5, nil
	case "string", "encryptedstring":
		return lorem.Word(1, rand.Intn(int(f.Length))), nil
	case "datetime", "date":
		return time.Now(), nil
	case "reference":
		return GetRelatedId(f.ReferenceTo[0], c)
	case "currency", "double":
		return rand.Intn(int(f.Precision)) / int(math.Pow10(int(f.Scale))), nil
	case "email":
		return lorem.Email(), nil
	case "location":
		return nil, fmt.Errorf("location value not implemented yet")
	case "percent":
		return float32(rand.Intn(100)), nil
	case "phone":
		return nil, fmt.Errorf("phone value not implemented yet")
	case "picklist":
		return f.PicklistValues[rand.Intn(len(f.PicklistValues))].Value, nil
	case "multipicklist":
		return f.PicklistValues[rand.Intn(len(f.PicklistValues))].Value, nil
	case "textarea":
		l := f.Length
		s := lorem.Sentence(1, int(l))
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

// two queries happen
// first gets the one record, then look at the
// total records size param.
// then random offset using this number
func GetRelatedId(obj string, c *soapforce.Client) (string, error) {
	// first query gets a total
	q := fmt.Sprintf("select id from %v", "account")
	qr, err := c.QueryAll(q)
	if err != nil {
		return "", err
	}
	q2 := fmt.Sprintf("select id from %v limit 1 offset %d", obj, rand.Intn(int(qr.Size)))
	fmt.Printf("%v\n", q2)
	qr2, err := c.Query(q2)
	if err != nil {
		return "", err
	}
	return qr2.Records[0].Id, nil
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
		return "", fmt.Errorf("unhandled response code from call to mockaro %v\n", resp.Status)
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
