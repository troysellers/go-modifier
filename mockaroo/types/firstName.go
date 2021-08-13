package types

type FirstName struct {
	*Field
}

func (fn FirstName) GetField() *Field {
	return fn.Field
}
func (fn FirstName) SetFormula(f string) {
	fn.Formula = f
}
func NewFirstName(m map[string]interface{}) *FirstName {
	return &FirstName{
		Field: &Field{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "First Name",
		},
	}
}
