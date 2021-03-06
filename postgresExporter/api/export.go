package api

import (
	"context"
	"psqlexport/config"
	"psqlexport/database"
	"sync"

	"github.com/jackc/pgx"
)

// WorkerExport exports data from PostgreSQL database to .CSV files.
func WorkerExport(ctx context.Context, wg *sync.WaitGroup, conn *pgx.ConnPool, outputDir string, tasks <-chan config.Table) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-tasks:
			if ok {
				recordChan := make(chan []string)
				go CreateCSV(ctx, task.MaxLines, outputDir, task.Name, recordChan)
				database.Query(conn, task.Query, task.Name, recordChan)
			} else {
				return
			}
		}
	}
}
