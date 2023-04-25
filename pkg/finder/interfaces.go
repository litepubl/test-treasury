// Package finder implements application business logic. Each logic group in own file.
package finder

import (
	"context"

	"github.com/litepubl/test-treasury/pkg/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=importer_test
type  (
	DataFinder interface {
GetNames(ctx context.Context, name string, strong bool) ([]entity.Person, error)
	}

	PersonRepo interface {
GetNames(ctx context.Context, firstName string, lastName string, strong bool) ([]entity.Person, error)
}
)
