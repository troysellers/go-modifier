package types

type Number struct {
	*Field
	Decimals int `json:"decimals"`
	Max      int `json:"max"`
	Min      int `json:"min"`
}

func (n Number) GetField() *Field {
	return n.Field
}
func (n Number) SetFormula(f string) {
	n.Formula = f
}
func NewNumber(m map[string]interface{}) *Number {
	return &Number{
		Field: &Field{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Number",
		},
		Decimals: 0,
		Max:      100,
		Min:      0,
	}
}
