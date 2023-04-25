package finder

import (
	"context"
	"strings"

	"github.com/litepubl/test-treasury/pkg/entity"
)

type Finder struct {
	repo PersonRepo
}

func New(r PersonRepo) *Finder {
	return &Finder{
		repo: r,
	}
}

func (finder *Finder) GetNames(ctx context.Context, name string, strong bool) ([]entity.Person, error) {
	names := strings.Split(name, " ")
	firstName := strings.TrimSpace(names[0])
	lastName := ""
	if len(names) > 1 {
		lastName = strings.TrimSpace(names[1])
	}

	return finder.repo.GetNames(context.Background(), firstName, lastName, strong)
}
