package service

import (
	"context"
	"math/rand"
	"strings"

	"soccer-api/internal/domain/entity"
)

type Country struct {
	countryList entity.ListOfCountries
	countries   []entity.Country
}

func New() *Country {
	return &Country{
		countryList: entity.CountryList,
		countries:   entity.Countries,
	}
}

func (c *Country) GetCountries(ctx context.Context) []entity.Country {
	return c.countries
}

func (c *Country) GetCountryByName(ctx context.Context, name string) (entity.Country, bool) {
	countryName := strings.ToLower(strings.TrimSpace(name))
	countryName = strings.ToUpper(countryName[:1]) + countryName[1:]
	country, exists := c.countryList[countryName]
	return country, exists
}

func (c *Country) Random() entity.Country {
	return c.countries[rand.Intn(len(c.countries))]
}
