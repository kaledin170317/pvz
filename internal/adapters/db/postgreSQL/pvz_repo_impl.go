package postgreSQL

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"pvZ/internal/adapters/db"
	"pvZ/internal/domain/models"
	"time"
)

type pvzRepositoryImpl struct {
	db        *sqlx.DB
	statement sq.StatementBuilderType
}

func NewPVZRepository(db *sqlx.DB) db.PVZRepository {
	return &pvzRepositoryImpl{
		db:        db,
		statement: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *pvzRepositoryImpl) Create(ctx context.Context, pvz *models.Pvz) (*models.Pvz, error) {
	query, args, err := r.statement.
		Insert("pvz").
		Columns("city", "registration_date").
		Values(pvz.City, sq.Expr("NOW()")).
		Suffix("RETURNING id, registration_date, city").
		ToSql()
	if err != nil {
		return nil, err
	}

	var created models.Pvz
	err = r.db.GetContext(ctx, &created, query, args...)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			return nil, err
		}
	}

	return &created, nil
}

func (r *pvzRepositoryImpl) GetByID(ctx context.Context, id string) (*models.Pvz, error) {
	query, args, err := r.statement.
		Select("id", "registration_date", "city").
		From("pvz").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var pvz models.Pvz
	err = r.db.GetContext(ctx, &pvz, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &pvz, nil
}

func (r *pvzRepositoryImpl) ListWithDateRange(ctx context.Context, startDate, endDate *time.Time, limit, offset int) ([]models.Pvz, error) {
	builder := r.statement.
		Select("id", "registration_date", "city").
		From("pvz").
		OrderBy("registration_date DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	if startDate != nil {
		builder = builder.Where(sq.GtOrEq{"registration_date": *startDate})
	}
	if endDate != nil {
		builder = builder.Where(sq.LtOrEq{"registration_date": *endDate})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var result []models.Pvz
	err = r.db.SelectContext(ctx, &result, query, args...)
	return result, err
}
