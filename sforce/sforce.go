package sforce

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/simpleforce/simpleforce"
	"github.com/troysellers/go-modifier/config"
	"github.com/troysellers/go-modifier/file"
	"github.com/troysellers/go-modifier/lorem"
	"github.com/tzmfreedom/go-soapforce"
)

const bulkQueryEndpoint string = "/services/data/v%.1f/jobs/query"
const bulkIngestEndpoint string = "/services/data/v%.1f/jobs/ingest/"

// used for managing the bulk query
type QueryJob struct {
	Create       BulkQueryJobCreate
	SFClient     *simpleforce.Client
	SFObjectMeta *simpleforce.SObjectMeta
	SessionId    string
	SFEndpoint   string
	ApiVersion   float32
	BulkJob      BulkJob
	QueryData    [][]string
	FileName     string
	FilePath     string
}

// creats the BulkQuery
type BulkQueryJobCreate struct {
	Operation string `json:"operation"`
	Query     string `json:"query"`
}

// gets the results and status etc.
type BulkJob struct {
	Id                     string  `json:"id"`
	Operation              string  `json:"operation"`
	Object                 string  `json:"object"`
	CreatedById            string  `json:"createdById"`
	CreatedDate            string  `json:"createdDate"`
	SystemModstamp         string  `json:"systemModstamp"`
	State                  string  `json:"state"`
	ConcurrencyMode        string  `json:"concurrencyMode"`
	ContentType            string  `json:"contentType"`
	ApiVersion             float32 `json:"apiVersion"`
	LineEnding             string  `json:"lineEnding"`
	ColumnDelimiter        string  `json:"columnDelimiter"`
	NumberRecordsProcessed int     `json:"numberOfRecordsProcessed"`
	Retries                int     `json:"retries"`
	TotalProcessingTime    int     `json:"totalProcessingTime"`
}

// used for managing bulk ingest
type UpsertJob struct {
	Create       BulkUpsertJobCreate
	Object       string
	Job          BulkUpsertJob
	Close        BulkUpsertJobClose
	SessionId    string  // the salesforce session id to use
	SFEndpoint   string  // the salesforce endpoint to use
	ApiVersion   float32 // the Salesforce api version
	ModifiedFile string  // file on disk of modified data
	Cfg          *config.Config
}

// creates the bulk ingest job
type BulkUpsertJobCreate struct {
	Object              string `json:"object"`
	ExternalIdFieldName string `json:"externalIdFieldName"`
	ContentType         string `json:"contentType"`
	Operation           string `json:"operation"`
}

// results and
type BulkUpsertJob struct {
	Id                      string  `json:"id"`
	Operation               string  `json:"operation"`
	Object                  string  `json:"object"`
	CreatedBy               string  `json:"createdBy"`
	CreatedDate             string  `json:"createdDate"`
	SystemModstamp          string  `json:"systemModstamp"`
	State                   string  `json:"state"`
	ExternalIdFieldName     string  `json:"externalIdFieldName"`
	ConcurrencyMode         string  `json:"Parallel"`
	ContentType             string  `json:"contentType"`
	ApiVersion              float32 `json:"apiVersion"`
	ContentUrl              string  `json:"contentUrl"` // url to write batches to
	LineEnding              string  `json:"lineEnding"`
	ColumnDelimiter         string  `json:"columnDelimiter"`
	NumberRecordsProcessed  int     `json:"numberRecordsProcessed"`
	NumberRecordsFailed     int     `json:"numberRecordsFailed"`
	Retries                 int     `json:"retries"`
	TotalProcessingTime     int     `json:"totalProcessingTime"`
	ApiActiveProcessingTime int     `json:"apiActiveProcessingTime"`
	ApexProcessingTime      int     `json:"apexProcessingTime"`
}

// closing the upsert
type BulkUpsertJobClose struct {
	State string `json:"state"`
}

// creates an authorised REST client
func NewRestClient(cfg *config.SFConfig) (*simpleforce.Client, error) {

	c := simpleforce.NewClient(fmt.Sprintf("https://%v", cfg.LoginUrl), simpleforce.DefaultClientID, fmt.Sprintf("%.1f", cfg.ApiVersion))
	if c == nil {
		return nil, fmt.Errorf("unable to establish a Salesforce REST Client")
	}
	err := c.LoginPassword(cfg.Username, cfg.Password, cfg.Token)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// creates an authenticated Salesforce session
// using the SOAP client. Don't think it is used but?
func NewSoapClient(cfg *config.SFConfig) (*soapforce.Client, error) {

	sfc := soapforce.NewClient()
	sfc.SetLoginUrl(cfg.LoginUrl)
	sfc.SetApiVersion(fmt.Sprintf("%.1f", cfg.ApiVersion))
	sfc.SetDebug(cfg.SfDebug)

	res, err := sfc.Login(cfg.Username, cfg.Password)
	if err != nil {
		return nil, err
	}
	if res.SessionId != "" {
		log.Printf("Authenticated to the org at %v", res.ServerUrl)
	}

	sfc.SetBatchSize(cfg.SfBatchSize)
	return sfc, nil
}

func getObjectNameFromQuery(q string) string {
	tokens := strings.Split(q, " ")
	for i, t := range tokens {
		if strings.EqualFold(t, "from") {
			return tokens[i+1]
		}
	}
	return ""
}

// fetches metadata for the object that has been queried.
func (qj *QueryJob) fetchMetaForObj() {

	sobj := qj.SFClient.SObject(qj.BulkJob.Object)
	meta := sobj.Describe()

	qj.SFObjectMeta = meta
}

func getField(fname string, fields []interface{}) map[string]interface{} {

	for _, f := range fields {
		field := f.(map[string]interface{})
		if strings.EqualFold(field["name"].(string), fname) {
			return field
		}
	}
	return nil
}

/*
	changes the data in the file on disk..
*/
func (qj *QueryJob) ModifyData(cfg *config.Config, objIds *sync.Map, c *simpleforce.Client) error {

	log.Println("\n\nwe are trying to modify things")

	// for each header (field name)
	for i, fieldName := range qj.QueryData[0] {
		// get the SF metadata for this field
		f := getField(fieldName, (*qj.SFObjectMeta)["fields"].([]interface{}))
		// if field is updateable
		if f["updateable"].(bool) {
			// loop through each row in the file
			for _, row := range qj.QueryData[1:] {

				val, err := GetValueForType(cfg, f, qj.SFClient, objIds)
				if err != nil {
					log.Printf("%v", err)
				} else {
					if cfg.ModifyWithNull {
						row[i] = "null"
					} else {
						// update the column with this random value
						row[i] = fmt.Sprintf("%v", val)
					}
					log.Printf("update %v to %v", fieldName, row[i])
				}
			}
		}
	}
	return nil
}

// returns the filepath of downloaded CSV of results.
// this is blocking
// returns the Object the query was for, the filepath of the downloaded results.. and an error
func GetBulkQuery(cfg *config.Config, c *simpleforce.Client, q string) (QueryJob, error) {

	// query the data
	queryJob := &QueryJob{
		Create: BulkQueryJobCreate{
			Operation: "query",
			Query:     q,
		},
		SessionId:  c.GetSid(),
		SFEndpoint: c.GetLoc(),
		ApiVersion: cfg.SF.ApiVersion,
		SFClient:   c,
	}
	if err := queryJob.createQueryJob(); err != nil {
		return QueryJob{}, err
	}
	// wait for query to finish
	for {
		log.Printf("Job state %v with %v records", queryJob.BulkJob.State, queryJob.BulkJob.NumberRecordsProcessed)
		if queryJob.BulkJob.State == "JobComplete" {
			break
		}
		time.Sleep(10 * time.Second)
		if err := queryJob.getBulkJobState(); err != nil {
			return *queryJob, err
		}
	}
	// get the results of that data
	if err := queryJob.getResults(); err != nil {
		return *queryJob, err
	}
	var err error
	queryJob.FileName = fmt.Sprintf("%v-query.csv", queryJob.BulkJob.Object)
	queryJob.FilePath, err = file.BuildFilePath(queryJob.FileName, cfg)
	if err != nil {
		return *queryJob, err
	}
	// write the data to a temporary file
	d, err := file.WriteCsv(queryJob.FilePath, queryJob.QueryData)
	if err != nil {
		return *queryJob, err
	}
	queryJob.FilePath = d

	return *queryJob, nil
}

func UploadCSVToSalesforce(cfg *config.Config, c *simpleforce.Client, csvfile string, obj string) error {

	// update the object if we are loading personaccounts
	if strings.EqualFold(obj, "personaccount") {
		obj = "account"
	}

	// create the bulk update job
	uj := UpsertJob{
		SessionId:    c.GetSid(),
		SFEndpoint:   c.GetLoc(),
		ApiVersion:   cfg.SF.ApiVersion,
		ModifiedFile: csvfile,
		Cfg:          cfg,
		Create: BulkUpsertJobCreate{
			Object:              obj,
			ExternalIdFieldName: "Id",
			ContentType:         "CSV",
			Operation:           "upsert",
		},
	}

	if err := uj.createBulkIngest(); err != nil {
		return err
	}
	if err := uj.sendData(); err != nil {
		return err
	}

	if err := uj.closeJob(); err != nil {
		return err
	}
	if err := uj.GetJobStatus(); err != nil {
		return err
	}
	return nil
}

// creates the BulkV2 Query
func (qj *QueryJob) createQueryJob() error {

	h := make(map[string]string)
	h["Content-Type"] = "application/json; charset=UTF-8"
	h["Accept"] = "application/json"
	h["Authorization"] = fmt.Sprintf("Bearer %v", qj.SessionId)
	url := fmt.Sprintf("%v%v", qj.SFEndpoint, fmt.Sprintf(bulkQueryEndpoint, qj.ApiVersion))

	b, err := json.Marshal(qj.Create)
	if err != nil {
		return err
	}

	_, r, err := doHttp(url, qj.SessionId, b, "POST", h)
	if err != nil {
		return err
	}
	qj.BulkJob = BulkJob{}
	if err := json.Unmarshal(r, &qj.BulkJob); err != nil {
		return err
	}
	// fetch the object metadata
	qj.fetchMetaForObj()
	return nil
}

// creates the initial BulkIngest Bulk V2 job
func (uj *UpsertJob) createBulkIngest() error {
	h := make(map[string]string)
	h["Content-Type"] = "application/json; charset=UTF-8"
	h["Accept"] = "application/json"
	h["Authorization"] = fmt.Sprintf("Bearer %v", uj.SessionId)

	url := fmt.Sprintf("%v%v", uj.SFEndpoint, fmt.Sprintf(bulkIngestEndpoint, uj.ApiVersion))
	b, err := json.Marshal(uj.Create)
	if err != nil {
		return err
	}
	_, responseBytes, err := doHttp(url, uj.SessionId, b, "POST", h)
	if err != nil {
		return err
	}
	var upsertJob BulkUpsertJob
	if err := json.Unmarshal(responseBytes, &upsertJob); err != nil {
		panic(err)
	}
	uj.Job = upsertJob
	return nil
}

// loads the CSV into the BulkV2 upsert job
func (uj *UpsertJob) sendData() error {
	h := make(map[string]string)
	h["Authorization"] = fmt.Sprintf("Bearer %v", uj.SessionId)
	h["Content-Type"] = "text/csv"
	h["Accept"] = "application/json"

	url := fmt.Sprintf("%v/%v", uj.SFEndpoint, uj.Job.ContentUrl)
	body, err := getCSVBytes(uj.ModifiedFile)
	if err != nil {
		return err
	}
	fmt.Printf("-- WE HAVE %d encoded bytes to send ---\n", len(body))
	_, responseBytes, err := doHttp(url, uj.SessionId, body, "PUT", h)
	if err != nil {
		return err
	}
	fmt.Println(string(responseBytes))
	return nil
}

// Gets the status of the BulkV2 query
func (qj *QueryJob) getBulkJobState() error {

	h := make(map[string]string)
	h["Content-Type"] = "application/json; charset=UTF-8"
	h["Accept"] = "application/json"
	h["Authorization"] = fmt.Sprintf("Bearer %v", qj.SessionId)
	url := fmt.Sprintf("%v%v/%v/", qj.SFEndpoint, fmt.Sprintf(bulkQueryEndpoint, qj.ApiVersion), qj.BulkJob.Id)

	_, rb, err := doHttp(url, qj.SessionId, nil, "GET", h)
	if err != nil {
		return err
	}
	qj.BulkJob = BulkJob{}
	if err := json.Unmarshal(rb, &qj.BulkJob); err != nil {
		return err
	}
	return nil
}

// downloads the Bulk V2 Query results
func (qj *QueryJob) getResults() error {

	h := make(map[string]string)
	h["Content-Type"] = "application/json; charset=UTF-8"
	h["Accept"] = "application/json"
	h["Authorization"] = fmt.Sprintf("Bearer %v", qj.SessionId)
	url := fmt.Sprintf("%v%v/%v/results", qj.SFEndpoint, fmt.Sprintf(bulkQueryEndpoint, qj.BulkJob.ApiVersion), qj.BulkJob.Id)
	resHeaders, resBytes, err := doHttp(url, qj.SessionId, nil, "GET", h)
	for {
		if err != nil {
			return err
		}

		lines, err := csv.NewReader(bytes.NewBuffer(resBytes)).ReadAll()
		if err != nil {
			return err
		}
		//numRecords, _ := strconv.Atoi(resHeaders.Get("Sforce-NumberOfRecords"))
		//limits := resHeaders.Header.Get("Sforce-Limit-Info")
		locator := resHeaders.Get("Sforce-Locator")

		qj.QueryData = append(qj.QueryData, lines...)
		if locator == "null" {
			break
		}
		resHeaders, resBytes, err = doHttp(fmt.Sprintf("%v?locator=%v", url, locator), qj.SessionId, nil, "GET", h)
	}
	return nil
}

// reads the entire file into a byte []
func getCSVBytes(file string) ([]byte, error) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return f, nil
}
func (uj *UpsertJob) GetJobStatus() error {
	h := make(map[string]string)
	h["Content-Type"] = "application/json; charset=UTF-8"
	h["Accept"] = "application/json"
	h["Authorization"] = fmt.Sprintf("Bearer %v", uj.SessionId)

	for {
		url := fmt.Sprintf("%v%v%v/", uj.SFEndpoint, fmt.Sprintf(bulkIngestEndpoint, uj.ApiVersion), uj.Job.Id)
		_, responseBytes, err := doHttp(url, uj.SessionId, nil, "GET", h)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(responseBytes, &uj.Job); err != nil {
			return err
		}
		if uj.Job.State == "JobComplete" {
			break
		}
		log.Println("Waiting 30 seconds before trying again")
		time.Sleep(30 * time.Second) // wait 30 secionds
	}
	if uj.Job.NumberRecordsFailed > 0 {
		uj.fetchFailedRecords()
	}
	log.Printf("Job %s [%s on %s] complete in %dms\n", uj.Job.Id, uj.Job.Operation, uj.Job.Object, uj.Job.TotalProcessingTime)
	log.Printf("Total Records %d\n", uj.Job.NumberRecordsProcessed)
	log.Printf("Records Failed %d\n", uj.Job.NumberRecordsFailed)
	if uj.Job.NumberRecordsFailed > 0 {
		if err := uj.fetchFailedRecords(); err != nil {
			return err
		}
	}
	return nil
}

func (uj *UpsertJob) fetchFailedRecords() error {
	h := make(map[string]string)
	h["Content-Type"] = "application/json; charset=UTF-8"
	h["Accept"] = "application/json"
	h["Authorization"] = fmt.Sprintf("Bearer %v", uj.SessionId)

	url := fmt.Sprintf("%v%v%v/failedResults/", uj.SFEndpoint, fmt.Sprintf(bulkIngestEndpoint, uj.ApiVersion), uj.Job.Id)
	_, responseBytes, err := doHttp(url, uj.SessionId, nil, "GET", h)
	if err != nil {
		return err
	}
	lines, err := csv.NewReader(bytes.NewBuffer(responseBytes)).ReadAll()
	if err != nil {
		return err
	}
	fPath, err := file.BuildFilePath("unsuccessful.csv", uj.Cfg)
	if err != nil {
		return err
	}
	if _, err := file.WriteCsv(fPath, lines); err != nil {
		return err
	}
	return nil
}
func (uj *UpsertJob) closeJob() error {
	h := make(map[string]string)
	h["Content-Type"] = "application/json; charset=UTF-8"
	h["Accept"] = "application/json"
	h["Authorization"] = fmt.Sprintf("Bearer %v", uj.SessionId)

	uj.Close = BulkUpsertJobClose{
		State: "UploadComplete",
	}
	b, err := json.Marshal(uj.Close)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%v%v%v/", uj.SFEndpoint, fmt.Sprintf(bulkIngestEndpoint, uj.ApiVersion), uj.Job.Id)

	_, responseBytes, err := doHttp(url, uj.SessionId, b, "PATCH", h)
	if err != nil {
		return err
	}
	fmt.Println(string(responseBytes))
	return nil
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
		return res.Header, nil, fmt.Errorf("unsuccesful attempt to create Bulk Job %v\n%v", res.Status, string(bytes))
	}

	log.Printf("Received %d bytes with http response %v\n", len(bytes), res.Status)
	return res.Header, bytes, nil
}

func GetValueForType(cfg *config.Config, f map[string]interface{}, c *simpleforce.Client, objIds *sync.Map) (interface{}, error) {

	/* if can be empty, retun empty on a 10%
	if f["nillable"].(bool) && rand.Intn(10) < 2 {
		return nil, nil
	}*/
	switch f["type"].(string) {
	case "id":
		return nil, fmt.Errorf("id values are not supported for generation")
	case "boolean":
		return rand.Intn(10) >= 5, nil
	case "string", "encryptedstring":
		if f["unique"].(bool) {
			return uuid.New(), nil
		}
		l := int(f["length"].(float64))
		return lorem.Word(1, rand.Intn(l)), nil
	case "datetime", "date":
		d := time.Now()
		d = d.AddDate(0, rand.Intn(12), rand.Intn(30))
		return d, nil
	case "reference":
		//	log.Printf("REFERENCCE %v\n", f)
		// get the name of the object this field references
		rt := f["referenceTo"].([]interface{})
		referenceTo := rt[0].(string)
		log.Printf("reference to [%v]\n", referenceTo)
		// have we already got the complete list of ids?
		_, ok := objIds.Load(referenceTo)
		if !ok {
			// if not, get and cache in the sync.Map
			ids, err := GetAllObjIds(cfg, referenceTo, c)
			if err != nil {
				return nil, err
			}
			objIds.Store(referenceTo, ids)
		}
		i, _ := objIds.Load(referenceTo)
		ids := i.([]string)
		return ids[rand.Intn(len(ids))], nil
	case "currency", "double":
		p := f["precision"].(float64)
		s := f["scale"].(float64)
		return rand.Intn(int(p)) / int(math.Pow10(int(s))), nil
	case "email":
		return lorem.Email(), nil
	case "location":
		return nil, fmt.Errorf("location value not implemented yet")
	case "percent":
		return float32(rand.Intn(100)), nil
	case "int":
		d := int(f["digits"].(float64))
		return rand.Intn(int(math.Pow10(d))), nil
	case "phone":
		return "000000000", nil
		//return nil, fmt.Errorf("phone value not implemented yet")
	case "picklist", "multipicklist":
		plv := f["picklistValues"].([]interface{})
		selected := plv[rand.Intn(len(plv))]
		return selected.(map[string]interface{})["value"], nil
	case "textarea":
		/*
			l := int(f["length"].(float64))
			s := lorem.Sentence(1, l)
			if len(s) < int(l) {
				return s, nil
			} else {
				return s[:l], nil
			} */
		return "broken", nil

	case "time":
		return nil, fmt.Errorf("phone value not implemented yet")
	case "url":
		return lorem.Url(), nil
	}

	return nil, nil
}

// returns object, allIds and an error
func GetAllObjIds(cfg *config.Config, obj string, c *simpleforce.Client) ([]string, error) {

	q := fmt.Sprintf("select id from %v", obj)

	if strings.EqualFold(obj, "user") {
		q += " where isActive = true and userType = 'standard'"
	}
	log.Printf("Downloading all IDS [%v]. This could take a while... ", q)
	qj, err := GetBulkQuery(cfg, c, q)
	if err != nil {
		return nil, err
	}

	b, err := file.GetCSVBytes(qj.FilePath)
	if err != nil {
		return nil, err
	}
	rows, err := csv.NewReader(bytes.NewReader(b)).ReadAll()
	if err != nil {
		return nil, err
	}

	var results []string
	for _, r := range rows[1:] { // ignore the header row
		if len(r) > 1 {
			return nil, fmt.Errorf("there has been an error downloading IDs. More than one value per record returned")
		}
		results = append(results, r[0])
	}
	log.Printf("Found %d id values for %v", len(results), obj)
	return results, nil
}
