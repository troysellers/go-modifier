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
