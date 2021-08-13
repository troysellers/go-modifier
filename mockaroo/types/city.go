package types

type City struct {
	*Field
}

func (c City) GetField() *Field {
	return c.Field
}
func (c City) SetFormula(f string) {
	c.Formula = f
}

func NewCity(m map[string]interface{}) *City {
	return &City{
		Field: &Field{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "City",
		},
	}
}
