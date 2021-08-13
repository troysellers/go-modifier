package types

type URL struct {
	*Field
	IncludeHost        bool `json:"includeHost"`
	IncludePath        bool `json:"includePath"`
	IncludeProtocol    bool `json:"includeProtocol"`
	IncludeQueryString bool `json:"includeQueryString"`
}

func (u URL) GetField() *Field {
	return u.Field
}
func (u URL) SetFormula(f string) {
	u.Formula = f
}
func NewURL(m map[string]interface{}) *URL {
	return &URL{
		Field: &Field{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "URL",
		},
		IncludeHost:        true,
		IncludePath:        true,
		IncludeProtocol:    true,
		IncludeQueryString: false,
	}
}
