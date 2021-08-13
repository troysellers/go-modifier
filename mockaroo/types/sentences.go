package types

type Sentences struct {
	*Field
	Max int `json:"max"`
	Min int `json:"min"`
}

func (s Sentences) GetField() *Field {
	return s.Field
}
func (s Sentences) SetFormula(f string) {
	s.Formula = f
}
func NewSentences(m map[string]interface{}) *Sentences {
	return &Sentences{
		Field: &Field{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Sentences",
		},
		Max: 1,
		Min: 1,
	}
}
