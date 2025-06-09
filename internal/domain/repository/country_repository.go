package repository

import "soccer-api/internal/domain/entity"

type Country interface {
	CheckCountry(country string) (entity.Country, bool)
	Random() entity.Country
	GetCountries() []entity.Country
}
