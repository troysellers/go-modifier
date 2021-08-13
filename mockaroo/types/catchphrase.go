package types

type CatchPhrase struct {
	*Field
}

func (c CatchPhrase) GetField() *Field {
	return c.Field
}
func (c CatchPhrase) SetFormula(f string) {
	c.Formula = f
}
func NewCatchPhrase(m map[string]interface{}) *CatchPhrase {
	return &CatchPhrase{
		Field: &Field{
			FieldType:  "Catch Phrase",
			Name:       m["name"].(string),
			SforceMeta: m,
		},
	}
}
