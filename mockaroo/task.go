package mockaroo

import (
	"log"
	"time"

	"github.com/troysellers/go-modifier/mockaroo/types"
)

func getSchemaForTask(fields []interface{}, personAccounts bool) []types.IField {
	var mockFields []types.IField

	for _, f := range fields {
		field := f.(map[string]interface{})
		if shouldGetData(field) {
			var mf types.IField
			switch field["name"].(string) {
			case "WhoId", "WhatId":
				w := types.NewWords(field)
				w.Max = 0
				w.Min = 0
				mf = w
			case "Description":
				s := types.NewSentences(field)
				s.Max = 100
				s.Min = 1
				mf = s
			case "CompletedDateTime":
				dt := types.NewDatetime(field)
				dt.Max = time.Now().Format("01-02-2006")
			case "Priority", "ActivityDate", "Status", "Subject", "CallType", "Type", "OwnerId":
				mf = getMockTypeForField(field)
			default:
				log.Printf("TASK : Ignoring %v\n", field["name"])
				//mf = getMockTypeForField(field)
			}
			if mf != nil {
				mockFields = append(mockFields, mf)
			}
		}
	}
	return mockFields
}
