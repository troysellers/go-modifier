package mockaroo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/simpleforce/simpleforce"
	"github.com/troysellers/go-modifier/config"
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
}

// fetches mockaroo CSV for the object specified.
// returns a string that is the full path
func GetDataForObj(req *MockarooRequest) (string, error) {

	mockFields := getSchemaForObjectType(req.SObject, req.PersonAccounts)
	b, err := json.Marshal(mockFields)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("https://api.mockaroo.com/api/generate.csv?key=%v&count=%d", req.Cfg.Mockaroo.Key, req.Count)
	_, respbytes, err := doHttp(url, "", b, "POST", nil)
	if err != nil {
		return "", err
	}
	fName := fmt.Sprintf("%v%v.csv", req.Cfg.Mockaroo.DataDir, (*req.SObject)["name"])
	if err = os.WriteFile(fName, respbytes, 0777); err != nil {
		return "", nil
	}
	return fName, nil
}

// captures top level logic on whether the field data should
// be generated.
// Returns true if all these conditions are satisfied
// - updateable = true
// - doesn't belong to a managed package (i.e. name doesn't contain two occurences of '__' )
// - field isn't self referencing (i.e. name != 'ParentId')
func shouldGetData(f map[string]interface{}) bool {

	// we only want to get data for fields we can update
	if !f["updateable"].(bool) {
		return false
	}
	// we aren't going to try to populate managed package data fields
	if strings.Count(f["name"].(string), "__") == 2 {
		return false
	}
	// we don't support parent relationships yet
	if f["name"].(string) == "ParentId" {
		return false
	}
	return true
}

func getMockTypeForField(f map[string]interface{}) interface{} {

	// get the mock type first
	var mockType interface{}
	switch f["type"].(string) {
	case "id":
		return nil
	case "boolean":
		mockType = NewBoolean(f["name"].(string))
	case "string", "encryptedstring":
		if f["externalId"].(bool) {
			e := NewGUID(f["name"].(string))
			e.Formula = fmt.Sprintf("this[0,%d]", int(f["length"].(float64)))
			mockType = e
		} else {
			w := NewWords(f["name"].(string))
			w.Formula = fmt.Sprintf("this[0,%d]", int(f["length"].(float64)))
			mockType = w
		}
	case "datetime", "date":
		d := NewDatetime(f["name"].(string))
		mockType = d
	case "reference":
		w := NewWords(f["name"].(string))
		w.Max = 0
		w.Min = 0
		mockType = w
	case "currency", "double", "percent", "int":
		p := int(f["precision"].(float64))
		s := int(f["scale"].(float64))
		n := NewNumber(f["name"].(string))
		n.Decimals = s
		n.Max = p*10 - 1
		mockType = n
	case "email":
		e := NewEmailAddress(f["name"].(string))
		mockType = e
	case "phone":
		p := NewPhone(f["name"].(string))
		mockType = p
	case "picklist", "multipicklist":
		plv := f["picklistValues"].([]interface{})
		l := NewCustomList(f["name"].(string))
		for _, v := range plv {
			val := v.(map[string]interface{})
			l.Values = append(l.Values, val["value"].(string))
		}
		mockType = l
	case "textarea":
		t := NewSentences(f["name"].(string))
		t.Max = 100
		t.Min = 1
		t.Formula = fmt.Sprintf("this[0,%d]", int(f["length"].(float64)))
		mockType = t
	case "url":
		u := NewURL(f["name"].(string))
		mockType = u
	default:
		log.Printf("%v type has not been mapped to a mockaroo data type yet", f["type"])
		return nil
	}
	return mockType

}

// returns a mocktype that is ideal for the object
// or defaults for custom object
func getSchemaForObjectType(obj *simpleforce.SObjectMeta, personAccounts bool) []interface{} {

	var fields = (*obj)["fields"].([]interface{})
	switch (*obj)["name"].(string) {
	case "Account":
		log.Println("getting data for the account")
		return getSchemaForAccount(fields, personAccounts)
	case "Contact":
		return getSchemaForContact(fields, personAccounts)
	default:
		return getSchemaForGenericObj(fields)
	}
}

// returns the mockaroo schema for any object we haven't
// specifically coded for.
func getSchemaForGenericObj(fields []interface{}) []interface{} {

	var mockFields []interface{}
	for _, f := range fields {
		field := f.(map[string]interface{})
		if shouldGetData(field) {
			mockFields = append(mockFields, getMockTypeForField(field))
		}
	}
	return mockFields
}

// returns true if the string is set as the Name value of one of the
// schema elements passed in the array
func hasBeenPopulated(fieldName string, schema []FieldSpec) bool {

	for _, s := range schema {
		if s.Name == fieldName {
			return true
		}
	}
	return false
}

// returns the response body bytes if we had a 200 response.
// errors for all others.
func doHttp(url string, sid string, body []byte, method string, headers map[string]string) (http.Header, []byte, error) {

	log.Printf("METHOD : %v \nURL : %v\n", method, url)
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
