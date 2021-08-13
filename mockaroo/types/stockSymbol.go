package types

type StockSymbol struct {
	*Field
}

func (s StockSymbol) GetField() *Field {
	return s.Field
}
func (s StockSymbol) SetFormula(f string) {
	s.Formula = f
}
func NewStockSymbol(m map[string]interface{}) *StockSymbol {
	return &StockSymbol{
		Field: &Field{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Stock Symbol",
		},
	}
}
