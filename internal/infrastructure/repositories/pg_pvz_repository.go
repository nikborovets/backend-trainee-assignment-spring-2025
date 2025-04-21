package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
)

// PGPVZRepository — реализация PVZRepository для PostgreSQL (Squirrel, без ORM)
type PGPVZRepository struct {
	db *sql.DB
	qb squirrel.StatementBuilderType
}

// NewPGPVZRepository создаёт новый PGPVZRepository
func NewPGPVZRepository(db *sql.DB) *PGPVZRepository {
	return &PGPVZRepository{
		db: db,
		qb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// Save сохраняет (insert) PVZ
func (r *PGPVZRepository) Save(ctx context.Context, pvz entities.PVZ) (entities.PVZ, error) {
	q := r.qb.Insert("pvz").
		Columns("id", "registration_date", "city").
		Values(pvz.ID, pvz.RegistrationDate, pvz.City).
		Suffix("ON CONFLICT (id) DO UPDATE SET registration_date = EXCLUDED.registration_date, city = EXCLUDED.city RETURNING id")
	row := q.RunWith(r.db).QueryRowContext(ctx)
	var id uuid.UUID
	if err := row.Scan(&id); err != nil {
		return entities.PVZ{}, err
	}
	pvz.ID = id
	return pvz, nil
}

// List возвращает список PVZ с фильтрами по дате и пагинацией
func (r *PGPVZRepository) List(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]entities.PVZ, error) {
	q := r.qb.Select("id", "registration_date", "city").From("pvz")
	if startDate != nil {
		q = q.Where(squirrel.GtOrEq{"registration_date": *startDate})
	}
	if endDate != nil {
		q = q.Where(squirrel.LtOrEq{"registration_date": *endDate})
	}
	if limit > 0 {
		q = q.Limit(uint64(limit))
	}
	if page > 0 && limit > 0 {
		offset := uint64((page - 1) * limit)
		q = q.Offset(offset)
	}
	rows, err := q.RunWith(r.db).QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []entities.PVZ
	for rows.Next() {
		var pvz entities.PVZ
		if err := rows.Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City); err != nil {
			return nil, err
		}
		res = append(res, pvz)
	}
	return res, rows.Err()
}
