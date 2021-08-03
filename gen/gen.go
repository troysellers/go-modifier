package gen

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/troysellers/go-modifier/config"
	"github.com/troysellers/go-modifier/file"
)

// pass a mockaroo schema
// receive a path to the downloaded CSV file
func GetDataFromMockaroo(cfg *config.MockarooConfig, s string, r int) (string, error) {

	resp, err := http.Get(fmt.Sprintf("https://api.mockaroo.com/api/generate.csv?key=%v&count=%v&schema=%v", cfg.Key, r, s))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	fPath, err := file.BuildFilePath(fmt.Sprintf("%v.csv", s))
	if err != nil {
		return "", err
	}
	out, err := os.Create(fPath)
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
