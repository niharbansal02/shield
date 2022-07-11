package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"

	"github.com/odpf/shield/internal/schema"
	"github.com/odpf/shield/model"
)

type Namespace struct {
	Id        string    `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func buildGetNamespaceQuery(dialect goqu.DialectWrapper) (string, error) {
	getNamespaceQuery, _, err := dialect.From("namespaces").Where(goqu.Ex{
		"id": goqu.L("$1"),
	}).ToSQL()

	return getNamespaceQuery, err
}
func buildCreateNamespaceQuery(dialect goqu.DialectWrapper) (string, error) {
	createNamespaceQuery, _, err := dialect.Insert("namespaces").Rows(
		goqu.Record{
			"id":   goqu.L("$1"),
			"name": goqu.L("$2"),
		}).OnConflict(goqu.DoUpdate("id", goqu.Record{
		"name": goqu.L("$2"),
	})).Returning(&Namespace{}).ToSQL()

	return createNamespaceQuery, err
}
func buildListNamespacesQuery(dialect goqu.DialectWrapper) (string, error) {
	listNamespacesQuery, _, err := dialect.From("namespaces").ToSQL()

	return listNamespacesQuery, err
}
func buildUpdateNamespaceQuery(dialect goqu.DialectWrapper) (string, error) {
	updateNamespaceQuery, _, err := dialect.Update("namespaces").Set(
		goqu.Record{
			"id":         goqu.L("$2"),
			"name":       goqu.L("$3"),
			"updated_at": goqu.L("now()"),
		}).Where(goqu.Ex{
		"id": goqu.L("$1"),
	}).Returning(&Namespace{}).ToSQL()

	return updateNamespaceQuery, err
}

func (s Store) GetNamespace(ctx context.Context, id string) (model.Namespace, error) {
	fetchedNamespace, err := s.selectNamespace(ctx, id, nil)
	return fetchedNamespace, err
}

func (s Store) selectNamespace(ctx context.Context, id string, txn *sqlx.Tx) (model.Namespace, error) {
	var fetchedNamespace Namespace
	getNamespaceQuery, err := buildGetNamespaceQuery(dialect)
	if err != nil {
		return model.Namespace{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	err = s.DB.WithTimeout(ctx, func(ctx context.Context) error {
		return s.DB.GetContext(ctx, &fetchedNamespace, getNamespaceQuery, id)
	})

	if errors.Is(err, sql.ErrNoRows) {
		return model.Namespace{}, schema.NamespaceDoesntExist
	} else if err != nil && fmt.Sprintf("%s", err.Error()[0:38]) == "pq: invalid input syntax for type uuid" {
		// TODO: this uuid syntax is a error defined in db, not in library
		// need to look into better ways to implement this
		return model.Namespace{}, schema.InvalidUUID
	} else if err != nil {
		return model.Namespace{}, fmt.Errorf("%w: %s", dbErr, err)
	}

	transformedNamespace, err := transformToNamespace(fetchedNamespace)
	if err != nil {
		return model.Namespace{}, fmt.Errorf("%w: %s", parseErr, err)
	}

	return transformedNamespace, nil
}

func (s Store) CreateNamespace(ctx context.Context, namespaceToCreate model.Namespace) (model.Namespace, error) {
	var newNamespace Namespace
	createNamespaceQuery, err := buildCreateNamespaceQuery(dialect)
	if err != nil {
		return model.Namespace{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	err = s.DB.WithTimeout(ctx, func(ctx context.Context) error {
		return s.DB.GetContext(ctx, &newNamespace, createNamespaceQuery, namespaceToCreate.Id, namespaceToCreate.Name)
	})

	if err != nil {
		return model.Namespace{}, fmt.Errorf("%w: %s", dbErr, err)
	}

	transformedNamespace, err := transformToNamespace(newNamespace)
	if err != nil {
		return model.Namespace{}, fmt.Errorf("%w: %s", parseErr, err)
	}

	return transformedNamespace, nil
}

func (s Store) ListNamespaces(ctx context.Context) ([]model.Namespace, error) {
	var fetchedNamespaces []Namespace
	listNamespacesQuery, err := buildListNamespacesQuery(dialect)
	if err != nil {
		return []model.Namespace{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	err = s.DB.WithTimeout(ctx, func(ctx context.Context) error {
		return s.DB.SelectContext(ctx, &fetchedNamespaces, listNamespacesQuery)
	})

	if errors.Is(err, sql.ErrNoRows) {
		return []model.Namespace{}, schema.NamespaceDoesntExist
	}

	if err != nil {
		return []model.Namespace{}, fmt.Errorf("%w: %s", dbErr, err)
	}

	var transformedNamespaces []model.Namespace

	for _, o := range fetchedNamespaces {
		transformedNamespace, err := transformToNamespace(o)
		if err != nil {
			return []model.Namespace{}, fmt.Errorf("%w: %s", parseErr, err)
		}

		transformedNamespaces = append(transformedNamespaces, transformedNamespace)
	}

	return transformedNamespaces, nil
}

func (s Store) UpdateNamespace(ctx context.Context, id string, toUpdate model.Namespace) (model.Namespace, error) {
	var updatedNamespace Namespace
	updateNamespaceQuery, err := buildUpdateNamespaceQuery(dialect)
	if err != nil {
		return model.Namespace{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	err = s.DB.WithTimeout(ctx, func(ctx context.Context) error {
		return s.DB.GetContext(ctx, &updatedNamespace, updateNamespaceQuery, id, toUpdate.Id, toUpdate.Name)
	})

	if errors.Is(err, sql.ErrNoRows) {
		return model.Namespace{}, schema.NamespaceDoesntExist
	} else if err != nil {
		return model.Namespace{}, fmt.Errorf("%w: %s", dbErr, err)
	}

	transformedNamespace, err := transformToNamespace(updatedNamespace)
	if err != nil {
		return model.Namespace{}, fmt.Errorf("%s: %w", parseErr, err)
	}

	return transformedNamespace, nil
}

func transformToNamespace(from Namespace) (model.Namespace, error) {
	return model.Namespace{
		Id:        from.Id,
		Name:      from.Name,
		CreatedAt: from.CreatedAt,
		UpdatedAt: from.UpdatedAt,
	}, nil
}
