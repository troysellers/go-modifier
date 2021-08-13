package types

type Country struct {
	*Field
	RestrictTo []string `json:"countries"`
}

func (c Country) GetField() *Field {
	return c.Field
}
func (c Country) SetFormula(f string) {
	c.Formula = f
}
func NewCountry(m map[string]interface{}) *Country {
	return &Country{
		Field: &Field{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Country",
		},
		RestrictTo: make([]string, 0),
	}
}
