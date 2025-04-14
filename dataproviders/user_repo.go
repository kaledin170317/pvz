package dataproviders

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"pvZ/dataproviders/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.UserModel) (*models.UserModel, error)
	GetByEmail(ctx context.Context, email string) (*models.UserModel, error)
}

type userRepositoryImpl struct {
	db        *sqlx.DB
	statement squirrel.StatementBuilderType
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepositoryImpl{
		db:        db,
		statement: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// ✅ Создание пользователя
func (r *userRepositoryImpl) Create(ctx context.Context, user *models.UserModel) (*models.UserModel, error) {
	query, args, err := r.statement.
		Insert("users").
		Columns("email", "password", "role").
		Values(user.Email, user.Password, user.Role).
		Suffix("RETURNING id, email, role").
		ToSql()
	if err != nil {
		return nil, err
	}

	var created models.UserModel
	err = r.db.GetContext(ctx, &created, query, args...)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

// ✅ Поиск по email
func (r *userRepositoryImpl) GetByEmail(ctx context.Context, email string) (*models.UserModel, error) {
	query, args, err := r.statement.
		Select("id", "email", "password", "role").
		From("users").
		Where(squirrel.Eq{"email": email}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var user models.UserModel
	err = r.db.GetContext(ctx, &user, query, args...)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
