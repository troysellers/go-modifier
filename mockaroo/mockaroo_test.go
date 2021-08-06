package mockaroo

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/troysellers/go-modifier/config"
)

func TestStructs(t *testing.T) {

	g := NewGUID("My GUID")
	b, err := json.Marshal(g)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("%s", string(b))
}

func TestArray(t *testing.T) {
	var fields []interface{}
	fields = append(fields, NewGUID("myGUid"))
	fields = append(fields, NewCustomList("mycustomerlist"))
	fields = append(fields, NewFakeCompanyName("companyname"))
	b, err := json.Marshal(fields)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("%s", string(b))
}

func TestField(t *testing.T) {
	cfg := &config.Config{
		Mockaroo: config.MockarooConfig{
			Key: "c04c9a30",
		},
	}
	var fields []interface{}
	w := NewWords("mynewwords")
	w.Formula = "this[0,255]"
	fields = append(fields, w)

	b, err := json.Marshal(fields)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("%v\n", string(b))
	url := fmt.Sprintf("https://api.mockaroo.com/api/generate.csv?key=%v&count=%d", cfg.Mockaroo.Key, 20)
	header, respbytes, err := doHttp(url, "", b, "POST", nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("%v", header)
	log.Printf("%s", string(respbytes))
}

func TestAllFields(t *testing.T) {
	cfg := &config.Config{
		Mockaroo: config.MockarooConfig{
			Key: "c04c9a30",
		},
	}
	var fields []interface{}
	fields = append(fields, NewDUNSNumber("mynewDUNS"))
	customList := NewCustomList("mynewcustomlist")
	customList.Values = []string{"val1", "val2", "val3", "val4"}
	fields = append(fields, customList)
	fields = append(fields, NewDatetime("mynewdatetime"))
	digitSequence := NewDigitSequence("mydigisequence")
	digitSequence.Format = "###-###-###"
	fields = append(fields, digitSequence)
	fields = append(fields, NewFakeCompanyName("myfakecompany"))
	fields = append(fields, NewNumber("mynewnumber"))
	fields = append(fields, NewPhone("mynewphone"))
	fields = append(fields, NewPostalCode("mynewpostal"))
	fields = append(fields, NewGUID("mynewguid"))
	fields = append(fields, NewSentences("mynewsentences"))
	fields = append(fields, NewState("mynewstate"))
	fields = append(fields, NewStreetAddress("mynewstreetaddress"))
	fields = append(fields, NewStreetName("mynewstreetname"))
	fields = append(fields, NewURL("mynewurl"))
	fields = append(fields, NewWords("mynewwords"))
	fields = append(fields, NewBoolean("mynewbool"))
	fields = append(fields, NewBuzzword("mynewbuzz"))
	fields = append(fields, NewCatchPhrase("mynewcatch"))
	fields = append(fields, NewCity("mynewcity"))
	fields = append(fields, NewCountry("mynewcountry"))

	b, err := json.Marshal(fields)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("%v\n", string(b))
	url := fmt.Sprintf("https://api.mockaroo.com/api/generate.csv?key=%v&count=%d", cfg.Mockaroo.Key, 20)
	header, respbytes, err := doHttp(url, "", b, "POST", nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("%v", header)
	log.Printf("%s", string(respbytes))
	fName := "/tmp/mockaroo-data/test.csv"
	os.WriteFile(fName, respbytes, 0766)
}
