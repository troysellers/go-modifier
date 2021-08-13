package types

type State struct {
	*Field
	OnlyUS bool `json:"onlyUSPlaces"`
}

func (s State) GetField() *Field {
	return s.Field
}
func (s State) SetFormula(f string) {
	s.Formula = f
}
func NewState(m map[string]interface{}) *State {
	return &State{
		Field: &Field{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "State",
		},
		OnlyUS: false,
	}
}
