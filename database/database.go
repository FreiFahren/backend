package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TicketInfo struct {
	Timestamp    time.Time
	Station_ID   string
	Line         sql.NullString
	Direction_ID sql.NullString
}

var pool *pgxpool.Pool

func Config() *pgxpool.Config {
	const defaultMaxConns = int32(4)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	dbConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		log.Println("Acquiring a connection...")
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		log.Println("Connection has been released")
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		log.Println("Closed the connection pool")
	}

	return dbConfig
}

func CreatePool() {
	var err error

	p, err := pgxpool.NewWithConfig(context.Background(), Config())
	if err != nil {
		log.Fatal("Error while creating connection to the database!!")
	}

	pool = p

}
func ClosePool() {
	if pool != nil {
		pool.Close()
	}
}

func CreateTicketInfoTable() {
	sql := `
	CREATE TABLE IF NOT EXISTS ticket_info (
		id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
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

	_, err := pool.Exec(context.Background(), sql)
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

	_, err := pool.Exec(context.Background(), sql, values...)
	log.Println("Inserting ticket info...")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to insert ticket info: %v\n", err)
		return err
	}
	return nil
}

func GetHistoricStations(timestamp time.Time) ([]TicketInfo, error) {
	// Extract hour and weekday
	hour := timestamp.Hour()
	weekday := timestamp.Weekday()

	// get top 5 stations for the given hour and weekday
	sql := `
        SELECT station_id
        FROM ticket_info
        WHERE EXTRACT(HOUR FROM timestamp) = $1 AND EXTRACT(DOW FROM timestamp) = $2
		AND station_name IS NOT NULL
		AND station_id IS NOT NULL
        GROUP BY station_id
        ORDER BY COUNT(station_id) DESC
        LIMIT 20;
    `
	// Execute query
	rows, err := pool.Query(context.Background(), sql, hour, weekday)
	if err != nil {
		return nil, fmt.Errorf("query execution error: %w", err)
	}
	defer rows.Close()

	var ticketInfoList []TicketInfo

	for rows.Next() {
		var ticketInfo TicketInfo

		if err := rows.Scan(&ticketInfo.Station_ID); err != nil {
			return nil, fmt.Errorf("error scanning row (historic data): %w", err)
		}
		ticketInfoList = append(ticketInfoList, ticketInfo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows (historic data): %w", err)
	}

	if len(ticketInfoList) == 0 {
		fmt.Println("No historic data found")
	}

	return ticketInfoList, nil
}

func GetLatestStationCoordinates() ([]TicketInfo, error) {
	sql := `SELECT timestamp, station_id, direction_id, line
            FROM ticket_info
            WHERE timestamp >= NOW() - INTERVAL '180 minutes'
            AND station_name IS NOT NULL
			AND station_id IS NOT NULL;`

	rows, err := pool.Query(context.Background(), sql)
	log.Println("Getting Data...")

	if err != nil {
		return nil, fmt.Errorf("query execution error: %w", err)
	}

	defer rows.Close()

	var ticketInfoList []TicketInfo

	for rows.Next() {
		var ticketInfo TicketInfo
		if err := rows.Scan(&ticketInfo.Timestamp, &ticketInfo.Station_ID, &ticketInfo.Direction_ID, &ticketInfo.Line); err != nil {
			return nil, fmt.Errorf("error scanning row (latest station coordinate data): %w", err)
		}

		ticketInfoList = append(ticketInfoList, ticketInfo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows (latest station coordinate data): %w", err)
	}

	return ticketInfoList, nil
}
