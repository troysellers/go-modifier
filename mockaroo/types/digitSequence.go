package types

type DigitSequence struct {
	*Field
	Format string `json:"format"`
}

func (d DigitSequence) GetField() *Field {
	return d.Field
}
func (d DigitSequence) SetFormula(f string) {
	d.Formula = f
}

/*
Format
	Use "#" for a random digit.
	Use "@" for a random lower case letter.
	Use "^" for a random upper case letter.
	Use "*" for a random digit or letter.
	Use "$" for a random digit or lower case letter.
	Use "%" for a random digit or upper case letter.
	Any other character will be included verbatim.
	Examples

	###-##-#### => 232-66-7439
	***-## => A0c-34
	^222-##:### => Cght-87:485
*/
func NewDigitSequence(m map[string]interface{}) *DigitSequence {
	return &DigitSequence{
		Field: &Field{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Digit Sequence",
		},
	}
}
