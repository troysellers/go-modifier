package types

type CustomList struct {
	*Field
	Distribution   string   `json:"distribution"`
	SelectionStype string   `json:"selectionStyle"`
	Values         []string `json:"values"`
}

func (c CustomList) GetField() *Field {
	return c.Field
}

func (c CustomList) SetFormula(f string) {
	c.Formula = f
}

func NewCustomList(m map[string]interface{}) *CustomList {
	return &CustomList{
		Field: &Field{
			FieldType:  "Custom List",
			Name:       m["name"].(string),
			SforceMeta: m,
		},
		SelectionStype: "random",
	}
}
