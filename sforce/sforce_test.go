package sforce

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/troysellers/go-modifier/config"
	"github.com/tzmfreedom/go-soapforce"
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Print("No .env file found")
	}
}

func TestSyncMap(t *testing.T) {
	var sm sync.Map

	key := "USER"
	val1 := "some value"

	v, ok := sm.Load(key)
	fmt.Printf("empty load should be empty [%v] [%v]", v, ok)

	sm.Store(key, val1)

	v, ok = sm.Load(key)
	fmt.Printf("empty load should be populated [%v] [%v]", v, ok)
}

func TestGetObjFromQuery(t *testing.T) {
	query := "select id, name, descriptions from account where field = value"
	obj := getObjectNameFromQuery(query)
	if !strings.EqualFold(obj, "account") {
		t.Error()
	}
}

func TestNewSoapClient(t *testing.T) {
	cfg := config.NewConfig()
	c, err := NewSoapClient(&cfg.SF)
	if err != nil {
		t.Fatal(err)
	}
	if c == nil {
		t.Fatal("Returned client is nil")
	}
	if c.SessionId == "" {
		t.Fatal("no session id")
	}
}
func TestNewRestClient(t *testing.T) {
	cfg := config.NewConfig()
	c, err := NewRestClient(&cfg.SF)
	if err != nil {
		t.Fatal(err)
	}
	if c == nil {
		t.Fatal("Returned client is nil")
	}
}

func TestGetBatches(t *testing.T) {
	jobId := "7507h000006y7MSAAY"
	log.Printf("%v", jobId)
}

func TestDate(t *testing.T) {

	d := "2006-01-02T15:04:05+0000"
	d2 := "2021-07-27T03:29:17.000+0000"
	t2, _ := time.Parse(d, d2)
	fmt.Printf("%v", t2)

}

func TestGetStuff(t *testing.T) {
	cfg := &config.Config{
		SF: config.SFConfig{
			Username:   "troy@grax.perf",
			Password:   "demo1234",
			LoginUrl:   "login.salesforce.com",
			Token:      "",
			ApiVersion: 51.0,
			SfDebug:    true,
		},
	}
	c, err := NewSoapClient(&cfg.SF)
	if err != nil {
		t.Fatal(err)
	}
	dsor, err := c.DescribeSObject("Account")
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range dsor.Fields {
		//log.Printf("%v\n", f.Name)
		if f.Name == "CreatedDate" || f.Name == "CreatedById" || f.Name == "LastModifiedDate" || f.Name == "LastModifiedById" {
			log.Printf("%s.updateable : %v", f.Name, f.Updateable)
		}
	}

}
func TestWriteCSV(t *testing.T) {
	cfg := &config.Config{
		SF: config.SFConfig{
			Username:   "troy@grax.sdo.allRecords",
			Password:   "Demo1234",
			LoginUrl:   "test.salesforce.com",
			Token:      "",
			ApiVersion: 52.0,
			SfDebug:    false,
		},
	}
	c, err := NewRestClient(&cfg.SF)
	if err != nil {
		t.Fatal(err)
	}
	UploadCSVToSalesforce(cfg, c, "/tmp/mockaroo-data/account-update.csv", "Account")
}

func TestCreateAccount(t *testing.T) {

	cfg := &config.Config{
		SF: config.SFConfig{
			Username:   "troy@grax.perf",
			Password:   "demo1234",
			LoginUrl:   "login.salesforce.com",
			Token:      "",
			ApiVersion: 51.0,
			SfDebug:    true,
		},
	}
	c, err := NewSoapClient(&cfg.SF)
	if err != nil {
		t.Fatal(err)
	}
	fields := make(map[string]interface{})
	fields["Name"] = "Test-SOAP-Troy"
	fields["CreatedDate"] = "2021-08-03T12:00:00.000Z"
	fields["LastModifiedDate"] = "2021-08-03T12:00:00.000Z"
	fields["LastModifiedById"] = "0056g0000022cohAAA"
	fields["CreatedById"] = "0056g0000022cohAAA"
	obj := &soapforce.SObject{
		Fields: fields,
		Type:   "Account",
		Id:     "0016g00000IGY2xAAH",
	}
	var objects []*soapforce.SObject
	objects = append(objects, obj)
	sr, err := c.Update(objects)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range sr {
		if !s.Success {
			for _, e := range s.Errors {
				log.Printf("%v\n", e.Message)
			}
		}
	}
}
