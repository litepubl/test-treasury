// Package importer implements application business logic. Each logic group in own file.
package importer

import (
	"context"

	"github.com/litepubl/test-treasury/pkg/entity"
	"github.com/litepubl/test-treasury/pkg/xmldata"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=importer_test

type (
	DataUpdater interface {
		Update(ctx context.Context) error
		GetState() State
	}

	DataImporter interface {
		Import(ctx context.Context) error
	}

	PersonRepo interface {
		Store(ctx context.Context, p *entity.Person) error
		Update(ctx context.Context, p *entity.Person) error
		DeleteById(ctx context.Context, idList []int) error
		GetAll(ctx context.Context) (map[int]entity.Person, error)
	}

	Downloader interface {
		Download(ctx context.Context) (*xmldata.SdnList, error)
	}
)
