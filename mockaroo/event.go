package mockaroo

import (
	"fmt"
	"log"

	"github.com/troysellers/go-modifier/mockaroo/types"
)

func getSchemaForEvent(fields []interface{}, personAccounts bool) []types.IField {
	var mockFields []types.IField

	for _, f := range fields {
		field := f.(map[string]interface{})
		/* formula sytax can include Ruby code. These formulas are designed to adjust start and end times around the activity date */
		if shouldGetData(field) {
			var mf types.IField
			switch field["name"].(string) {
			case "WhoId", "WhatId":
				w := types.NewWords(field)
				w.Max = 0
				w.Min = 0
				mf = w
			case "Description":
				s := types.NewCatchPhrase(field)
				mf = s
			case "ActivityDate":
				dt := types.NewDatetime(field)
				// sets about 70% of the records in the past
				dt.Formula = "if random(0,10) <= 7 then Date.today - random(0,365) else Date.today + random(0,365) end"
				mf = dt
			case "StartDateTime":
				dt := types.NewDatetime(field)
				// if this isn't an all day event, set a random start time on the same day as the activity date
				dt.Formula = "if field('IsAllDayEvent') == true then '' else (field('ActivityDate') + random(0.1,0.5)).strftime('%Y-%m-%dT%H:%M:%S.%L%z') end"
				//dt.Formula = strings.ReplaceAll(dt.Formula, "%d", "\\%d")
				mf = dt
			case "EndDateTime":
				dt := types.NewDatetime(field)
				// if this isn't an all day event, set a random end time some time after the start time
				dt.Formula = "if field('IsAllDayEvent') == false then (date(field('StartDateTime')) + random(30,240)).strftime('%Y-%m-%dT%H:%M:%S.%L%z') else field('StartDateTime') end"
				mf = dt
			case "Priority", "Status", "Subject", "CallType", "Type", "OwnerId", "ShowAs", "IsAllDayEvent", "Location":
				mf = getMockTypeForField(field)
			default:
				log.Printf("TASK : Ignoring %v\n", field["name"])
				//mf = getMockTypeForField(field)
			}
			if mf != nil {
				if mf.GetField().Formula == "" {
					l := int(field["length"].(float64))
					if l > 0 {
						// formula to ensure we don't generate data longer than we can insert into SF
						mf.SetFormula(fmt.Sprintf("this[0,%d]", l))
					}
				}
				mockFields = append(mockFields, mf)
			}
		}
	}
	return mockFields
}
