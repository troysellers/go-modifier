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

	if len(os.Args) < 2 {
		printUsage()
	}
	var op = flag.String("op", "", "Specify the operation to run [create|update]")
	var query = flag.Bool("query", true, "Run the Salesforce query only to test what would be modified")
	var count = flag.Int("count", 10, "Defines how many records need to be created when using the mockaroo generate")
	var obj = flag.String("obj", "", "Which salesforce object do you want to fetch data for")
	var references = flag.Bool("references", false, "set to true if you want to query for records to relate this to. ")
	//var mockSchema = flag.String("s", "", "Defines the schema to download from mockaroo. If flag not used, will get data from the queries in the .env file")
	var fetchOnly = flag.Bool("fetch", false, "When true will fetch and merge mockaroo data but will not send to Salesforce.")
	var whoObj = flag.String("who", "", "If creating activities (tasks/events) you need to specify the who object (user|contact)")
	var whatObj = flag.String("what", "", "If creating activities (tasks/events) you need to specify the what object (any activity enabled obj)")

	flag.Parse()

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
		/*
			if *mockSchema == "" {
				log.Printf("You need to specify the mockaroo schema (-schema flag) name for this to operate\n----")
				printUsage()
			}
		*/
		//csvFile, err := downloadFromMockaroo(cfg, *mockSchema, *count)
		o := c.SObject(*obj)
		m := o.Describe()
		csvFile, err := mockaroo.GetDataForObj(m, cfg, *count)
		if err != nil {
			panic(err)
		}
		/*switch *mockSchema {
		case "case", "opportunity", "contact":
			if err := updateIds(cfg, csvFile, "account", "accountId", &objIds, c); err != nil {
				panic(err)
			}
			fmt.Printf("Updated to point to random, existing Account IDs")

		case "task", "event":
			if *whatObj == "" {
				log.Printf("\n**This will create %d %v records without setting the WhatID field Is this OK?\n", *count, *mockSchema)
				time.Sleep(time.Second * 5)
			} else {
				if err := updateIds(cfg, csvFile, *whatObj, "WhatId", &objIds, c); err != nil {
					panic(err)
				}
			}
			if *whoObj == "" {
				log.Printf("\n**This will create %d %v records without setting the WhoID field. Is this OK?\n", *count, *mockSchema)
				time.Sleep(time.Second * 5)
			} else {
				if err := updateIds(cfg, csvFile, *whoObj, "WhoId", &objIds, c); err != nil {
					panic(err)
				}
			}
		} */

		if *references {
			fields := (*m)["fields"].([]interface{})
			for _, f := range fields {
				// TODO : handle polymorphic keys better than this...
				field := f.(map[string]interface{})
				if strings.EqualFold(*obj, "task") || strings.EqualFold(*obj, "event") {
					if field["relationshipName"] == "Who" {
						field["referenceTo"] = []string{*whoObj}
					}
					if field["relationshipName"] == "What" {
						field["referenceTo"] = []string{*whatObj}
					}
				}
				if field["referenceTo"].([]interface{}) != nil {
					// fetch all the possible Ids for this.
					_, err := sforce.GetValueForType(cfg, field, c, &objIds)
					if err != nil {
						panic(err)
					}
					fieldName := field["name"].(string)
					rt := field["referenceTo"].([]interface{})
					referenceTo := rt[0].(string)
					if err := updateIds(cfg, csvFile, referenceTo, fieldName, &objIds, c); err != nil {
						panic(err)
					}
				}
			}
		}

		// always update the owner
		if err := updateIds(cfg, csvFile, "user", "ownerId", &objIds, c); err != nil {
			panic(err)
		}
		if !*fetchOnly {
			// write data into Salesforce
			if err := sforce.UploadCSVToSalesforce(cfg, c, csvFile, *obj); err != nil {
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

/*
func downloadFromMockaroo(cfg *config.Config, s string, r int) (string, error) {

	log.Printf("downloading %v", s)
	file, err := gen.GetDataFromMockaroo(&cfg.Mockaroo, s, r)
	if err != nil {
		return "", err
	}
	return file, nil
}
*/
func modify(q string, cfg *config.Config, wg *sync.WaitGroup, queryOnly bool, c *simpleforce.Client, objIds *sync.Map) {

	defer wg.Done()
	log.Printf("Query to run %v : query only %v", q, queryOnly)
	queryJob, err := sforce.GetBulkQuery(cfg, c, q)
	if err != nil {
		panic(err)
	}

	if !queryOnly {
		// change the fields in the data
		// depending on the query, this can take some time if it is populating referenced fields randomly.
		if err := queryJob.ModifyData(cfg, objIds, c); err != nil {
			panic(err)
		}
		// write the CSV back to file
		fPath, err := file.BuildFilePath(fmt.Sprintf("%v-query-modified.csv", queryJob.BulkJob.Object))
		if err != nil {
			panic(err)
		}
		d2, err := file.WriteCsv(fPath, queryJob.QueryData)
		if err != nil {
			panic(err)
		}
		if err := sforce.UploadCSVToSalesforce(cfg, c, d2, queryJob.BulkJob.Object); err != nil {
			panic(err)
		}
	}
}
