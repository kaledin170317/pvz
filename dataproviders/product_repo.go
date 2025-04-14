package dataproviders

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"pvZ/dataproviders/models"
)

type ProductRepository interface {
	AddProduct(ctx context.Context, receptionID string, productType string) (*models.ProductModel, error)
	DeleteLastProduct(ctx context.Context, receptionID string) error
	GetLastInReception(ctx context.Context, receptionID string) (*models.ProductModel, error)
}

type productRepositoryImpl struct {
	db        *sqlx.DB
	statement squirrel.StatementBuilderType
}

func NewProductRepository(db *sqlx.DB) ProductRepository {
	return &productRepositoryImpl{
		db:        db,
		statement: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *productRepositoryImpl) AddProduct(ctx context.Context, receptionID string, productType string) (*models.ProductModel, error) {
	query, args, err := r.statement.
		Insert("product").
		Columns("reception_id", "type", "date_time").
		Values(receptionID, productType, squirrel.Expr("NOW()")).
		Suffix("RETURNING id, reception_id, type, date_time").
		ToSql()
	if err != nil {
		return nil, err
	}

	var product models.ProductModel
	err = r.db.GetContext(ctx, &product, query, args...)
	return &product, err
}

func (r *productRepositoryImpl) GetLastInReception(ctx context.Context, receptionID string) (*models.ProductModel, error) {
	query, args, err := r.statement.
		Select("id", "reception_id", "type", "date_time").
		From("product").
		Where(squirrel.Eq{"reception_id": receptionID}).
		OrderBy("date_time DESC").
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	var product models.ProductModel
	err = r.db.GetContext(ctx, &product, query, args...)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepositoryImpl) DeleteLastProduct(ctx context.Context, receptionID string) error {
	last, err := r.GetLastInReception(ctx, receptionID)
	if err != nil {
		return err
	}

	query, args, err := r.statement.
		Delete("product").
		Where(squirrel.Eq{"id": last.ID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}
