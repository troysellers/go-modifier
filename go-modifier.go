package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/simpleforce/simpleforce"
	"github.com/troysellers/go-modifier/config"
	"github.com/troysellers/go-modifier/file"
	"github.com/troysellers/go-modifier/mockaroo"
	"github.com/troysellers/go-modifier/mockaroo/types"
	"github.com/troysellers/go-modifier/sforce"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	if len(os.Args) < 2 {
		printUsage()
	}
	var op = flag.String("op", "", "Specify the operation to run [create|update]")
	var query = flag.Bool("query", true, "Update Flag : Run the Salesforce query only to test what would be modified")
	var count = flag.Int("count", 10, "Create Flag : Defines how many records need to be created when using the mockaroo generate")
	var obj = flag.String("obj", "", "Create Flag : Which salesforce object do you want to fetch data for")
	var references = flag.Bool("references", false, "Create Flag : set to true if you want to query for records to relate this to. ")
	var fetchOnly = flag.Bool("fetch", false, "Create Flag : When true will fetch and merge mockaroo data but will not send to Salesforce.")
	var whoObj = flag.String("who", "", "Create Flag : If creating activities (tasks/events) you need to specify the who object (user|contact)")
	var whatObj = flag.String("what", "", "Create Flag : If creating activities (tasks/events) you need to specify the what object (any activity enabled obj)")
	var personAccounts = flag.Bool("personaccounts", false, "Create Flag : Set to true if you want to create person accounts. If this is true it will also set any lookups that reference accounts to reference person accounts (i.e. you might want to create opportunities for person accounts)")

	flag.Parse()

	cfg := config.NewConfig()
	c, err := sforce.NewRestClient(&cfg.SF)
	if err != nil {
		panic(err)
	}
	// get a syncMap to store any downloaded Ids so we only do this once.
	var objIds sync.Map

	switch *op {
	case "upload":
		uploadToSF("/tmp/mockaroo-data/lead-update.csv", cfg, c)
	case "getstuff":
		getstuff()
	case "update":
		var wg sync.WaitGroup
		for _, q := range cfg.SF.Queries {
			wg.Add(1)
			go modify(q, cfg, &wg, *query, c, &objIds)
		}
		wg.Wait()
	case "create":
		o := c.SObject(*obj)
		mr := &mockaroo.MockarooRequest{
			SObject:        o.Describe(),
			Cfg:            cfg,
			Count:          *count,
			PersonAccounts: *personAccounts,
		}
		if err := mr.GetDataForObj(); err != nil {
			panic(err)
		}
		// TODO : check that we haven't excluded reference ids in the schema creations
		if *references {
			fields := mr.Schema
			for _, f := range fields {
				log.Printf("MockField %v\n", f.GetField().Name)
				// TODO : handle polymorphic keys better than this...
				field := f.GetField().SforceMeta

				if strings.EqualFold(*obj, "task") || strings.EqualFold(*obj, "event") {
					log.Printf("\v%v\v", "handling tasks and events!")
					if field["relationshipName"] == "Who" {
						field["referenceTo"] = []string{*whoObj}
					}
					if field["relationshipName"] == "What" {
						field["referenceTo"] = []string{*whatObj}
					}
				}
				// look for the relationship fields that have been included in the schema
				if field["relationshipName"] != nil {
					// fetch all the possible Ids for this.
					fieldName := field["name"].(string)
					rt := field["referenceTo"].([]interface{})
					var referenceTo string
					if fieldName == "OwnerId" {
						referenceTo = "User"
					} else {
						referenceTo = rt[0].(string)
					}

					_, ok := objIds.Load(referenceTo)
					if !ok {
						// if not, get and cache in the sync.Map
						ids, err := sforce.GetAllObjIds(cfg, referenceTo, c)
						if err != nil {
							panic(err)
						}
						objIds.Store(referenceTo, ids)
					}
					if err := updateIds(cfg, mr.FilePath, referenceTo, fieldName, &objIds, c); err != nil {
						panic(err)
					}
				}
			}
		} else if err := updateIds(cfg, mr.FilePath, "user", "ownerId", &objIds, c); err != nil { // always update the owner
			panic(err)
		}
		if !*fetchOnly {
			// write data into Salesforce
			if err := sforce.UploadCSVToSalesforce(cfg, c, mr.FilePath, *obj); err != nil {
				panic(err)
			}
		}
	}
}

//
// f - filename that contains the CSV to modify
// obj - the object name of the referenced field (e.g. if you want to update the ownerId col, this should be user)
// col - the column in the CSV that references this object (e.g. it would be "ownerId" if you wanted to update record owners)
// objIds - the syncMap that contains all the previously downloaded sets of Ids - trying to save some time.
// c - the salesforce REST Client in case we need to get some more id values
//
// function updates a column with random values selected from the complete set of possibles out of salesforce.
// this can take some time, we execute bulk queries in case you want to randomly select from 1 million accounts (as an example)
func updateIds(cfg *config.Config, f string, obj string, col string, objIds *sync.Map, c *simpleforce.Client) error {
	var ids []string
	var err error
	i, ok := objIds.Load(obj) // test if we have the list of IDs in our cache
	if !ok {
		ids, err = sforce.GetAllObjIds(cfg, obj, c) // download them if we don't
		if err != nil {
			return err
		}
		objIds.Store(obj, ids) // store in our cache
	} else {
		ids = i.([]string)
	}
	if err := file.UpdateColumn(f, col, ids); err != nil {
		return err
	}
	return nil
}

func printUsage() {
	log.Println("The modifier has two functions")
	log.Println("\t1 '-o update' use queries in .env file to modify existing salesforce data")
	log.Println("\t2 '-o create' use mockaroo to add new data in Salesforce.")
	log.Println("Create new data uses Mockaroo as the data source - https://www.mockaroo.com/projects/25058")
	log.Println("Modifying data will query Salesforce and then update the fields in this query with lipsum generated random stuff")
	log.Println("Usage:")
	log.Println("\tgo-modifier -o update -query=false")
	log.Println("\t\twill update using the Salesforce queries in the .env file")
	log.Println("\tgo-modifier -o create -c 100 -s contact")
	log.Println("\t\twill download 100 contact records from mockaroo and load them into salesforce.")
	os.Exit(1)
}
func modify(q string, cfg *config.Config, wg *sync.WaitGroup, queryOnly bool, c *simpleforce.Client, objIds *sync.Map) {

	defer wg.Done()
	log.Printf("Query to run %v : query only %v", q, queryOnly)
	queryJob, err := sforce.GetBulkQuery(cfg, c, q)
	if err != nil {
		panic(err)
	}
	// change the fields in the data
	// depending on the query, this can take some time if it is populating referenced fields randomly.
	if err := queryJob.ModifyData(cfg, objIds, c); err != nil {
		panic(err)
	}
	// write the CSV back to file
	fPath, err := file.BuildFilePath(fmt.Sprintf("%v-query-modified.csv", queryJob.BulkJob.Object), cfg)
	if err != nil {
		panic(err)
	}
	d2, err := file.WriteCsv(fPath, queryJob.QueryData)
	if err != nil {
		panic(err)
	}
	if !queryOnly {

		if err := sforce.UploadCSVToSalesforce(cfg, c, d2, queryJob.BulkJob.Object); err != nil {
			panic(err)
		}
	}
}

func uploadToSF(f string, cfg *config.Config, c *simpleforce.Client) {
	if err := sforce.UploadCSVToSalesforce(cfg, c, f, "Lead"); err != nil {
		panic(err)
	}
}

func getstuff() {
	schema := []types.IField{}
	firstName := make(map[string]interface{})
	firstName["name"] = "FirstName"
	schema = append(schema, types.NewFirstName(firstName))

	lastName := make(map[string]interface{})
	lastName["name"] = "LastName"
	schema = append(schema, types.NewLastName(lastName))

	company := make(map[string]interface{})
	company["name"] = "Company"
	schema = append(schema, types.NewFakeCompanyName(company))

	b, err := json.Marshal(schema)
	if err != nil {
		panic(err)
	}
	url := fmt.Sprintf("https://api.mockaroo.com/api/generate.csv?key=%v&count=%d&include_header=false", "c04c9a30", 5000)
	filePath := "/tmp/mockaroo-data/names.csv"

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err = f.Write([]byte("FirstName,LastName, Company\n")); err != nil {
		panic(err)
	}
	i := 0
	for {

		_, respbytes, err := mockaroo.DoHttp(url, "", b, "POST", nil)
		if err != nil {
			panic(err)
		}
		if _, err = f.Write(respbytes); err != nil {
			panic(err)
		}
		i++
		fmt.Printf("fetched %d\n", 5000*i)
		if i == 200 {
			break
		}
	}
	log.Printf("written to %v\n", filePath)
}
