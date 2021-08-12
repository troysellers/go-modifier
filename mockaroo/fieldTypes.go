package mockaroo

import (
	"time"
)

//https://www.mockaroo.com/docs
type FieldSpec struct {
	Name         string                 `json:"name"`         // name of the field
	PercentBlank int                    `json:"percentBlank"` // integer between 0 and 100
	Formula      string                 `json:"formula"`      // formula to alter mockaroo values
	FieldType    string                 `json:"type"`         // one of mockaroo types
	SforceMeta   map[string]interface{} `json:"-"`            // the salesforce metadata for this field
}

type FieldSpecInterface interface {
	GetFieldSpec() *FieldSpec
}

type GUID struct {
	*FieldSpec
}

func (g GUID) GetFieldSpec() *FieldSpec {
	return g.FieldSpec
}

func NewGUID(m map[string]interface{}) *GUID {
	return &GUID{
		FieldSpec: &FieldSpec{
			FieldType:  "GUID",
			Name:       m["name"].(string),
			SforceMeta: m,
		},
	}
}

type FakeCompanyName struct {
	*FieldSpec
}

func (f FakeCompanyName) GetFieldSpec() *FieldSpec {
	return f.FieldSpec
}
func NewFakeCompanyName(m map[string]interface{}) *FakeCompanyName {
	return &FakeCompanyName{
		FieldSpec: &FieldSpec{
			FieldType:  "Fake Company Name",
			Name:       m["name"].(string),
			SforceMeta: m,
		},
	}
}

type CustomList struct {
	*FieldSpec
	Distribution   string   `json:"distribution"`
	SelectionStype string   `json:"selectionStyle"`
	Values         []string `json:"values"`
}

func (c CustomList) GetFieldSpec() *FieldSpec {
	return c.FieldSpec
}
func NewCustomList(m map[string]interface{}) *CustomList {
	return &CustomList{
		FieldSpec: &FieldSpec{
			FieldType:  "Custom List",
			Name:       m["name"].(string),
			SforceMeta: m,
		},
		SelectionStype: "random",
	}
}

type StreetName struct {
	*FieldSpec
}

func (s StreetName) GetFieldSpec() *FieldSpec {
	return s.FieldSpec
}
func NewStreetName(m map[string]interface{}) *StreetName {
	return &StreetName{
		FieldSpec: &FieldSpec{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Street Name",
		},
	}
}

type City struct {
	*FieldSpec
}

func (c City) GetFieldSpec() *FieldSpec {
	return c.FieldSpec
}
func NewCity(m map[string]interface{}) *City {
	return &City{
		FieldSpec: &FieldSpec{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "City",
		},
	}
}

type State struct {
	*FieldSpec
	OnlyUS bool `json:"onlyUSPlaces"`
}

func (s State) GetFieldSpec() *FieldSpec {
	return s.FieldSpec
}
func NewState(m map[string]interface{}) *State {
	return &State{
		FieldSpec: &FieldSpec{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "State",
		},
		OnlyUS: false,
	}
}

type Country struct {
	*FieldSpec
	RestrictTo []string `json:"countries"`
}

func (c Country) GetFieldSpec() *FieldSpec {
	return c.FieldSpec
}
func NewCountry(m map[string]interface{}) *Country {
	return &Country{
		FieldSpec: &FieldSpec{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Country",
		},
		RestrictTo: make([]string, 0),
	}
}

type PostalCode struct {
	*FieldSpec
}

func (s PostalCode) GetFieldSpec() *FieldSpec {
	return s.FieldSpec
}
func NewPostalCode(m map[string]interface{}) *PostalCode {
	return &PostalCode{
		FieldSpec: &FieldSpec{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Postal Code",
		},
	}
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
type Phone struct {
	*FieldSpec
	Format string `json:"format"`
}

func (p Phone) GetFieldSpec() *FieldSpec {
	return p.FieldSpec
}

func NewPhone(m map[string]interface{}) *Phone {
	return &Phone{
		FieldSpec: &FieldSpec{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Phone",
		},
		Format: "+# ### ### ####",
	}
}

type URL struct {
	*FieldSpec
	IncludeHost        bool `json:"includeHost"`
	IncludePath        bool `json:"includePath"`
	IncludeProtocol    bool `json:"includeProtocol"`
	IncludeQueryString bool `json:"includeQueryString"`
}

func (u URL) GetFieldSpec() *FieldSpec {
	return u.FieldSpec
}
func NewURL(m map[string]interface{}) *URL {
	return &URL{
		FieldSpec: &FieldSpec{
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

type Number struct {
	*FieldSpec
	Decimals int `json:"decimals"`
	Max      int `json:"max"`
	Min      int `json:"min"`
}

func (n Number) GetFieldSpec() *FieldSpec {
	return n.FieldSpec
}
func NewNumber(m map[string]interface{}) *Number {
	return &Number{
		FieldSpec: &FieldSpec{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Number",
		},
		Decimals: 0,
		Max:      100,
		Min:      0,
	}
}

type Sentences struct {
	*FieldSpec
	Max int `json:"max"`
	Min int `json:"min"`
}

func (s Sentences) GetFieldSpec() *FieldSpec {
	return s.FieldSpec
}
func NewSentences(m map[string]interface{}) *Sentences {
	return &Sentences{
		FieldSpec: &FieldSpec{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Sentences",
		},
		Max: 1,
		Min: 1,
	}
}

type DUNSNumber struct {
	*FieldSpec
}

func (d DUNSNumber) GetFieldSpec() *FieldSpec {
	return d.FieldSpec
}
func NewDUNSNumber(m map[string]interface{}) *DUNSNumber {
	return &DUNSNumber{
		FieldSpec: &FieldSpec{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "DUNS Number",
		},
	}
}

type Words struct {
	*FieldSpec
	Max int `json:"max"`
	Min int `json:"min"`
}

func (w Words) GetFieldSpec() *FieldSpec {
	return w.FieldSpec
}
func NewWords(m map[string]interface{}) *Words {
	return &Words{
		FieldSpec: &FieldSpec{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Words",
		},
		Max: 5,
		Min: 1,
	}
}

type DigitSequence struct {
	*FieldSpec
	Format string `json:"format"`
}

func (d DigitSequence) GetFieldSpec() *FieldSpec {
	return d.FieldSpec
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
		FieldSpec: &FieldSpec{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Digit Sequence",
		},
	}
}

type Buzzword struct {
	*FieldSpec
}

func (b Buzzword) GetFieldSpec() *FieldSpec {
	return b.FieldSpec
}
func NewBuzzword(m map[string]interface{}) *Buzzword {
	return &Buzzword{
		FieldSpec: &FieldSpec{
			FieldType:  "Buzzword",
			Name:       m["name"].(string),
			SforceMeta: m,
		},
	}
}

type CatchPhrase struct {
	*FieldSpec
}

func (c CatchPhrase) GetFieldSpec() *FieldSpec {
	return c.FieldSpec
}
func NewCatchPhrase(m map[string]interface{}) *CatchPhrase {
	return &CatchPhrase{
		FieldSpec: &FieldSpec{
			FieldType:  "Catch Phrase",
			Name:       m["name"].(string),
			SforceMeta: m,
		},
	}
}

type Boolean struct {
	*FieldSpec
}

func (b Boolean) GetFieldSpec() *FieldSpec {
	return b.FieldSpec
}
func NewBoolean(m map[string]interface{}) *Boolean {
	return &Boolean{
		FieldSpec: &FieldSpec{
			FieldType:  "Boolean",
			Name:       m["name"].(string),
			SforceMeta: m,
		},
	}
}

type StreetAddress struct {
	*FieldSpec
}

func (s StreetAddress) GetFieldSpec() *FieldSpec {
	return s.FieldSpec
}
func NewStreetAddress(m map[string]interface{}) *StreetAddress {
	return &StreetAddress{
		FieldSpec: &FieldSpec{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Street Address",
		},
	}
}

type Datetime struct {
	*FieldSpec
	Max string `json:"max"`
	Min string `json:"min"`
}

func (d Datetime) GetFieldSpec() *FieldSpec {
	return d.FieldSpec
}
func NewDatetime(m map[string]interface{}) *Datetime {

	return &Datetime{
		FieldSpec: &FieldSpec{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Datetime",
		},
		Min: time.Now().AddDate(-1, 0, 0).Format("01/02/2006"),
		Max: time.Now().AddDate(1, 0, 0).Format("01/02/2006"),
	}
}

type EmailAddress struct {
	*FieldSpec
}

func (e EmailAddress) GetFieldSpec() *FieldSpec {
	return e.FieldSpec
}
func NewEmailAddress(m map[string]interface{}) *EmailAddress {
	return &EmailAddress{
		FieldSpec: &FieldSpec{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Email Address",
		},
	}
}

type FirstName struct {
	*FieldSpec
}

func (f FirstName) GetFieldSpec() *FieldSpec {
	return f.FieldSpec
}
func NewFirstName(m map[string]interface{}) *FirstName {
	return &FirstName{
		FieldSpec: &FieldSpec{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "First Name",
		},
	}
}

type LastName struct {
	*FieldSpec
}

func (l LastName) GetFieldSpec() *FieldSpec {
	return l.FieldSpec
}
func NewLastName(m map[string]interface{}) *LastName {
	return &LastName{
		FieldSpec: &FieldSpec{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Last Name",
		},
	}
}

type JobTitle struct {
	*FieldSpec
}

func (j JobTitle) GetFieldSpec() *FieldSpec {
	return j.FieldSpec
}
func NewJobTitle(m map[string]interface{}) *JobTitle {
	return &JobTitle{
		FieldSpec: &FieldSpec{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Job Title",
		},
	}
}

type StockSymbol struct {
	*FieldSpec
}

func (s StockSymbol) GetFieldSpec() *FieldSpec {
	return s.FieldSpec
}
func NewStockSymbol(m map[string]interface{}) *StockSymbol {
	return &StockSymbol{
		FieldSpec: &FieldSpec{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Stock Symbol",
		},
	}
}

type Latitude struct {
	*FieldSpec
}

func (l Latitude) GetFieldSpec() *FieldSpec {
	return l.FieldSpec
}
func NewLatitude(m map[string]interface{}) *Latitude {
	return &Latitude{
		FieldSpec: &FieldSpec{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Latitude",
		},
	}
}

type Longitude struct {
	*FieldSpec
}

func (l Longitude) GetFieldSpec() *FieldSpec {
	return l.FieldSpec
}
func NewLongitude(m map[string]interface{}) *Longitude {
	return &Longitude{
		FieldSpec: &FieldSpec{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Longitude",
		},
	}
}

type FullName struct {
	*FieldSpec
}

func (f FullName) GetFieldSpec() *FieldSpec {
	return f.FieldSpec
}
func NewFullName(m map[string]interface{}) *FullName {
	return &FullName{
		FieldSpec: &FieldSpec{
			Name:       m["name"].(string),
			SforceMeta: m,
			FieldType:  "Full Name",
		},
	}
}
