package mockaroo

import (
	"fmt"
	"log"

	"github.com/troysellers/go-modifier/mockaroo/types"
)

func getSchemaForContact(fields []interface{}, personAccounts bool) []types.IField {
	var mockFields []types.IField
	for _, f := range fields {
		field := f.(map[string]interface{})
		if shouldGetData(field) {
			var mf types.IField
			switch field["name"].(string) {
			case "ParentId", "IndividualId", "ReportsToId":
				log.Printf("Skipping %v - yet to be handled \n", field)
			case "LastName":
				mf = types.NewLastName(field)
			case "FirstName":
				mf = types.NewFirstName(field)
			case "JobTitle":
				mf = types.NewJobTitle(field)
			case "MailingLatitude", "OtherLatitude":
				mf = types.NewLatitude(field)
			case "MailingLongitude", "OtherLongitude":
				mf = types.NewLongitude(field)
			case "MailingStreet", "OtherStreet":
				mf = types.NewStreetAddress(field)
			case "MailingCity", "OtherCity":
				mf = types.NewCity(field)
			case "MailingState", "OtherState":
				mf = types.NewState(field)
			case "MailingCountry", "OtherCountry":
				mf = types.NewCountry(field)
			case "MailingPostalCode", "OtherPostalCode":
				mf = types.NewLongitude(field)
			default:
				mf = getMockTypeForField(field)
			}
			if mf != nil {
				// if there is a length set the formula to ensure Mockaroo will truncate the
				// generated value. This is to ensure we don't get errors on (mainly) text fields
				// when inserting into salesforce.
				l := int(field["length"].(float64))
				if l > 0 {
					if l > 1000 { //lets not go crazy with text
						l = 1000
					}
					mf.SetFormula(fmt.Sprintf("this[0,%d]", l))
				}
				mockFields = append(mockFields, mf)
			}
		}
	}
	return mockFields
}
