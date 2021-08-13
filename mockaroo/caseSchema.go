package mockaroo

import (
	"fmt"
	"log"

	"github.com/troysellers/go-modifier/mockaroo/types"
)

func getSchemaForCase(fields []interface{}, personAccounts bool) []types.IField {
	var mockFields []types.IField
	for _, f := range fields {
		field := f.(map[string]interface{})
		if shouldGetData(field) {

			switch field["name"].(string) {
			case "EntitlementId", "ParentId", "RecordTypeId", "AccountId", "SourceId":
				log.Printf("Skipping %v\n", field["name"])
			case "ContactId":
				mf := types.NewWords(field)
				mf.Max = 0
				mf.Min = 0 // we will populate this one ourselves
				mockFields = append(mockFields, mf)
			case "Subject":
				mf := types.NewCatchPhrase(field)
				mockFields = append(mockFields, mf)
			case "SuppliedName":
				mf := types.NewFullName(field)
				mf.Formula = fmt.Sprintf("this[0,%d]", int(field["length"].(float64)))
				mockFields = append(mockFields, mf)
			case "SuppliedEmail":
				mf := types.NewEmailAddress(field)
				mf.Formula = fmt.Sprintf("this[0,%d]", int(field["length"].(float64)))
				mockFields = append(mockFields, mf)
			case "SuppliedPhone":
				mf := types.NewPhone(field)
				mockFields = append(mockFields, mf)
			case "SuppliedCompany":
				mf := types.NewFakeCompanyName(field)
				mf.Formula = fmt.Sprintf("this[0,%d]", int(field["length"].(float64)))
				mockFields = append(mockFields, mf)
			default:
				mf := getMockTypeForField(field)
				mockFields = append(mockFields, mf)
			}
		}
	}

	return mockFields
}
