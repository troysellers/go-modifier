package types

type PostalCode struct {
	*Field
}

func (p PostalCode) GetField() *Field {
	return p.Field
}
func (p PostalCode) SetFormula(f string) {
	p.Formula = f
}
func NewPostalCode(m map[string]interface{}) *PostalCode {
	return &PostalCode{
		Field: &Field{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Postal Code",
		},
	}
}
