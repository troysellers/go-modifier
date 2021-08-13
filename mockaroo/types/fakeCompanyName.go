package types

type FakeCompanyName struct {
	*Field
}

func (fcn FakeCompanyName) GetField() *Field {
	return fcn.Field
}
func (fcn FakeCompanyName) SetFormula(f string) {
	fcn.Formula = f
}
func NewFakeCompanyName(m map[string]interface{}) *FakeCompanyName {
	return &FakeCompanyName{
		Field: &Field{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Fake Company Name",
		},
	}
}
