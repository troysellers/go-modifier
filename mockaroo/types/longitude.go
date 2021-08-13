package types

type Longitude struct {
	*Field
}

func (l Longitude) GetField() *Field {
	return l.Field
}
func (l Longitude) SetFormula(f string) {
	l.Formula = f
}
func NewLongitude(m map[string]interface{}) *Longitude {
	return &Longitude{
		Field: &Field{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Longitude",
		},
	}
}
