package infrastructure

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/internal/entities"
)

// PGProductRepository — реализация ProductRepository для PostgreSQL (Squirrel, без ORM)
type PGProductRepository struct {
	db *sql.DB
	qb squirrel.StatementBuilderType
}

// NewPGProductRepository создаёт новый PGProductRepository
func NewPGProductRepository(db *sql.DB) *PGProductRepository {
	return &PGProductRepository{
		db: db,
		qb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// Save сохраняет (insert) товар
func (r *PGProductRepository) Save(ctx context.Context, p entities.Product) (entities.Product, error) {
	q := r.qb.Insert("product").
		Columns("id", "reception_id", "type", "received_at").
		Values(p.ID, p.ReceptionID, p.Type, p.ReceivedAt).
		Suffix("RETURNING id")
	row := q.RunWith(r.db).QueryRowContext(ctx)
	var id uuid.UUID
	if err := row.Scan(&id); err != nil {
		return entities.Product{}, err
	}
	p.ID = id
	return p, nil
}

// DeleteLast удаляет последний добавленный товар по приёмке (LIFO)
func (r *PGProductRepository) DeleteLast(ctx context.Context, receptionID uuid.UUID) (*entities.Product, error) {
	q := r.qb.Select("id", "reception_id", "type", "received_at").
		From("product").
		Where(squirrel.Eq{"reception_id": receptionID}).
		OrderBy("received_at DESC").
		Limit(1)
	row := q.RunWith(r.db).QueryRowContext(ctx)
	var p entities.Product
	var typ string
	if err := row.Scan(&p.ID, &p.ReceptionID, &typ, &p.ReceivedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	p.Type = entities.ProductType(typ)
	// Удаляем найденный товар
	delQ := r.qb.Delete("product").Where(squirrel.Eq{"id": p.ID})
	_, err := delQ.RunWith(r.db).ExecContext(ctx)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// ListByReception возвращает все товары по приёмке
func (r *PGProductRepository) ListByReception(ctx context.Context, receptionID uuid.UUID) ([]entities.Product, error) {
	q := r.qb.Select("id", "reception_id", "type", "received_at").
		From("product").
		Where(squirrel.Eq{"reception_id": receptionID}).
		OrderBy("received_at ASC")
	rows, err := q.RunWith(r.db).QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []entities.Product
	for rows.Next() {
		var p entities.Product
		var typ string
		if err := rows.Scan(&p.ID, &p.ReceptionID, &typ, &p.ReceivedAt); err != nil {
			return nil, err
		}
		p.Type = entities.ProductType(typ)
		res = append(res, p)
	}
	return res, rows.Err()
}
