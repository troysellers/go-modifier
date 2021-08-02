package gen

import (
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/troysellers/go-modifier/config"
	"github.com/troysellers/go-modifier/sforce"
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Print("No .env file found")
	}
}

func TestGetFieldValue(t *testing.T) {
	cfg := config.NewConfig()
	c, err := sforce.NewSoapClient(&cfg.SF)
	if err != nil {
		t.Fatal(err)
	}
	dr, err := c.DescribeSObject("every_field__c")
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range dr.Fields {
		if *f.Type_ != "id" {
			val, err := GetValueForType(f, c)
			if err != nil {
				t.Fatal(err)
			}
			log.Printf("Type [%v] Val [%v]\n", *f.Type_, val)
		}
	}
}

func TestGetMockaroo(t *testing.T) {
	ms := []string{"account", "contact", "case", "task", "event", "casecomment", "lead", "opportunity", "personaccount"}
	cfg := config.NewConfig()
	for _, s := range ms {
		f, err := GetDataFromMockaroo(&cfg.Mockaroo, s, 10)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("We donwloaded the file to %v\n", f)
	}

}
