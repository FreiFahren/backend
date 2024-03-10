package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
)

var conn *pgx.Conn

func CreateConnection() {
	var err error
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	conn, err = pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Connected to the database.")
}

func CloseConnection() {
	if conn != nil {
		err := conn.Close(context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close database connection: %v\n", err)
		} else {
			fmt.Println("Database connection closed.")
		}
	}
}

func CreateTicketInfoTable() {
	sql := `
    CREATE TABLE IF NOT EXISTS ticket_info (
        id SERIAL PRIMARY KEY,
        timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
        message TEXT,
        author BIGINT,
        line VARCHAR(3),
        station_name VARCHAR(255),
        station_id VARCHAR(10),
        direction_name VARCHAR(255),
        direction_id VARCHAR(10)
    );
    `
	_, err := conn.Exec(context.Background(), sql)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create table: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Table created or already exists.")
}

func InsertTicketInfo(timestamp *time.Time, message *string, author *int64, line, stationName, stationId, directionName, directionId *string) error {
	sql := `
    INSERT INTO ticket_info (timestamp, message, author, line, station_name, station_id, direction_name, direction_id)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
    `
	// Convert *string and *int64 directly to interface{} for pgx
	values := []interface{}{timestamp, message, author, line, stationName, stationId, directionName, directionId}

	_, err := conn.Exec(context.Background(), sql, values...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to insert ticket info: %v\n", err)
		return err
	}
	return nil
}

// a function that accepts a timestamp and will return the top 10 station ids for that time of day and day of the week
func GetHistoricStations(timestamp time.Time) ([]string, error) {
	// Extract hour and weekday
	hour := timestamp.Hour()
	weekday := timestamp.Weekday()

	// Construct SQL query
	sql := `
        SELECT station_id
        FROM ticket_info
        WHERE EXTRACT(HOUR FROM timestamp) = $1 AND EXTRACT(DOW FROM timestamp) = $2
        GROUP BY station_id
        ORDER BY COUNT(station_id) DESC
        LIMIT 5
    `
	// Execute query
	rows, err := conn.Query(context.Background(), sql, hour, weekday)
	if err != nil {
		return nil, fmt.Errorf("query execution error: %w", err)
	}
	defer rows.Close()

	// Collect station IDs
	var stationIDs []string
	for rows.Next() {
		var stationID string
		if err := rows.Scan(&stationID); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		stationIDs = append(stationIDs, stationID)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return stationIDs, nil
}
