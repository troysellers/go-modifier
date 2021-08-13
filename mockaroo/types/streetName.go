package types

type StreetName struct {
	*Field
}

func (s StreetName) GetField() *Field {
	return s.Field
}
func (s StreetName) SetFormula(f string) {
	s.Formula = f
}

func NewStreetName(m map[string]interface{}) *StreetName {
	return &StreetName{
		Field: &Field{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Street Name",
		},
	}
}
