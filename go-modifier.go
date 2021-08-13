package main

import (
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
	"github.com/troysellers/go-modifier/sforce"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	var op = flag.String("op", "", "create | update")
	var query = flag.Bool("query", true, "(update) run the query only, do not execute the update in Salesforce")
	var count = flag.Int("count", 10, "(create) how many records to get from mockaroo")
	var obj = flag.String("obj", "", "(create) specify which salesforce object do you want to create")
	var references = flag.Bool("references", true, "(create) set to true if you want to populate reference fields to random data in the Salesforce org. ")
	var fetchOnly = flag.Bool("fetch", false, "(create) When true will fetch and merge mockaroo data but will not send to Salesforce.")
	var whoObj = flag.String("who", "", "(create) If creating activities (tasks/events) you need to specify the who object (user|contact)")
	var whatObj = flag.String("what", "", "(create) If creating activities (tasks/events) you need to specify the what object (any activity enabled obj)")
	var personAccounts = flag.Bool("personaccounts", false, "(create) Set to true if you want to create person accounts (or relate other objects to person accounts).")

	flag.Parse()

	if *op == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	cfg := config.NewConfig()
	c, err := sforce.NewRestClient(&cfg.SF)
	if err != nil {
		panic(err)
	}
	// get a syncMap to store any downloaded Ids so we only do this once.
	var objIds sync.Map

	switch *op {
	case "update":
		var wg sync.WaitGroup
		for _, q := range cfg.SF.Queries {
			wg.Add(1)
			go modify(q, cfg, &wg, *query, c, &objIds)
		}
		wg.Wait()
	case "create":
		log.Printf("Creating for %v\n", *obj)
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
