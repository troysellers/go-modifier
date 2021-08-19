package mockaroo

import (
	"log"

	"github.com/troysellers/go-modifier/mockaroo/types"
)

func getSchemaForContact(fields []interface{}, personAccounts bool) []types.IField {
	var mockFields []types.IField
	for _, f := range fields {
		field := f.(map[string]interface{})
		if shouldGetData(field) {
			switch field["name"].(string) {
			case "ParentId", "IndividualId", "ReportsToId", "Jigsaw", "CleanStatus":
				log.Printf("Skipping %v - yet to be handled \n", field["name"])
			case "LastName":
				mockFields = append(mockFields, types.NewLastName(field))
			case "FirstName":
				mockFields = append(mockFields, types.NewFirstName(field))
			case "JobTitle":
				mockFields = append(mockFields, types.NewJobTitle(field))
			case "MailingLatitude", "OtherLatitude":
				mockFields = append(mockFields, types.NewLatitude(field))
			case "MailingLongitude", "OtherLongitude":
				mockFields = append(mockFields, types.NewLongitude(field))
			case "MailingStreet", "OtherStreet":
				mockFields = append(mockFields, types.NewStreetAddress(field))
			case "MailingCity", "OtherCity":
				mockFields = append(mockFields, types.NewCity(field))
			case "MailingState", "OtherState":
				mockFields = append(mockFields, types.NewState(field))
			case "MailingCountry", "OtherCountry":
				mockFields = append(mockFields, types.NewCountry(field))
			case "MailingPostalCode", "OtherPostalCode":
				mockFields = append(mockFields, types.NewPostalCode(field))
			default:
				mf := getMockTypeForField(field)
				if mf != nil {
					mockFields = append(mockFields, mf)
				}
			}
		}
	}
	return mockFields
}
