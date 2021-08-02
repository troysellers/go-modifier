package account

import (
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/tzmfreedom/go-soapforce"
)

func UpdateExternalId(c *soapforce.Client) error {

	qr, err := c.QueryAll("select id, external_id__c from account where External_Id__c = ''")
	var accounts []*soapforce.SObject

	var wg sync.WaitGroup

	for {
		if err != nil {
			return err
		}
		for _, r := range qr.Records {
			sobj := soapforce.SObject{
				Id:     r.Id,
				Type:   r.Type,
				Fields: make(map[string]interface{}),
			}
			sobj.Fields["External_Id__c"] = uuid.New().String()
			accounts = append(accounts, &sobj)
			if len(accounts)%int(c.BatchSize) == 0 {
				wg.Add(1)
				go updateAccounts(c, accounts, &wg)
				accounts = nil
			}
		}
		if qr.Done {
			break
		}
		log.Printf("Fetching the next set of records")
		qr, err = c.QueryMore(qr.QueryLocator)
	}
	if len(accounts) > 0 {
		updateAccounts(c, accounts, &wg)
	}
	wg.Wait()
	return nil
}

func updateAccounts(c *soapforce.Client, accs []*soapforce.SObject, wg *sync.WaitGroup) {

	defer wg.Done()
	log.Printf("Updating %v accounts", len(accs))
	sr, err := c.Update(accs)
	if err != nil {
		log.Printf("%v", err)
		panic(err)
	}
	for _, r := range sr {
		if !r.Success {
			for _, e := range r.Errors {
				log.Printf("%v", e.Message)
			}
		}
	}

}
