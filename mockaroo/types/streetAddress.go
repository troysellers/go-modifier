package types

type StreetAddress struct {
	*Field
}

func (s StreetAddress) GetField() *Field {
	return s.Field
}
func (s StreetAddress) SetFormula(f string) {
	s.Formula = f
}
func NewStreetAddress(m map[string]interface{}) *StreetAddress {
	return &StreetAddress{
		Field: &Field{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Street Address",
		},
	}
}
