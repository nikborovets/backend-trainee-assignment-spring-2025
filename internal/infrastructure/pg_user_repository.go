package infrastructure

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
)

// PGUserRepository — реализация UserRepository для PostgreSQL (Squirrel, без ORM)
type PGUserRepository struct {
	db *sql.DB
	qb squirrel.StatementBuilderType
}

// NewPGUserRepository создаёт новый PGUserRepository
func NewPGUserRepository(db *sql.DB) *PGUserRepository {
	return &PGUserRepository{
		db: db,
		qb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// Create сохраняет пользователя и хэш пароля
func (r *PGUserRepository) Create(ctx context.Context, user entities.User, passwordHash string) (entities.User, error) {
	q := r.qb.Insert("users").
		Columns("id", "email", "role", "registration_date", "password_hash").
		Values(user.ID, user.Email, user.Role, user.RegistrationDate, passwordHash).
		Suffix("RETURNING id")
	row := q.RunWith(r.db).QueryRowContext(ctx)
	var id uuid.UUID
	if err := row.Scan(&id); err != nil {
		return entities.User{}, err
	}
	user.ID = id
	return user, nil
}

// GetByEmail ищет пользователя по email, возвращает User и passwordHash
func (r *PGUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, string, error) {
	q := r.qb.Select("id", "email", "role", "registration_date", "password_hash").
		From("users").
		Where(squirrel.Eq{"email": email})
	row := q.RunWith(r.db).QueryRowContext(ctx)
	var user entities.User
	var hash string
	if err := row.Scan(&user.ID, &user.Email, &user.Role, &user.RegistrationDate, &hash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", nil
		}
		return nil, "", err
	}
	return &user, hash, nil
}
