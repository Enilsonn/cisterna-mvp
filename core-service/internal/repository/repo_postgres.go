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
	CreateCistern(ctx context.Context, cis domain.Cisterna) (int64, error)
	SavePosition(ctx context.Context, pos domain.TruckPosition) error
	GetCisterns(ctx context.Context) ([]domain.Cisterna, error)
	GetTruckCurrrentLocation(ctx context.Context, id string) (*domain.TruckStatus, error)
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

func (r *postgresRepo) GetCisterns(ctx context.Context) ([]domain.Cisterna, error) {
	query := `
			SELECT id, nome, capacity_liters, ST_Y(location) as lat, ST_X(location) as lon
			FROM cisterns`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error to get cisterns: %v", err)
	}
	defer rows.Close()

	var cisterns []domain.Cisterna
	for rows.Next() {
		var c domain.Cisterna
		if err := rows.Scan(
			&c.ID,
			&c.Nome,
			&c.CapacityLiters,
			&c.Latitude,
			&c.Longitude,
		); err != nil {
			return nil, fmt.Errorf("error to scan cistern")
		}
		cisterns = append(cisterns, c)
	}

	return cisterns, nil
}

func (r *postgresRepo) GetTruckCurrrentLocation(ctx context.Context, id string) (*domain.TruckStatus, error) {
	query := `
			SELECT truck_id, ST_Y(location) as lat, ST_X(location) as lon, last_seen
			FROM truck_current_status
			WHERE truck_id = $1
			`
	var t domain.TruckStatus
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.TruckID,
		&t.Latitude,
		&t.Longitude,
		&t.LastSeen,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("truck not found: %v", err)
		} else {
			return nil, fmt.Errorf("error to search truck: %v", err)
		}
	}

	return &t, nil
}

func (r *postgresRepo) CreateCistern(ctx context.Context, cis domain.Cisterna) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("error to initialize transation: %v", err)
	}
	defer tx.Rollback()

	var id int64
	query := `
		INSERT INTO cisterns (name, capacity_liters, location)
		VALUES ($1, $2, ST_SetSRID(ST_MakePoint($3, $4), 4326))
		RETURNING id
	`
	if err = tx.QueryRowContext(ctx, query,
		cis.Nome,
		cis.CapacityLiters,
		cis.Longitude,
		cis.Latitude).Scan(&id); err != nil {
		return 0, fmt.Errorf("error to save cistern: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("error to commit: %v", err)
	}

	return id, nil
}
