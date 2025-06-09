package entity

type Country string

const (
	GERMANY     Country = "Germany"
	BELGIUM     Country = "Belgium"
	FRANCE      Country = "France"
	PORTUGAL    Country = "Portugal"
	SPAIN       Country = "Spain"
	SCOTLAND    Country = "Scotland"
	TURKEY      Country = "Turkey"
	AUSTRIA     Country = "Austria"
	ENGLAND     Country = "England"
	HUNGARY     Country = "Hungary"
	SLOVAKIA    Country = "Slovakia"
	ALBANIA     Country = "Albania"
	DENMARK     Country = "Denmark"
	NETHERLANDS Country = "Netherlands"
	ROMANIA     Country = "Romania"
	SWITZERLAND Country = "Switzerland"
	SERBIA      Country = "Serbia"
	ITALY       Country = "Italy"
	CZECHIA     Country = "Czechia"
	SLOVENIA    Country = "Slovenia"
	CROATIA     Country = "Croatia"
	GEORGIA     Country = "Georgia"
	UKRAINE     Country = "Ukraine"
	POLAND      Country = "Poland"
)

type ListOfCountries map[string]Country

var CountryList = ListOfCountries{
	"Germany":        GERMANY,
	"Belgium":        BELGIUM,
	"France":         FRANCE,
	"Portugal":       PORTUGAL,
	"Spain":          SPAIN,
	"Scotland":       SCOTLAND,
	"Turkey":         TURKEY,
	"Austria":        AUSTRIA,
	"England":        ENGLAND,
	"Hungary":        HUNGARY,
	"Slovakia":       SLOVAKIA,
	"Albania":        ALBANIA,
	"Denmark":        DENMARK,
	"Netherlands":    NETHERLANDS,
	"Romania":        ROMANIA,
	"Switzerland":    SWITZERLAND,
	"Serbia":         SERBIA,
	"Italy":          ITALY,
	"Czechia":        CZECHIA,
	"Czech Republic": CZECHIA,
	"Slovenia":       SLOVENIA,
	"Croatia":        CROATIA,
	"Georgia":        GEORGIA,
	"Ukraine":        UKRAINE,
	"Poland":         POLAND,
}

var Countries = []Country{
	GERMANY,
	BELGIUM,
	FRANCE,
	PORTUGAL,
	SPAIN,
	SCOTLAND,
	TURKEY,
	AUSTRIA,
	ENGLAND,
	HUNGARY,
	SLOVAKIA,
	ALBANIA,
	DENMARK,
	NETHERLANDS,
	ROMANIA,
	SWITZERLAND,
	SERBIA,
	ITALY,
	CZECHIA,
	SLOVENIA,
	CROATIA,
	GEORGIA,
	UKRAINE,
	POLAND,
}
