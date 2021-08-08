package mockaroo

import "time"

//https://www.mockaroo.com/docs
type FieldSpec struct {
	Name         string `json:"name"`         // name of the field
	PercentBlank int    `json:"percentBlank"` // integer between 0 and 100
	Formula      string `json:"formula"`      // formula to alter mockaroo values
	FieldType    string `json:"type"`         // one of mockaroo types
}

type GUID struct {
	FieldSpec `default:"{\"FieldType\":\"GUID\"}"`
}

func NewGUID(n string) *GUID {
	return &GUID{
		FieldSpec: FieldSpec{
			FieldType: "GUID",
			Name:      n,
		},
	}
}

type FakeCompanyName struct {
	FieldSpec
}

func NewFakeCompanyName(n string) *FakeCompanyName {
	return &FakeCompanyName{
		FieldSpec: FieldSpec{
			FieldType: "Fake Company Name",
			Name:      n,
		},
	}
}

type CustomList struct {
	FieldSpec
	Distribution   string   `json:"distribution"`
	SelectionStype string   `json:"selectionStyle"`
	Values         []string `json:"values"`
}

func NewCustomList(n string) *CustomList {
	return &CustomList{
		FieldSpec: FieldSpec{
			FieldType: "Custom List",
			Name:      n,
		},
		SelectionStype: "random",
	}
}

type StreetName struct {
	FieldSpec
}

func NewStreetName(n string) *StreetName {
	return &StreetName{
		FieldSpec: FieldSpec{
			Name:      n,
			FieldType: "Street Name",
		},
	}
}

type City struct {
	FieldSpec
}

func NewCity(n string) *City {
	return &City{
		FieldSpec: FieldSpec{
			Name:      n,
			FieldType: "City",
		},
	}
}

type State struct {
	FieldSpec
	OnlyUS bool `json:"onlyUSPlaces"`
}

func NewState(n string) *State {
	return &State{
		FieldSpec: FieldSpec{
			Name:      n,
			FieldType: "State",
		},
		OnlyUS: false,
	}
}

type Country struct {
	FieldSpec
	RestrictTo []string `json:"countries"`
}

func NewCountry(n string) *Country {
	return &Country{
		FieldSpec: FieldSpec{
			Name:      n,
			FieldType: "Country",
		},
		RestrictTo: make([]string, 0),
	}
}

type PostalCode struct {
	FieldSpec
}

func NewPostalCode(n string) *PostalCode {
	return &PostalCode{
		FieldSpec: FieldSpec{
			Name:      n,
			FieldType: "Postal Code",
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
	FieldSpec
	Format string `json:"format"`
}

func NewPhone(n string) *Phone {
	return &Phone{
		FieldSpec: FieldSpec{
			Name:      n,
			FieldType: "Phone",
		},
		Format: "+# ### ### ####",
	}
}

type URL struct {
	FieldSpec
	IncludeHost        bool `json:"includeHost"`
	IncludePath        bool `json:"includePath"`
	IncludeProtocol    bool `json:"includeProtocol"`
	IncludeQueryString bool `json:"includeQueryString"`
}

func NewURL(n string) *URL {
	return &URL{
		FieldSpec: FieldSpec{
			Name:      n,
			FieldType: "URL",
		},
		IncludeHost:        true,
		IncludePath:        true,
		IncludeProtocol:    true,
		IncludeQueryString: false,
	}
}

type Number struct {
	FieldSpec
	Decimals int `json:"decimals"`
	Max      int `json:"max"`
	Min      int `json:"min"`
}

func NewNumber(n string) *Number {
	return &Number{
		FieldSpec: FieldSpec{
			Name:      n,
			FieldType: "Number",
		},
		Decimals: 0,
		Max:      100,
		Min:      0,
	}
}

type Sentences struct {
	FieldSpec
	Max int `json:"max"`
	Min int `json:"min"`
}

func NewSentences(n string) *Sentences {
	return &Sentences{
		FieldSpec: FieldSpec{
			Name:      n,
			FieldType: "Sentences",
		},
		Max: 1,
		Min: 1,
	}
}

type DUNSNumber struct {
	FieldSpec
}

func NewDUNSNumber(n string) *DUNSNumber {
	return &DUNSNumber{
		FieldSpec: FieldSpec{
			Name:      n,
			FieldType: "DUNS Number",
		},
	}
}

type Words struct {
	FieldSpec
	Max int `json:"max"`
	Min int `json:"min"`
}

func NewWords(n string) *Words {
	return &Words{
		FieldSpec: FieldSpec{
			Name:      n,
			FieldType: "Words",
		},
		Max: 5,
		Min: 1,
	}
}

type DigitSequence struct {
	FieldSpec
	Format string `json:"format"`
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
func NewDigitSequence(n string) *DigitSequence {
	return &DigitSequence{
		FieldSpec: FieldSpec{
			Name:      n,
			FieldType: "Digit Sequence",
		},
	}
}

type Buzzword struct {
	FieldSpec
}

func NewBuzzword(n string) *Buzzword {
	return &Buzzword{
		FieldSpec: FieldSpec{
			FieldType: "Buzzword",
			Name:      n,
		},
	}
}

type CatchPhrase struct {
	FieldSpec
}

func NewCatchPhrase(n string) *CatchPhrase {
	return &CatchPhrase{
		FieldSpec: FieldSpec{
			FieldType: "Catch Phrase",
			Name:      n,
		},
	}
}

type Boolean struct {
	FieldSpec
}

func NewBoolean(n string) *Boolean {
	return &Boolean{
		FieldSpec: FieldSpec{
			FieldType: "Boolean",
			Name:      n,
		},
	}
}

type StreetAddress struct {
	FieldSpec
}

func NewStreetAddress(n string) *StreetAddress {
	return &StreetAddress{
		FieldSpec: FieldSpec{
			Name:      n,
			FieldType: "Street Address",
		},
	}
}

type Datetime struct {
	FieldSpec
	Max string `json:"max"`
	Min string `json:"min"`
}

func NewDatetime(n string) *Datetime {

	return &Datetime{
		FieldSpec: FieldSpec{
			Name:      n,
			FieldType: "Datetime",
		},
		Min: time.Now().AddDate(-1, 0, 0).Format("01/02/2006"),
		Max: time.Now().AddDate(1, 0, 0).Format("01/02/2006"),
	}
}

type EmailAddress struct {
	FieldSpec
}

func NewEmailAddress(n string) *EmailAddress {
	return &EmailAddress{
		FieldSpec: FieldSpec{
			Name:      n,
			FieldType: "Email Address",
		},
	}
}

type FirstName struct {
	FieldSpec
}

func NewFirstName(n string) *FirstName {
	return &FirstName{
		FieldSpec: FieldSpec{
			Name:      n,
			FieldType: "First Name",
		},
	}
}

type LastName struct {
	FieldSpec
}

func NewLastName(n string) *LastName {
	return &LastName{
		FieldSpec: FieldSpec{
			Name:      n,
			FieldType: "Last Name",
		},
	}
}

type JobTitle struct {
	FieldSpec
}

func NewJobTitle(n string) *JobTitle {
	return &JobTitle{
		FieldSpec: FieldSpec{
			Name:      n,
			FieldType: "Job Title",
		},
	}
}

type StockSymbol struct {
	FieldSpec
}

func NewStockSymbol(n string) *StockSymbol {
	return &StockSymbol{
		FieldSpec: FieldSpec{
			Name:      n,
			FieldType: "Stock Symbol",
		},
	}
}

type Latitude struct {
	FieldSpec
}

func NewLatitude(n string) *Latitude {
	return &Latitude{
		FieldSpec: FieldSpec{
			Name:      n,
			FieldType: "Latitude",
		},
	}
}

type Longitude struct {
	FieldSpec
}

func NewLongitude(n string) *Longitude {
	return &Longitude{
		FieldSpec: FieldSpec{
			Name:      n,
			FieldType: "Longitude",
		},
	}
}
