package mockaroo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/simpleforce/simpleforce"
	"github.com/troysellers/go-modifier/config"
	"github.com/troysellers/go-modifier/mockaroo/types"
)

//Output formats
const JSONFormat string = "generate.json"
const CSVFormat string = "generate.csv"
const TXTFormat string = "generate.txt"
const CustomFormat string = "generate.custom"
const SQLFormat string = "generate.sql"
const XMLFormat string = "generate.xml"

type MockarooRequest struct {
	SObject        *simpleforce.SObjectMeta
	Cfg            *config.Config
	Count          int
	PersonAccounts bool
	Schema         []types.IField
	FilePath       string
}

// fetches mockaroo CSV for the object specified.
// returns a string that is the full path
func (r *MockarooRequest) GetDataForObj() error {

	r.Schema = getSchemaForObjectType(r.SObject, r.PersonAccounts)

	b, err := json.Marshal(r.Schema)
	if err != nil {
		return err
	}

	header := true
	// mockaroo has a 5000 record api limit.
	mockLimit := 1000
	var wg sync.WaitGroup
	var files sync.Map
	var index int

	numBatches := r.Count / mockLimit

	for i := 1; i <= numBatches; i++ {
		log.Printf("fetching %d to %d dummy data\n", (i-1)*mockLimit, i*mockLimit)
		wg.Add(1)
		fname := fmt.Sprintf("%v%v-%d.csv", r.Cfg.Mockaroo.DataDir, (*r.SObject)["name"].(string), i)
		go fetchMockarooBatch(fname, r.Cfg.Mockaroo.Key, mockLimit, b, header, &wg, &files, i)

		if header {
			header = false
		}
		index = i
		if math.Mod(float64(i), 4) == 0 {
			wg.Wait()
		}
	}
	// mod gives us the remaining records to get.
	remainder := int(math.Mod(float64(r.Count), float64(mockLimit)))
	if remainder > 0 {
		fname := fmt.Sprintf("%v%v-%d.csv", r.Cfg.Mockaroo.DataDir, (*r.SObject)["name"].(string), index)
		wg.Add(1)
		go fetchMockarooBatch(fname, r.Cfg.Mockaroo.Key, remainder, b, header, &wg, &files, index)
	}

	wg.Wait()
	r.FilePath, err = mergeFiles(&files, r.Cfg.Mockaroo.DataDir, (*r.SObject)["name"].(string))
	if err != nil {
		return err
	}
	return nil
}

func mergeFiles(files *sync.Map, dir string, obj string) (string, error) {

	final, err := os.Create(fmt.Sprintf("%v%v.csv", dir, obj))
	if err != nil {
		return "", err
	}
	defer final.Close()
	var keys []int

	files.Range(func(k, s interface{}) bool {
		keys = append(keys, k.(int))
		return true
	})

	sort.Ints(keys)
	for _, k := range keys {
		s, _ := files.Load(k)
		log.Printf("merge %v\n", s)
		f, err := os.Open(s.(string))
		if err != nil {
			return "", err
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return "", err
		}
		final.Write(b)
	}
	return final.Name(), nil
}

func setFormula(f *types.Field) {
	l := int(f.SforceMeta["length"].(float64))
	if l > 0 {
		if l > 1000 { //lets not go crazy with text
			l = 1000
		}
		f.Formula = fmt.Sprintf("if this.nil? then '' else this[0,%d] end", l)
	}
}

func fetchMockarooBatch(fname string, key string, records int, schema []byte, header bool, wg *sync.WaitGroup, files *sync.Map, mapkey int) {

	log.Printf("Fetching mockaroo schema \n%v\n", string(schema))
	defer wg.Done()
	f, err := os.Create(fname)
	if err != nil {
		log.Printf("%v", err)
	}
	defer f.Close()

	headers := make(map[string]string)
	headers["Accept"] = "application/json"
	headers["Content-Type"] = "application/json"
	url := fmt.Sprintf("https://api.mockaroo.com/api/generate.csv?key=%v&count=%d&include_header=%v", key, records, header)

	_, respbytes, err := doHttp(url, "", schema, "POST", headers)
	if err != nil {
		log.Printf("%v", err)
		panic(err)
	}

	if _, err := f.Write(respbytes); err != nil {
		panic(err)
	}
	files.Store(mapkey, f.Name())
}

// captures top level logic on whether the field data should
// be generated.
// Returns true if all these conditions are satisfied
// - updateable = true
// - doesn't belong to a managed package (i.e. name doesn't contain two occurences of '__' )
func shouldGetData(f map[string]interface{}) bool {

	// we only want to get data for fields we can update
	if !f["updateable"].(bool) {
		return false
	}
	// we aren't going to try to populate managed package data fields
	if strings.Count(f["name"].(string), "__") == 2 {
		return false
	}
	return true
}

func getMockTypeForField(f map[string]interface{}) types.IField {

	// get the mock type first
	var mockType types.IField
	switch f["type"].(string) {
	case "id":
		return nil
	case "boolean":
		mockType = types.NewBoolean(f)
	case "string", "encryptedstring":
		if f["externalId"].(bool) || f["unique"].(bool) {
			mockType = types.NewGUID(f)
		} else {
			mockType = types.NewWords(f)
		}
	case "datetime", "date":
		mockType = types.NewDatetime(f)
	case "reference":
		w := types.NewWords(f)
		w.Max = 0
		w.Min = 0
		mockType = w
	case "currency", "double", "percent", "int":
		p := int(f["precision"].(float64))
		s := int(f["scale"].(float64))
		n := types.NewNumber(f)
		n.Decimals = s
		n.Max = p*10 - 1
		mockType = n
	case "email":
		mockType = types.NewEmailAddress(f)
	case "phone":
		mockType = types.NewPhone(f)
	case "picklist", "multipicklist":
		plv := f["picklistValues"].([]interface{})
		l := types.NewCustomList(f)
		for _, v := range plv {
			val := v.(map[string]interface{})
			l.Values = append(l.Values, val["value"].(string))
		}
		mockType = l
	case "textarea":
		t := types.NewSentences(f)
		t.Max = 100
		t.Min = 1
		mockType = t
	case "url":
		mockType = types.NewURL(f)
	default:
		log.Printf("%v type has not been mapped to a mockaroo data type yet", f["type"])
		return nil
	}
	return mockType

}

// returns a mocktype that is ideal for the object
// or defaults for custom object
func getSchemaForObjectType(obj *simpleforce.SObjectMeta, personAccounts bool) []types.IField {

	var schema []types.IField
	var fields = (*obj)["fields"].([]interface{})
	switch (*obj)["name"].(string) {
	case "Account", "account":
		schema = getSchemaForAccount(fields, personAccounts)
	case "Contact", "contact":
		schema = getSchemaForContact(fields, personAccounts)
	case "Case", "case":
		schema = getSchemaForCase(fields, personAccounts)
	case "Opportunity", "opportunity":
		schema = getSchemaForOpportunity(fields, personAccounts)
	case "Lead", "lead":
		schema = getSchmeaForLead(fields, personAccounts)
	case "Task", "task":
		schema = getSchemaForTask(fields, personAccounts)
	case "Event", "event":
		schema = getSchemaForEvent(fields, personAccounts)
	default:
		schema = getSchemaForGenericObj(fields)
	}
	for _, f := range schema {
		setFormula(f.GetField())
	}
	return schema
}

// returns the mockaroo schema for any object we haven't
// specifically coded for.
func getSchemaForGenericObj(fields []interface{}) []types.IField {

	var mockFields []types.IField
	for _, f := range fields {
		field := f.(map[string]interface{})
		if shouldGetData(field) {
			mf := getMockTypeForField(field)
			l := int(field["length"].(float64))
			if l > 0 {
				if l > 1000 { //lets not go crazy with text
					l = 1000
				}
				mf.SetFormula(fmt.Sprintf("this.nil? then '' else this[0,%d] end", l))
			}
			mockFields = append(mockFields, mf)
		}
	}
	return mockFields
}

func DoHttp(url string, sid string, body []byte, method string, headers map[string]string) (http.Header, []byte, error) {
	return doHttp(url, sid, body, method, headers)
}

// returns the response body bytes if we had a 200 response.
// errors for all others.
func doHttp(url string, sid string, body []byte, method string, headers map[string]string) (http.Header, []byte, error) {

	defer trackTime(time.Now(), fmt.Sprintf("%v:%v", method, url))

	client := &http.Client{}
	var r *bytes.Reader
	if body != nil {
		r = bytes.NewReader(body)
	} else {
		r = bytes.NewReader(nil)
	}
	req, err := http.NewRequest(method, url, r)
	if err != nil {
		return nil, nil, err
	}
	for header, value := range headers {
		req.Header.Add(header, value)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}
	if !(res.StatusCode >= 200 && res.StatusCode < 300) {
		return res.Header, nil, fmt.Errorf("unsuccesful attempt to call endpoint %v\n%v", res.Status, string(bytes))
	}

	log.Printf("Received %d bytes with http response %v\n", len(bytes), res.Status)
	return res.Header, bytes, nil
}

func trackTime(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s\n", name, elapsed)
}
