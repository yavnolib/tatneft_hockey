package postgre

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

const (
	ROWS = "SELECT tablename FROM pg_tables WHERE schemaname = current_schema()"
	DROP = "DROP TABLE IF EXISTS %s CASCADE"
)

func RemoveTables(db *pgxpool.Pool) {
	defer db.Close()
	ctx := context.Background()
	rows, err := db.Query(ctx, ROWS)
	if err != nil {
		log.Fatalf("Failed to retrieve tables: %v\n", err)
	}
	defer rows.Close()
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Fatalf("Failed to scan table name: %v\n", err)
		}
		_, err = db.Exec(ctx, fmt.Sprintf(DROP, tableName))
		log.Printf("[INFO] Table %s dropped \n", tableName)

	}
	if err := rows.Err(); err != nil {
		log.Fatalf("Error iterating through tables: %v\n", err)
	}
}
