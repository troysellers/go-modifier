package types

type DUNSNumber struct {
	*Field
}

func (d DUNSNumber) GetField() *Field {
	return d.Field
}
func (d DUNSNumber) SetFormula(f string) {
	d.Formula = f
}
func NewDUNSNumber(m map[string]interface{}) *DUNSNumber {
	return &DUNSNumber{
		Field: &Field{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "DUNS Number",
		},
	}
}
