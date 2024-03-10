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
