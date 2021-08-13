package types

/*
	Define a base set of attributes for a mockaroo field and the interface
	that all types must implement

	https://www.mockaroo.com/docs
*/

type Field struct {
	Name         string                 `json:"name"`         // name of the field
	PercentBlank int                    `json:"percentBlank"` // integer between 0 and 100
	Formula      string                 `json:"formula"`      // formula to alter mockaroo values
	FieldType    string                 `json:"type"`         // one of mockaroo types
	SforceMeta   map[string]interface{} `json:"-"`            // the salesforce metadata for this field
}

type IField interface {
	GetField() *Field
	SetFormula(f string)
}
