package types

type Phone struct {
	*Field
	Format string `json:"format"`
}

func (p Phone) GetField() *Field {
	return p.Field
}
func (p Phone) SetFormula(f string) {
	p.Formula = f
}

/*
Format must be one of these
	###-###-####
	(###) ###-####
	### ### ####
	+# ### ### ####
	+# (###) ###-####
	+#-###-###-####
	#-(###)###-####
	##########
*/
func NewPhone(m map[string]interface{}) *Phone {
	return &Phone{
		Field: &Field{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Phone",
		},
		Format: "+# ### ### ####",
	}
}
