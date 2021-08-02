package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/troysellers/go-modifier/config"
	"github.com/troysellers/go-modifier/gen"
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
	var op = flag.String("o", "create", "Specify the operation to run [create|update]")
	var query = flag.Bool("q", true, "Run the Salesforce query only to test what would be modified")
	var count = flag.Int("c", 10, "Defines how many records need to be created when using the mockaroo generate")
	var mockSchema = flag.String("s", "", "Defines the schema to download from mockaroo. If flag not used, will get data from the queries in the .env file")
	flag.Parse()

	cfg := config.NewConfig()

	if strings.EqualFold(*op, "update") {
		var wg sync.WaitGroup
		for _, q := range cfg.SF.Queries {
			wg.Add(1)
			go modify(q, &cfg.SF, &wg, *query)
		}
		wg.Wait()
	}

	if strings.EqualFold(*op, "create") {
		if *mockSchema == "" {
			log.Printf("You need to specify the mockaroo schema (-schema flag) name for this to operate\n----")
			printUsage()
		}
		c, err := sforce.NewRestClient(&cfg.SF)
		if err != nil {
			panic(err)
		}
		csvFile, err := downloadFromMockaroo(cfg, *mockSchema, *count)
		if err != nil {
			panic(err)
		}
		if err := sforce.UploadCSVToSalesforce(&cfg.SF, c, csvFile, *mockSchema); err != nil {
			panic(err)
		}
	}
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

func downloadFromMockaroo(cfg *config.Config, s string, r int) (string, error) {

	file, err := gen.GetDataFromMockaroo(&cfg.Mockaroo, s, r)
	if err != nil {
		return "", err
	}
	return file, nil
}

func modify(q string, cfg *config.SFConfig, wg *sync.WaitGroup, queryOnly bool) {

	defer wg.Done()
	c, err := sforce.NewRestClient(cfg)
	if err != nil {
		panic(err)
	}
	log.Printf("Query to run %v : query only %v", q, queryOnly)
	obj, modifiedFile, err := sforce.GetBulkQuery(cfg, c, q)
	if err != nil {
		panic(err)
	}
	if !queryOnly {
		if err := sforce.UploadCSVToSalesforce(cfg, c, modifiedFile, obj); err != nil {
			panic(err)
		}
	}
}
