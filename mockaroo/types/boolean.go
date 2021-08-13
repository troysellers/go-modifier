package types

type Boolean struct {
	*Field
}

func (b Boolean) GetField() *Field {
	return b.Field
}
func (b Boolean) SetFormula(f string) {
	b.Formula = f
}
func NewBoolean(m map[string]interface{}) *Boolean {
	return &Boolean{
		Field: &Field{
			FieldType:  "Boolean",
			Name:       m["name"].(string),
			SforceMeta: m,
		},
	}
}
