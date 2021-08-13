package types

type FullName struct {
	*Field
}

func (fn FullName) GetField() *Field {
	return fn.Field
}
func (fn FullName) SetFormula(f string) {
	fn.Formula = f
}
func NewFullName(m map[string]interface{}) *FullName {
	return &FullName{
		Field: &Field{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Full Name",
		},
	}
}
