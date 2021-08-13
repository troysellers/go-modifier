package mockaroo

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/troysellers/go-modifier/mockaroo/types"
)

func TestNothing(t *testing.T) {

	schema := []types.IField{}
	firstName := make(map[string]interface{})
	firstName["name"] = "FirstName"
	schema = append(schema, types.NewFirstName(firstName))
	lastName := make(map[string]interface{})
	lastName["name"] = "LastName"
	schema = append(schema, types.NewLastName(lastName))
	b, err := json.Marshal(schema)
	if err != nil {
		t.Fatal(err)
	}
	url := fmt.Sprintf("https://api.mockaroo.com/api/generate.csv?key=%v&count=%d&include_header=false", "c04c9a30", 5000)
	filePath := "/tmp/mockaroo-data/names.csv"

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	i := 0
	for {
		_, respbytes, err := doHttp(url, "", b, "POST", nil)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = f.Write(respbytes); err != nil {
			t.Fatal(err)
		}
		i++
		if i == 200 {
			break
		}
	}
	log.Printf("written to %v\n", filePath)

}
