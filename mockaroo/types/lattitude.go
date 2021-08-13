package types

type Latitude struct {
	*Field
}

func (l Latitude) GetField() *Field {
	return l.Field
}
func (l Latitude) SetFormula(f string) {
	l.Formula = f
}
func NewLatitude(m map[string]interface{}) *Latitude {
	return &Latitude{
		Field: &Field{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Latitude",
		},
	}
}
