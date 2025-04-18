package postgreSQL

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"pvZ/internal/adapters/db"
	"pvZ/internal/domain/models"
)

type receptionRepositoryImpl struct {
	db        *sqlx.DB
	statement squirrel.StatementBuilderType
}

func NewReceptionRepository(db *sqlx.DB) db.ReceptionRepository {
	return &receptionRepositoryImpl{
		db:        db,
		statement: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *receptionRepositoryImpl) Create(ctx context.Context, pvzID string) (*models.Reception, error) {
	query, args, err := r.statement.
		Insert("reception").
		Columns("pvz_id", "date_time", "status").
		Values(pvzID, squirrel.Expr("NOW()"), "in_progress").
		Suffix("RETURNING id, pvz_id, date_time, status").
		ToSql()

	if err != nil {
		return nil, err
	}

	var reception models.Reception
	err = r.db.GetContext(ctx, &reception, query, args...)
	return &reception, err
}

func (r *receptionRepositoryImpl) GetLastInProgress(ctx context.Context, pvzID string) (*models.Reception, error) {
	query, args, err := r.statement.
		Select("id", "pvz_id", "date_time", "status").
		From("reception").
		Where(squirrel.Eq{
			"pvz_id": pvzID,
			"status": "in_progress",
		}).
		OrderBy("date_time DESC").
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	var reception models.Reception
	err = r.db.GetContext(ctx, &reception, query, args...)
	if err != nil {
		return nil, err
	}
	return &reception, nil
}

func (r *receptionRepositoryImpl) CloseLastReception(ctx context.Context, pvzID string) (*models.Reception, error) {
	reception, err := r.GetLastInProgress(ctx, pvzID)
	if err != nil {
		return nil, err
	}

	query, args, err := r.statement.
		Update("reception").
		Set("status", "close").
		Where(squirrel.Eq{"id": reception.ID}).
		Suffix("RETURNING id, pvz_id, date_time, status").
		ToSql()
	if err != nil {
		return nil, err
	}

	var updated models.Reception
	err = r.db.GetContext(ctx, &updated, query, args...)
	return &updated, err
}
