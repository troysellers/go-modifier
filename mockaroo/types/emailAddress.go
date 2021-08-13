package types

type EmailAddress struct {
	*Field
}

func (e EmailAddress) GetField() *Field {
	return e.Field
}
func (e EmailAddress) SetFormula(f string) {
	e.Formula = f
}
func NewEmailAddress(m map[string]interface{}) *EmailAddress {
	return &EmailAddress{
		Field: &Field{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Email Address",
		},
	}
}
