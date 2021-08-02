package sforce

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/simpleforce/simpleforce"
	"github.com/troysellers/go-modifier/config"
	"github.com/troysellers/go-modifier/file"
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
	DownloadFile string
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
	Id                  string  `json:"id"`
	Operation           string  `json:"operation"`
	Object              string  `json:"object"`
	CreatedBy           string  `json:"createdBy"`
	CreatedDate         string  `json:"createdDate"`
	SystemModstamp      string  `json:"systemModstamp"`
	State               string  `json:"state"`
	ExternalIdFieldName string  `json:"externalIdFieldName"`
	ConcurrencyMode     string  `json:"Parallel"`
	ContentType         string  `json:"contentType"`
	ApiVersion          float32 `json:"apiVersion"`
	ContentUrl          string  `json:"contentUrl"` // url to write batches to
	LineEnding          string  `json:"lineEnding"`
	ColumnDelimiter     string  `json:"columnDelimiter"`
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

// fetches metadata for the object that has been queried.
func (qj *QueryJob) fetchMetaForObj() {

	sobj := qj.SFClient.SObject(qj.BulkJob.Object)
	meta := sobj.Describe()

	qj.SFObjectMeta = meta
}

/*
	changes the data in the file on disk..
	TODO: random according to query
*/
func (qj *QueryJob) ModifyData() error {
	switch qj.BulkJob.Object {
	case "Account":
		var nameIndex int
		for i, header := range qj.QueryData[0] {
			if strings.EqualFold(header, "Name") {
				nameIndex = i
				break
			}
		}
		for h, _ := range qj.QueryData[1:] {
			qj.QueryData[h+1][nameIndex] = "broken-again"
		}
	case "Contact":
		var nameIndex int
		for i, header := range qj.QueryData[0] {
			if strings.EqualFold(header, "FirstName") {
				nameIndex = i
				break
			}
		}
		for h, _ := range qj.QueryData[1:] {
			qj.QueryData[h+1][nameIndex] = "Broken"
		}
	default:
		return fmt.Errorf("objec type %v not supported for update", qj.QueryData)
	}
	return nil
}

// returns the filepath of downloaded CSV of results.
// this is blocking
// returns the Object the query was for, the filepath of the downloaded results.. and an error
func GetBulkQuery(cfg *config.SFConfig, c *simpleforce.Client, q string) (string, string, error) {

	// query the data
	queryJob := &QueryJob{
		Create: BulkQueryJobCreate{
			Operation: "query",
			Query:     q,
		},
		SessionId:  c.GetSid(),
		SFEndpoint: c.GetLoc(),
		ApiVersion: cfg.ApiVersion,
		SFClient:   c,
	}
	if err := queryJob.createQueryJob(); err != nil {
		return "", "", err
	}
	// wait for query to finish
	for {
		log.Printf("Job state %v with %v records", queryJob.BulkJob.State, queryJob.BulkJob.NumberRecordsProcessed)
		if queryJob.BulkJob.State == "JobComplete" {
			break
		}
		time.Sleep(10 * time.Second)
		if err := queryJob.getBulkJobState(); err != nil {
			return "", "", err
		}
	}
	// get the results of that data
	if err := queryJob.getResults(); err != nil {
		return "", "", err
	}
	// write the data to a temporary file
	d, err := file.WriteCsv(fmt.Sprintf("%v-query.csv", queryJob.BulkJob.Object), queryJob.QueryData)
	if err != nil {
		return "", "", err
	}
	queryJob.DownloadFile = d
	// change the fields in the data
	if err := queryJob.ModifyData(); err != nil {
		return "", "", err
	}
	// write the CSV back to file
	d2, err := file.WriteCsv(fmt.Sprintf("%v-query-modified.csv", queryJob.BulkJob.Object), queryJob.QueryData)
	if err != nil {
		panic(err)
	}
	return queryJob.BulkJob.Object, d2, nil
}

func UploadCSVToSalesforce(cfg *config.SFConfig, c *simpleforce.Client, csvfile string, obj string) error {
	// create the bulk update job
	uj := UpsertJob{
		SessionId:    c.GetSid(),
		SFEndpoint:   c.GetLoc(),
		ApiVersion:   cfg.ApiVersion,
		ModifiedFile: csvfile,
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
