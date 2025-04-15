package postgreSQL

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"pvZ/internal/adapters/db"
	"pvZ/internal/domain/models"
)

type userRepositoryImpl struct {
	db        *sqlx.DB
	statement squirrel.StatementBuilderType
}

func NewUserRepository(db *sqlx.DB) db.UserRepository {
	return &userRepositoryImpl{
		db:        db,
		statement: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *models.User) (*models.User, error) {
	query, args, err := r.statement.
		Insert("users").
		Columns("email", "password", "role").
		Values(user.Email, user.Password, user.Role).
		Suffix("RETURNING id, email, role").
		ToSql()
	if err != nil {
		return nil, err
	}

	var created models.User
	err = r.db.GetContext(ctx, &created, query, args...)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *userRepositoryImpl) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query, args, err := r.statement.
		Select("id", "email", "password", "role").
		From("users").
		Where(squirrel.Eq{"email": email}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var user models.User
	err = r.db.GetContext(ctx, &user, query, args...)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
