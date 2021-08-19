package mockaroo

import (
	"log"
	"strings"

	"github.com/troysellers/go-modifier/mockaroo/types"
)

func getSchemaForAccount(fields []interface{}, personAccounts bool) []types.IField {
	var mockFields []types.IField

	for _, f := range fields {
		field := f.(map[string]interface{})
		if shouldGetData(field) {
			var mf types.IField
			switch field["name"].(string) {
			case "RecordTypeId", "IsPartner", "IsCustomerPortal", "ParentId", "CleanStatus", "Jigsaw":
				log.Printf("TODO : Have not implemented %v\n", field["name"].(string))
			case "Name":
				mf = types.NewFakeCompanyName(field)
			case "DunsNumber":
				mf = types.NewDUNSNumber(field)
			case "TickerSymbol":
				mf = types.NewStockSymbol(field)
			case "BillingStreet", "ShippingStreet", "PersonMailingStreet", "PersonOtherStreet":
				mf = types.NewStreetAddress(field)
			case "BillingCity", "ShippingCity", "PersonMailingCity", "PersonOtherCity":
				mf = types.NewCity(field)
			case "BillingState", "ShippingState", "PersonMailingState", "PersonOtherState":
				mf = types.NewState(field)
			case "BillingCountry", "ShippingCountry", "PersonMailingCountry", "PersonOtherCountry":
				mf = types.NewCountry(field)
			case "BillingLatitude", "ShippingLatitude", "PersonMailingLatitude", "PersonOtherLatitude":
				mf = types.NewLatitude(field)
			case "BillingLongitude", "ShippingLongitude", "PersonMailingLongitude", "PersonOtherLongitude":
				mf = types.NewLongitude(field)
			case "FirstName":
				mf = types.NewFirstName(field)
			case "LastName":
				mf = types.NewLastName(field)
			case "BillingPostalCode", "ShippingPostalCode", "PersonMailingPostalCode", "PersonOtherPostalCode":
				mf = types.NewPostalCode(field)
			case "PersonEmail":
				mf = types.NewEmailAddress(field)
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

/*
	function returns a collection of fields that represent either a Company or Person Account.
*/
func handlePersonAccounts(fields []types.IField, personAccounts bool) []types.IField {

	var newFields []types.IField
	for _, f := range fields {
		if personAccounts && (strings.Contains(f.GetField().Name, "__pc") || strings.Index(f.GetField().Name, "Person") == 0) {
			newFields = append(newFields, f)
		} else if !personAccounts && !(strings.Contains(f.GetField().Name, "__pc") || strings.Index(f.GetField().Name, "Person") == 0) {
			newFields = append(newFields, f)
		}
	}
	return newFields
}
