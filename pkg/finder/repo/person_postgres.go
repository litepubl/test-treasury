package repo

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/litepubl/test-treasury/pkg/entity"
	"github.com/litepubl/test-treasury/pkg/postgres"
)

const defaultEntityCap = 64

// PersonRepo -.
type PersonRepo struct {
	*postgres.Postgres
}

// New -.
func New(pg *postgres.Postgres) *PersonRepo {
	return &PersonRepo{pg}
}

func (r *PersonRepo) Names(
	ctx context.Context,
	firstName string,
	lastName string,
	strong bool) ([]entity.Person, error) {
	sql, args, err := r.Builder.
		Select("*").
		From("persons").
		Where(r.buildWhere(firstName, lastName, strong)).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("PersonRepo - GetPerson - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("PersonRepo - GetPerson - r.Pool.Query: %w", err)
	}

	defer rows.Close()

	entities := make([]entity.Person, 0, defaultEntityCap)

	for rows.Next() {
		e := entity.Person{}

		err = rows.Scan(&e.Uid, &e.FirstName, &e.LastName)
		if err != nil {
			return nil, fmt.Errorf("PersonRepo - GetPerson - rows.Scan: %w", err)
		}

		entities = append(entities, e)
	}

	return entities, nil
}

func (r *PersonRepo) buildWhere(
	firstName string,
	lastName string,
	strong bool) any {
	if strong {
		if lastName == "" {
			return sq.Eq{"first_name": firstName}
		} else {
			return sq.Eq{"first_name": firstName, "last_name": lastName}
		}
	} else {
		if lastName == "" {
			return sq.Or{
				sq.Eq{"first_name": firstName},
				sq.Eq{"last_name": firstName},
			}
		} else {
			return sq.Or{
				sq.Eq{"first_name": firstName},
				sq.Eq{"first_name": lastName},
				sq.Eq{"last_name": firstName},
				sq.Eq{"last_name": lastName},
			}
		}
	}
}
