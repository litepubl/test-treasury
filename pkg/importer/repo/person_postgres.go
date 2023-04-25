package repo

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/litepubl/test-treasury/pkg/entity"
	"github.com/litepubl/test-treasury/pkg/postgres"
	"github.com/rs/zerolog/log"
)

const _defaultEntityCap = 64

// PersonRepo -.
type PersonRepo struct {
	*postgres.Postgres
}

// New -.
func New(pg *postgres.Postgres) *PersonRepo {
	return &PersonRepo{pg}
}

// Store -.
func (r *PersonRepo) Store(ctx context.Context, p *entity.Person) error {
	sql, args, err := r.Builder.
		Insert("persons").
		Columns("uid, first_name, last_name").
		Values(p.Uid, p.FirstName, p.LastName).
		ToSql()
	if err != nil {
		return fmt.Errorf("PersonRepo - Store - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("PersonRepo - Store - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *PersonRepo) Update(ctx context.Context, p *entity.Person) error {
	sql, args, err := r.Builder.
		Update("persons").
		Set("first_name", p.FirstName).
		Set("last_name", p.LastName).
		Where(sq.Eq{"uid": p.Uid}).
		ToSql()
	if err != nil {
		return fmt.Errorf("PersonRepo - Update - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("PersonRepo - Update - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *PersonRepo) DeleteById(ctx context.Context, idList []int) error {
	sql, args, err := r.Builder.
		Delete("persons").
		Where(sq.Eq{"uid": idList}).
		ToSql()

	if err != nil {
		return fmt.Errorf("PersonRepo - DeeleteById - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("PersonRepo - Store - r.Pool.Exec: %w", err)
	}

	return nil
}

// GetAll -.
func (r *PersonRepo) GetAll(ctx context.Context) (map[int]entity.Person, error) {
	sql, _, err := r.Builder.
		Select("uid, first_name, last_name").
		From("persons").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("PersonRepo - GetPerson - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		log.Info().Err(err).Msg("GetAll persons")
		return nil, fmt.Errorf("PersonRepo - GetAll - r.Pool.Query: %w", err)
	}

	defer rows.Close()

	entities := make(map[int]entity.Person, _defaultEntityCap)

	for rows.Next() {
		e := entity.Person{}

		err = rows.Scan(&e.Uid, &e.FirstName, &e.LastName)
		if err != nil {
			return nil, fmt.Errorf("PersonRepo - GetPerson - rows.Scan: %w", err)
		}

		entities[e.Uid] = e
	}

	return entities, nil
}
