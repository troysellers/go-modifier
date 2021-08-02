package config

import (
	"log"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Print("No .env file found")
	}
}

func TestGetConfig(t *testing.T) {
	cfg := NewConfig()
	q := make([]string, 2)
	q[0] = "select Id, Name,Industry, Website, YearStarted,Tradestyle, Phone, NumberOfEmployees, Description, Fax from Account where isPersonAccount=false limit 50000"
	q[1] = "select Id, StageName, Amount, CloseDate from Opportunity where isClosed = false"

	testCfg := &Config{
		Mockaroo: MockarooConfig{
			Key: "mock_key",
		},
		SF: SFConfig{
			Username:    "sforce@user.com",
			Password:    "sforcepass",
			Token:       "sforcetoken",
			LoginUrl:    "test.salesforce.com",
			ApiVersion:  52.0,
			SfDebug:     false,
			SfBatchSize: 200,
			Queries:     q,
		},
	}
	log.Printf("%v", cfg)
	log.Printf("%v", testCfg)
	if !cmp.Equal(testCfg, cfg) {
		t.Error("config is not equal")
	}
}
