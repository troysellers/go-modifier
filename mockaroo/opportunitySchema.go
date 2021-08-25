package mockaroo

import (
	"log"
	"strings"
	"time"

	"github.com/troysellers/go-modifier/mockaroo/types"
)

func getSchemaForOpportunity(fields []interface{}, personAccounts bool) []types.IField {
	var mockFields []types.IField
	for _, f := range fields {
		field := f.(map[string]interface{})
		if shouldGetData(field) {
			var mf types.IField
			switch field["name"].(string) {
			case "RecordTypeId", "Probability", "ForecastCategoryName", "Territory2Id", "IsExcludedFromTerritory2Filter", "SyncedQuoteId":
				log.Printf("TODO : Have not implemented %v\n", field["name"].(string))
			case "Name":
				mf = types.NewConstructionSubContract(field)
			case "Amount":
				num := types.NewNumber(field)
				num.Min = 15000
				num.Max = 500000
				num.Decimals = 2
				mf = num
			case "Description":
				mf = types.NewCatchPhrase(field)
			case "StageName":
				list := types.NewCustomList(field)
				list.Values = getOpenOppStages(field)
				mf = list
			case "CloseDate":
				dt := types.NewDatetime(field)
				dt.Min = time.Now().Format("01/02/2006")
				dt.Max = time.Now().AddDate(1, 6, 0).Format("01/02/2006")
				mf = dt
			default:
				mf = getMockTypeForField(field)
			}
			if mf != nil {
				mockFields = append(mockFields, mf)
			}
		}
	}
	return mockFields
}

func getOpenOppStages(field map[string]interface{}) []string {
	v := make([]string, 0)

	plvs := field["picklistValues"].([]interface{})
	for _, plv := range plvs {
		if plv.(map[string]interface{})["value"] != nil && !strings.Contains(plv.(map[string]interface{})["value"].(string), "Closed") {
			v = append(v, plv.(map[string]interface{})["value"].(string))
		}
	}

	return v
}
