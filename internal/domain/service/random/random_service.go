package random

import (
	"context"
	"math/rand"

	"github.com/jaswdr/faker/v2"
)

type Random struct {
	fkr faker.Faker
}

func New() *Random {
	return &Random{fkr: faker.New()}
}

func (p *Random) FirstName(ctx context.Context) string {
	return p.fkr.Person().FirstName()
}

func (p *Random) LastName(ctx context.Context) string {
	return p.fkr.Person().LastName()
}

func (p *Random) Age(ctx context.Context, min, max int) int {
	if min > max {
		min, max = max, min
	}
	return min + rand.Intn(max-min+1)
}
