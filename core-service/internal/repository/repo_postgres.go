package repository

import (
	"cisterna-mvp/core-service/internal/domain"
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	_ "github.com/lib/pq"
)

//go:embed schema.sql
var schemaSQL string

type PositionRepository interface {
	SavePosition(ctx context.Context, pos domain.TruckPosition) error
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(connStr string) (PositionRepository, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	repo := &postgresRepo{
		db: db,
	}

	if err = repo.createTruckCoordinates(context.Background()); err != nil {
		return nil, fmt.Errorf("error to create table: %v", err)
	}

	return repo, nil
}

func (r *postgresRepo) createTruckCoordinates(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, schemaSQL)
	return err
}

func (r *postgresRepo) SavePosition(ctx context.Context, pos domain.TruckPosition) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error to inicialize transation: %v", err)
	}
	defer tx.Rollback()

	// -- INSERT
	queryHistory := `
		INSERT INTO truck_coordinates (truck_id, location, recorded_at)
		VALUES ($1, ST_SetSRID(ST_MakePoint($2, $3), 4326), $4)
	`
	_, err = tx.ExecContext(ctx, queryHistory,
		pos.TruckID,
		pos.Longitude, // vem primeiro
		pos.Latitude,
		pos.Timestamp)
	if err != nil {
		return fmt.Errorf("error to INSERT on history: %v", err)
	}

	// -- UPSERT
	queryCurrentStatus := `
		INSERT INTO truck_current_status (truck_id, location, last_seen)
		VALUES ($1, ST_SetSRID(ST_MakePoint($2, $3), 4326), $4)
		ON CONFLICT (truck_id) 
		DO UPDATE SET 
			location = EXCLUDED.location,
			last_seen = EXCLUDED.last_seen;
	`
	_, err = tx.ExecContext(ctx, queryCurrentStatus,
		pos.TruckID,
		pos.Longitude,
		pos.Latitude,
		pos.Timestamp)
	if err != nil {
		return fmt.Errorf("error to UPSERT on current status: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error to commit trasation: %v", err)
	}

	return nil
}
