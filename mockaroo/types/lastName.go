package types

type LastName struct {
	*Field
}

func (l LastName) GetField() *Field {
	return l.Field
}
func (l LastName) SetFormula(f string) {
	l.Formula = f
}
func NewLastName(m map[string]interface{}) *LastName {
	return &LastName{
		Field: &Field{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Last Name",
		},
	}
}
