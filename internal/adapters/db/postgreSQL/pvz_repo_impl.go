package postgreSQL

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"pvZ/internal/adapters/db"
	"pvZ/internal/domain/models"
	"pvZ/internal/logger"
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
		Columns("id", "city", "registration_date").
		Values(pvz.ID, pvz.City, pvz.RegistrationDate).
		Suffix("RETURNING id, registration_date, city").
		ToSql()
	if err != nil {
		return nil, err
	}

	var created models.Pvz
	err = r.db.GetContext(ctx, &created, query, args...)
	if err != nil {
		return nil, err
	}
	logger.Log.Info("pvz created", "id", created.ID, "city", created.City)
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

func (r *pvzRepositoryImpl) GetReceptionsByPVZ(ctx context.Context, pvzID string) ([]models.Reception, error) {
	query, args, err := r.statement.Select("id", "pvz_id", "date_time", "status").
		From("reception").Where(sq.Eq{"pvz_id": pvzID}).OrderBy("date_time").ToSql()
	if err != nil {
		return nil, err
	}
	var result []models.Reception
	err = r.db.SelectContext(ctx, &result, query, args...)
	return result, err
}

func (r *pvzRepositoryImpl) GetProductsByReception(ctx context.Context, receptionID string) ([]models.Product, error) {
	query, args, err := r.statement.Select("id", "reception_id", "type", "date_time").
		From("product").Where(sq.Eq{"reception_id": receptionID}).ToSql()
	if err != nil {
		return nil, err
	}
	var result []models.Product
	err = r.db.SelectContext(ctx, &result, query, args...)
	return result, err
}
