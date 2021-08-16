package mockaroo

import (
	"log"

	"github.com/troysellers/go-modifier/mockaroo/types"
)

func getSchmeaForLead(fields []interface{}, personAccounts bool) []types.IField {
	var mockFields []types.IField
	for _, f := range fields {
		field := f.(map[string]interface{})
		if shouldGetData(field) {
			switch field["name"].(string) {
			case "Jigsaw", "DandbCompanyId", "IndividualId":
				log.Printf("Skipping %v - yet to be handled \n", field["name"])
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
