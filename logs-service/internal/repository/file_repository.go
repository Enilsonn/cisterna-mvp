package repository

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"cisterna-mvp/logs-service/internal/consumer"
)

type FileRepository struct {
	path string
}

func NewFileRepository(basePath string) *FileRepository {
	return &FileRepository{path: filepath.Join(basePath, "logs.txt")}
}

func (r *FileRepository) Save(entry consumer.LogEntry) error {
	file, err := os.OpenFile(r.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = fmt.Fprintf(writer, "%s | truck=%s | lat=%.6f | lon=%.6f | received_at=%s\n",
		entry.Timestamp.Format(time.RFC3339),
		entry.TruckID,
		entry.Latitude,
		entry.Longitude,
		entry.ReceivedAt.Format(time.RFC3339),
	)
	if err != nil {
		return err
	}

	return writer.Flush()
}
