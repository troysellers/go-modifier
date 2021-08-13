package types

type JobTitle struct {
	*Field
}

func (j JobTitle) GetField() *Field {
	return j.Field
}
func (j JobTitle) SetFormula(f string) {
	j.Formula = f
}
func NewJobTitle(m map[string]interface{}) *JobTitle {
	return &JobTitle{
		Field: &Field{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Job Title",
		},
	}
}
