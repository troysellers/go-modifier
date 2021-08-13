package types

type Buzzword struct {
	*Field
}

func (b Buzzword) GetField() *Field {
	return b.Field
}
func (b Buzzword) SetFormula(f string) {
	b.Formula = f
}
func NewBuzzword(m map[string]interface{}) *Buzzword {
	return &Buzzword{
		Field: &Field{
			FieldType:  "Buzzword",
			Name:       m["name"].(string),
			SforceMeta: m,
		},
	}
}
