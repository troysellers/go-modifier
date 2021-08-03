package file

import (
	"encoding/csv"
	"fmt"
	"os"
	"testing"
)

func TestOpen(t *testing.T) {

	f, err := os.Open("/Users/troysellers/tmpdata/opportunity.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		t.Fatal("Unable to parse file as CSV for /Users/troysellers/tmpdata/opportunity.csv", err)
	}
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("We have %d lines\n", len(records))
}
