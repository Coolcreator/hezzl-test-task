package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"goods-service/internal/good/domain"
)

type GoodStorage struct {
	pool *pgxpool.Pool
}

func (s *GoodStorage) CreateGood(ctx context.Context, createGood domain.CreateGood) (good domain.Good, err error) {
	const query = `INSERT INTO goods (project_id, name) VALUES ($1, $2) RETURNING id, project_id, name, description, priority, removed, created_at;`
	row := s.pool.QueryRow(ctx, query, createGood.ProjectID, createGood.Name)
	err = row.Scan(&good.ID, &good.ProjectID, &good.Name, &good.Description, &good.Priority, &good.Removed, &good.CreatedAt)
	if err != nil {
		err = fmt.Errorf("insert query: %w", err)
		return
	}
	return
}

func (s *GoodStorage) UpdateGood(ctx context.Context, updateGood domain.UpdateGood) (good domain.Good, err error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		err = fmt.Errorf("begin tx: %w", err)
		return
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				rollbackErr = fmt.Errorf("tx rollback: %w", rollbackErr)
				err = errors.Join(err, rollbackErr)
			}
		}
	}()
	const selectQuery = `SELECT TRUE FROM goods WHERE id = $1 AND project_id = $2 FOR UPDATE;`
	var exists bool
	err = tx.QueryRow(ctx, selectQuery, updateGood.ID, updateGood.ProjectID).Scan(&exists)
	if err != nil {
		err = fmt.Errorf("check existance: %w", err)
		return
	}
	const updateQuery = `UPDATE goods SET name = $1, description = $2 WHERE id = $3 AND project_id = $4 RETURNING id, project_id, name, description, priority, removed, created_at;`
	row := tx.QueryRow(ctx, updateQuery, updateGood.Name, updateGood.Description, updateGood.ID, updateGood.ProjectID)
	err = row.Scan(&good.ID, &good.ProjectID, &good.Name, &good.Description, &good.Priority, &good.Removed, &good.CreatedAt)
	if err != nil {
		err = fmt.Errorf("update good: %w", err)
		return
	}

	eventID := uuid.New()
	const insertQuery = `INSERT INTO outbox(event_id, good_id, project_id) VALUES ($1, $2, $3);`
	_, err = tx.Exec(ctx, insertQuery, eventID, updateGood.ID, updateGood.ProjectID)
	if err != nil {
		err = fmt.Errorf("insert event: %w", err)
		return
	}
	err = tx.Commit(ctx)
	if err != nil {
		err = fmt.Errorf("tx commit: %w", err)
		return
	}
	return
}

func (s *GoodStorage) DeleteGood(ctx context.Context, deleteGood domain.DeleteGood) (err error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		err = fmt.Errorf("begin tx: %w", err)
		return
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				rollbackErr = fmt.Errorf("tx rollback: %w", rollbackErr)
				err = errors.Join(err, rollbackErr)
			}
		}
	}()
	const existanceQuery = `SELECT TRUE FROM goods WHERE id = $1 AND project_id = $2 FOR UPDATE;`
	var exists bool
	err = tx.QueryRow(ctx, existanceQuery, deleteGood.ID, deleteGood.ProjectID).Scan(&exists)
	if err != nil {
		err = fmt.Errorf("check existance: %w", err)
		return
	}
	const deleteQuery = `UPDATE goods SET removed = $1 WHERE id = $2 AND project_id = $3;`
	_, err = tx.Exec(ctx, deleteQuery, true, deleteGood.ID, deleteGood.ProjectID)
	if err != nil {
		err = fmt.Errorf("delete good: %w", err)
		return
	}
	eventID := uuid.New()
	const insertQuery = `INSERT INTO outbox(event_id, good_id, project_id) VALUES ($1, $2, $3);`
	_, err = tx.Exec(ctx, insertQuery, eventID, deleteGood.ID, deleteGood.ProjectID)
	if err != nil {
		err = fmt.Errorf("insert event: %w", err)
		return
	}
	err = tx.Commit(ctx)
	if err != nil {
		err = fmt.Errorf("tx commit: %w", err)
		return
	}
	return
}

func (s *GoodStorage) ListGoods(ctx context.Context, listGoods domain.ListGoods) (goodsList domain.GoodsList, err error) {
	const selectQuery = `SELECT id, project_id, name, description, priority, removed, created_at FROM goods LIMIT $1 OFFSET $2;`
	rows, err := s.pool.Query(ctx, selectQuery, listGoods.Limit, listGoods.Offset)
	if err != nil {
		err = fmt.Errorf("get goods list: %w", err)
		return
	}
	defer rows.Close()
	goods := make([]domain.Good, listGoods.Limit)
	meta := domain.Meta{}
	for rows.Next() {
		var good domain.Good
		err = rows.Scan(&good.ID, &good.ProjectID, &good.Name, &good.Description, &good.Priority, &good.Removed, &good.CreatedAt)
		if err != nil {
			err = fmt.Errorf("rows scan: %w", err)
			return
		}
		meta.Total++
		if good.Removed {
			meta.Removed++
		}
	}
	if err = rows.Err(); err != nil {
		err = fmt.Errorf("rows error: %w", err)
		return
	}
	goodsList = domain.GoodsList{
		Meta:  meta,
		Goods: goods,
	}
	return
}

func (s *GoodStorage) ReprioritizeGood(ctx context.Context, reprioritizeGood domain.ReprioritizeGood) (
	goodsPriorities []domain.GoodPriority, err error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		err = fmt.Errorf("begin tx: %w", err)
		return
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				rollbackErr = fmt.Errorf("tx rollback: %w", rollbackErr)
				err = errors.Join(err, rollbackErr)
			}
		}
	}()
	const existanceQuery = `SELECT TRUE FROM goods WHERE id = $1 AND project_id = $2 FOR UPDATE;`
	var exists bool
	err = tx.QueryRow(ctx, existanceQuery, reprioritizeGood.ID, reprioritizeGood.ProjectID).Scan(&exists)
	if err != nil {
		err = fmt.Errorf("check existance: %w", err)
		return
	}
	const reprioritizeQuery = `UPDATE goods SET priority = $1 WHERE id = $2 AND project_id = $3 RETURNING id, priority;`
	rows, err := tx.Query(ctx, reprioritizeQuery, true, reprioritizeGood.ID, reprioritizeGood.ProjectID)
	if err != nil {
		err = fmt.Errorf("reprioritize good: %w", err)
		return
	}
	goodsPriorities = make([]domain.GoodPriority, 0)
	for rows.Next() {
		var goodPriority domain.GoodPriority
		err = rows.Scan(&goodPriority.ID, &goodPriority.Priority)
		if err != nil {
			err = fmt.Errorf("rows scan: %w", err)
			return
		}
		goodsPriorities = append(goodsPriorities, goodPriority)
	}
	if err = rows.Err(); err != nil {
		err = fmt.Errorf("rows error: %w", err)
		return
	}
	eventID := uuid.New()
	const insertQuery = `INSERT INTO outbox(event_id, good_id, project_id) VALUES ($1, $2, $3);`
	_, err = tx.Exec(ctx, insertQuery, eventID, reprioritizeGood.ID, reprioritizeGood.ProjectID)
	if err != nil {
		err = fmt.Errorf("insert event: %w", err)
		return
	}
	err = tx.Commit(ctx)
	if err != nil {
		err = fmt.Errorf("tx commit: %w", err)
		return
	}
	return
}

func NewGoodStorage(pool *pgxpool.Pool) (storage *GoodStorage) {
	storage = &GoodStorage{
		pool: pool,
	}
	return
}
