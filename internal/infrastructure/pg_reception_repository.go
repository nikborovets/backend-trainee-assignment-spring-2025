package infrastructure

import (
	"context"
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
)

// PGReceptionRepository — реализация ReceptionRepository для PostgreSQL (Squirrel, без ORM)
type PGReceptionRepository struct {
	db *sql.DB
	qb squirrel.StatementBuilderType
}

// NewPGReceptionRepository создаёт новый PGReceptionRepository
func NewPGReceptionRepository(db *sql.DB) *PGReceptionRepository {
	return &PGReceptionRepository{
		db: db,
		qb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// Save сохраняет (insert/update) приёмку
func (r *PGReceptionRepository) Save(ctx context.Context, rec entities.Reception) (entities.Reception, error) {
	q := r.qb.Insert("reception").
		Columns("id", "pvz_id", "status", "date_time").
		Values(rec.ID, rec.PVZID, rec.Status, rec.DateTime).
		Suffix("ON CONFLICT (id) DO UPDATE SET pvz_id = EXCLUDED.pvz_id, status = EXCLUDED.status, date_time = EXCLUDED.date_time RETURNING id")
	row := q.RunWith(r.db).QueryRowContext(ctx)
	var id uuid.UUID
	if err := row.Scan(&id); err != nil {
		return entities.Reception{}, err
	}
	rec.ID = id
	return rec, nil
}

// GetActive возвращает открытую приёмку по PVZ (status = in_progress)
func (r *PGReceptionRepository) GetActive(ctx context.Context, pvzID uuid.UUID) (*entities.Reception, error) {
	q := r.qb.Select("id", "pvz_id", "status", "date_time").
		From("reception").
		Where(squirrel.Eq{"pvz_id": pvzID, "status": entities.ReceptionInProgress})
	row := q.RunWith(r.db).QueryRowContext(ctx)
	var rec entities.Reception
	var status string
	if err := row.Scan(&rec.ID, &rec.PVZID, &status, &rec.DateTime); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	rec.Status = entities.ReceptionStatus(status)
	return &rec, nil
}

// CloseLast закрывает последнюю открытую приёмку по PVZ (status = in_progress → close, date_time обновляется)
func (r *PGReceptionRepository) CloseLast(ctx context.Context, pvzID uuid.UUID, closedAt time.Time) error {
	q := r.qb.Update("reception").
		Set("status", entities.ReceptionClosed).
		Set("date_time", closedAt).
		Where(squirrel.Eq{"pvz_id": pvzID, "status": entities.ReceptionInProgress})
	res, err := q.RunWith(r.db).ExecContext(ctx)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
