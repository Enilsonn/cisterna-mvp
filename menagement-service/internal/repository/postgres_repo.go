package repository

import (
	"cisterna-mvp/menagement-service/internal/domain"
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	_ "embed"
)

type SighRepository interface {
	CreatePipeiro(ctx context.Context, pipeiro domain.Pipeiro) (string, error)
	CreateTruck(ctx context.Context, truck domain.Truck) (string, error)
	CreateCistern(ctx context.Context, cistern domain.Cistern) (string, error)
	CreateDelivery(ctx context.Context, delivery domain.Delivery) (string, error)

	UpdatePipeiro(ctx context.Context, pipeiro domain.Pipeiro) error
	UpdateTruck(ctx context.Context, truck domain.Truck) error
	UpdateCistern(ctx context.Context, cistern domain.Cistern) error
	UpdateDelivery(ctx context.Context, delivery domain.Delivery) error

	GetPipeiroByCPF(ctx context.Context, cpf string) (*domain.Pipeiro, error)
	GetTruckByPlate(ctx context.Context, palte string) (*domain.Truck, error)
	GetCisternByUUID(ctx context.Context, uuid string) (*domain.Cistern, error)
	GetDeliveryByUUID(ctx context.Context, uuid string) (*domain.Delivery, error)

	GetTruckByPipeiroUUID(ctx context.Context, uuid string) ([]*domain.Truck, error)
	GetCisterns(ctx context.Context) ([]*domain.Cistern, error)
	GetDeliveryByPipeiroUUID(ctx context.Context, uuid string) ([]*domain.Delivery, error)
	GetDeliveryByTruckUUID(ctx context.Context, uuid string) ([]*domain.Delivery, error)
}

type postgresRepo struct {
	repo *sql.DB
}

func NewPostgresRepo(db *sql.DB) SighRepository {
	return &postgresRepo{
		repo: db,
	}
}

func (r *postgresRepo) CreatePipeiro(ctx context.Context, pipeiro domain.Pipeiro) (string, error) {
	tx, err := r.repo.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("error to start a transaction: %v", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO pipeiros (name, cpf, cnh, phone)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id string
	if err := tx.QueryRowContext(ctx, query,
		pipeiro.ID,
		pipeiro.CPF,
		pipeiro.CNH,
		pipeiro.Phone).Scan(&id); err != nil {
		return "", fmt.Errorf("error to create pipeiro: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("error to commit pipeiro insert: %v", err)
	}

	return id, nil
}

func (r *postgresRepo) CreateTruck(ctx context.Context, truck domain.Truck) (string, error) {
	tx, err := r.repo.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("error to start a trasaction: %v", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO trucks (plate, capacity_liters, pipeiro_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var id string

	if err := tx.QueryRowContext(ctx, query,
		truck.Plate,
		truck.CapacityLiters,
		truck.PipeiroID).Scan(&id); err != nil {
		return "", fmt.Errorf("error to create truck")
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("error to commit truck insert: %v", err)
	}

	return id, nil

}

func (r *postgresRepo) CreateCistern(ctx context.Context, cistern domain.Cistern) (string, error) {
	tx, err := r.repo.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("error to start a trasaction: %v", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO cisterns (name, responsible_name, city, capacity_liters, latitude, longitude)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id string
	if err := tx.QueryRowContext(ctx, query,
		cistern.Name,
		cistern.ResponsabibleName,
		cistern.City,
		cistern.CapacityLiters,
		cistern.Latitude,
		cistern.Longitude,
	).Scan(&id); err != nil {
		return "", fmt.Errorf("error to create cistern: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("error to commit cistern insert: %v", err)
	}

	return id, nil
}

func (r *postgresRepo) CreateDelivery(ctx context.Context, delivery domain.Delivery) (string, error) {
	tx, err := r.repo.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("error to start a trasaction: %v", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO deliveries (cistern_id, truck_id, scheduled_date, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id string
	if err := tx.QueryRowContext(ctx, query,
		delivery.CisternID,
		delivery.TruckID,
		delivery.ScheduledDate,
		delivery.Status,
	).Scan(&id); err != nil {
		return "", fmt.Errorf("error to create delivery: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("error to commit delivery insert: %v", err)
	}

	return id, nil
}

func (r *postgresRepo) UpdatePipeiro(ctx context.Context, pipeiro domain.Pipeiro) error {
	tx, err := r.repo.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error to start a trasaction: %v", err)
	}
	defer tx.Rollback()

	query := `
		UPDATE pipeiros
		SET name = $1, cpf = $2, cnh = $3, phone = $4, is_active = $5
		WHERE id = $6
	`

	result, err := tx.ExecContext(ctx, query,
		pipeiro.Name,
		pipeiro.CPF,
		pipeiro.CNH,
		pipeiro.Phone,
		pipeiro.IsActive,
		pipeiro.ID,
	)
	if err != nil {
		return fmt.Errorf("error to update pipeiro: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error to vizualize affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no one row affected")
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error to commit pipeiro update: %v", err)
	}

	return nil
}

func (r *postgresRepo) UpdateTruck(ctx context.Context, truck domain.Truck) error {
	tx, err := r.repo.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error to start a trasaction: %v", err)
	}
	defer tx.Rollback()

	query := `
		UPDATE trucks
		SET plate = $1, capacity_liters = $2, pipeiro_id = $3
		WHERE id = $4
	`

	result, err := tx.ExecContext(ctx, query,
		truck.Plate,
		truck.CapacityLiters,
		truck.PipeiroID,
		truck.ID,
	)
	if err != nil {
		return fmt.Errorf("error to update truck: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error to vizualize rows affected")
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no one row affected")
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error to commit truck update: %v", err)
	}

	return nil
}

func (r *postgresRepo) UpdateCistern(ctx context.Context, cistern domain.Cistern) error {
	tx, err := r.repo.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error to start a trasaction: %v", err)
	}
	defer tx.Rollback()

	query := `
		UPDATE cisterns
		SET name = $1, responsible_name = $2, city = $3, capacity_liters = $4, latitude = $5, longitude = $6
		WHERE id = $7
	`

	result, err := tx.ExecContext(ctx, query,
		cistern.Name,
		cistern.ResponsabibleName,
		cistern.City,
		cistern.CapacityLiters,
		cistern.Latitude,
		cistern.Longitude,
		cistern.ID,
	)

	if err != nil {
		return fmt.Errorf("error to update cistern: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error to vizualize rows affected")
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no one row affected")
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error to commit cisterna update: %v", err)
	}

	return nil
}

func (r *postgresRepo) UpdateDelivery(ctx context.Context, delivery domain.Delivery) error {
	tx, err := r.repo.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error to start a trasaction: %v", err)
	}
	defer tx.Rollback()

	query := `
		UPDATE deliveries
		SET cistern_id = $1, truck_id = $2, scheduled_date = $3, status = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
	`

	result, err := tx.ExecContext(ctx, query,
		delivery.CisternID,
		delivery.TruckID,
		delivery.UpdatedAt,
		delivery.Status,
		delivery.ID,
	)
	if err != nil {
		return fmt.Errorf("error to update delivery: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error to vizualize rows affected")
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no one row affected")
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error to commit delivery update: %v", err)
	}

	return nil
}

func (r *postgresRepo) GetPipeiroByCPF(ctx context.Context, cpf string) (*domain.Pipeiro, error) {
	query := `
		SELECT id, name, cpf, cnh, phone, is_active, created_at 
		FROM pipeiros 
		WHERE cpf = $1
	`
	var pipeiro domain.Pipeiro
	if err := r.repo.QueryRowContext(ctx, query, cpf).Scan(
		&pipeiro.ID,
		&pipeiro.Name,
		&pipeiro.CPF,
		&pipeiro.CNH,
		&pipeiro.Phone,
		&pipeiro.IsActive,
		&pipeiro.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pipeiro(cpf: %v) not found: %v", cpf, err)
		}
		return nil, fmt.Errorf("internal error to search pipeiro(%v): %v", cpf, err)
	}

	return &pipeiro, nil
}

func (r *postgresRepo) GetTruckByPlate(ctx context.Context, plate string) (*domain.Truck, error) {
	query := `
		SELECT id, plate, capacity_liters, pipeiro_id, created_at 
		FROM trucks 
		WHERE license_plate = $1
	`

	var truck domain.Truck
	if err := r.repo.QueryRowContext(ctx, query, plate).Scan(
		&truck.ID,
		&truck.Plate,
		&truck.CapacityLiters,
		&truck.PipeiroID,
		&truck.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("truck(plate: %v) not found: %v", plate, err)
		}
		return nil, fmt.Errorf("internal error to search truck(plate:%v): %v", plate, err)
	}

	return &truck, nil
}

func (r *postgresRepo) GetCisternByUUID(ctx context.Context, uuid string) (*domain.Cistern, error) {
	query := `
		SELECT id, name, responsible_name, city, capacity_liters, latitude, longitude, created_at 
		FROM cisterns 
		WHERE id = $1
	`

	var cistern domain.Cistern
	if err := r.repo.QueryRowContext(ctx, query, uuid).Scan(
		&cistern.ID,
		&cistern.Name,
		&cistern.City,
		&cistern.CapacityLiters,
		&cistern.Latitude,
		&cistern.Longitude,
		&cistern.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cistern not found: %v", err)
		}
		return nil, fmt.Errorf("internal error to search cistern: %v", err)
	}

	return &cistern, nil
}

func (r *postgresRepo) GetDeliveryByUUID(ctx context.Context, uuid string) (*domain.Delivery, error) {
	query := `
		SELECT id, cistern_id, truck_id, scheduled_date, status, created_at, updated_at 
		FROM deliveries 
		WHERE id = $1
	`

	var delivery domain.Delivery
	if err := r.repo.QueryRowContext(ctx, query, uuid).Scan(
		&delivery.ID,
		&delivery.CisternID,
		&delivery.TruckID,
		&delivery.ScheduledDate,
		&delivery.Status,
		&delivery.CreatedAt,
		&delivery.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("delivery not found: %v", err)
		}
		return nil, fmt.Errorf("internal error to search delivery: %v", err)
	}

	return &delivery, nil
}

func (r *postgresRepo) GetTruckByPipeiroUUID(ctx context.Context, uuid string) ([]*domain.Truck, error) {
	query := `
		SELECT id, plate, capacity_liters, pipeiro_id, created_at 
		FROM trucks 
		WHERE pipeiro_id = $1
	`
	rows, err := r.repo.QueryContext(ctx, query, uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no one truck by pipeiro found: %v", err)
		}
		return nil, fmt.Errorf("internal error to search truck by pipeiro: %v", err)
	}
	defer rows.Close()

	var trucks []*domain.Truck
	for rows.Next() {
		var t domain.Truck
		if err := rows.Scan(
			&t.ID,
			&t.Plate,
			&t.CapacityLiters,
			&t.PipeiroID,
			&t.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error to read truck: %v", err)
		}
		trucks = append(trucks, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while getting trucks by pipeiro: %v", err)
	}

	return trucks, nil
}

func (r *postgresRepo) GetCisterns(ctx context.Context) ([]*domain.Cistern, error) {
	query := `
		SELECT id, name, responsible_name, city, capacity_liters, latitude, longitude, created_at 
		FROM cisterns
	`
	rows, err := r.repo.QueryContext(ctx, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no one cistern found: %v", err)
		}
		return nil, fmt.Errorf("internal error to search cistern: %v", err)
	}
	defer rows.Close()

	var cisterns []*domain.Cistern
	for rows.Next() {
		var c domain.Cistern
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.ResponsabibleName,
			&c.City,
			&c.CapacityLiters,
			&c.Latitude,
			&c.Longitude,
			&c.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error to read cistern: %v", err)
		}
		cisterns = append(cisterns, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while getting cisterns: %v", err)
	}

	return cisterns, nil
}

func (r *postgresRepo) GetDeliveryByPipeiroUUID(ctx context.Context, uuid string) ([]*domain.Delivery, error) {
	query := `
		SELECT d.id, d.cistern_id, d.truck_id, d.scheduled_date, d.status, d.created_at, d.updated_at
		FROM deliveries d
		INNER JOIN trucks t ON d.truck_id = t.id
		WHERE t.pipeiro_id = $1
	`
	rows, err := r.repo.QueryContext(ctx, query, uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no one delivery by pipeiro found: %v", err)
		}
		return nil, fmt.Errorf("internal error to search delivery by pipeiro: %v", err)
	}
	defer rows.Close()

	var deliveries []*domain.Delivery
	for rows.Next() {
		var d domain.Delivery
		if rows.Scan(
			&d.ID,
			&d.CisternID,
			&d.TruckID,
			&d.ScheduledDate,
			&d.Status,
			&d.CreatedAt,
			&d.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error to read delivery: %v", err)
		}
		deliveries = append(deliveries, &d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while getting deliveries by pipeiro: %v", err)
	}

	return deliveries, nil
}

func (r *postgresRepo) GetDeliveryByTruckUUID(ctx context.Context, uuid string) ([]*domain.Delivery, error) {
	query := `
		SELECT id, cistern_id, truck_id, scheduled_date, status, created_at, updated_at 
		FROM deliveries 
		WHERE truck_id = $1
	`
	rows, err := r.repo.QueryContext(ctx, query, uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no one delivery by truck found: %v", err)
		}
		return nil, fmt.Errorf("internal error to search delivery by truck: %v", err)
	}
	defer rows.Close()

	var deliveries []*domain.Delivery
	for rows.Next() {
		var d domain.Delivery
		if rows.Scan(
			&d.ID,
			&d.CisternID,
			&d.TruckID,
			&d.ScheduledDate,
			&d.Status,
			&d.CreatedAt,
			&d.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error to read delivery: %v", err)
		}
		deliveries = append(deliveries, &d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while getting deliveries by pipeiro: %v", err)
	}

	return deliveries, nil
}
