package types

type Words struct {
	*Field
	Max int `json:"max"`
	Min int `json:"min"`
}

func (w Words) GetField() *Field {
	return w.Field
}
func (w Words) SetFormula(f string) {
	w.Formula = f
}
func NewWords(m map[string]interface{}) *Words {
	return &Words{
		Field: &Field{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Words",
		},
		Max: 5,
		Min: 1,
	}
}
