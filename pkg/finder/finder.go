package finder

import (
	"context"
	"strings"

	"github.com/litepubl/test-treasury/pkg/entity"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Service Finder provide search persons in database
type Finder struct {
	repo PersonRepo
}

// Finder constructor
func New(r PersonRepo) *Finder {
	return &Finder{
		repo: r,
	}
}

// Implementation of DataFinder interface

func (finder *Finder) Names(ctx context.Context, name string, strong bool) ([]entity.Person, error) {
	names := strings.Split(name, " ")
	firstName := cases.Title(language.Und).String(strings.TrimSpace(names[0]))
	lastName := ""
	if len(names) > 1 {
		lastName = cases.Title(language.Und).String(strings.TrimSpace(names[1]))
	}

	return finder.repo.Names(context.Background(), firstName, lastName, strong)
}
