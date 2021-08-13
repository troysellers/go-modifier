package types

import "time"

type Datetime struct {
	*Field
	Max string `json:"max"`
	Min string `json:"min"`
}

func (d Datetime) GetField() *Field {
	return d.Field
}
func (d Datetime) SetFormula(f string) {
	d.Formula = f
}

func NewDatetime(m map[string]interface{}) *Datetime {

	return &Datetime{
		Field: &Field{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Datetime",
		},
		Min: time.Now().AddDate(-1, 0, 0).Format("01/02/2006"),
		Max: time.Now().AddDate(1, 0, 0).Format("01/02/2006"),
	}
}
