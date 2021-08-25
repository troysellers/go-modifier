package types

type ConstructionSubContract struct {
	*Field
}

func (c ConstructionSubContract) GetField() *Field {
	return c.Field
}
func (c ConstructionSubContract) SetFormula(f string) {
	c.Formula = f
}
func NewConstructionSubContract(m map[string]interface{}) *ConstructionSubContract {
	return &ConstructionSubContract{
		Field: &Field{
			FieldType:  "Construction Subcontract Category",
			Name:       m["name"].(string),
			SforceMeta: m,
		},
	}
}
