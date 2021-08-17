package mockaroo

import (
	"log"

	"github.com/troysellers/go-modifier/mockaroo/types"
)

func getSchemaForEvent(fields []interface{}, personAccounts bool) []types.IField {
	var mockFields []types.IField

	for _, f := range fields {
		field := f.(map[string]interface{})
		if shouldGetData(field) {
			var mf types.IField
			switch field["name"].(string) {
			case "RecordTypeId", "IsPartner", "IsCustomerPortal", "ParentId":
				log.Printf("TODO : Have not implemented %v\n", field["name"].(string))
			case "WhoId", "WhatId":
				words := types.NewWords(field)
				words.Max = 0
				words.Min = 0
				mf = words
			default:
				mf = getMockTypeForField(field)

			}
			if mf != nil {
				mockFields = append(mockFields, mf)
			}
		}
	}
	adjusted := handlePersonAccounts(mockFields, personAccounts)
	return adjusted
}
