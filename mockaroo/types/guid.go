package types

type GUID struct {
	*Field
}

func (g GUID) GetField() *Field {
	return g.Field
}
func (g GUID) SetFormula(f string) {
	g.Formula = f
}

func NewGUID(m map[string]interface{}) *GUID {
	return &GUID{
		Field: &Field{
			FieldType:  "GUID",
			Name:       m["name"].(string),
			SforceMeta: m,
		},
	}
}
