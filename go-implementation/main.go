package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// connectToPrimary connects to the primary PostgreSQL database, performs a select and insert, and then selects again.
func connectToPrimary(connectionConfig string) {
	db, err := sql.Open("postgres", connectionConfig)
	if err != nil {
		log.Fatalf("Error opening connection to Primary PostgreSQL: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to Primary PostgreSQL: %v", err)
	}
	fmt.Println("Connected to Primary PostgreSQL!")

	// Perform a SELECT query
	rows, err := db.Query(`SELECT * FROM my_table`)
	if err != nil {
		log.Fatalf("Error executing SELECT query on Primary: %v", err)
	}
	defer rows.Close()

	// Get column names
	columnNames, err := rows.Columns()
	if err != nil {
		log.Fatalf("Error getting column names: %v", err)
	}
	fmt.Println("Column Names:", columnNames)

	var columns = make([]interface{}, len(columnNames))
	columnPointers := make([]interface{}, len(columnNames))
	for i := range columns {
		columnPointers[i] = &columns[i]
	}

	fmt.Println("Data:")
	for rows.Next() {
		err := rows.Scan(columnPointers...)
		if err != nil {
			log.Fatalf("Error scanning row from Primary: %v", err)
		}
		fmt.Println(columns)
	}
	if err = rows.Err(); err != nil {
		log.Fatalf("Error iterating through rows from Primary: %v", err)
	}

	// Perform an INSERT query
	result, err := db.Exec(`INSERT INTO my_table (data) VALUES ($1)`, "Testing Insert - 6")
	if err != nil {
		log.Fatalf("Error executing INSERT query on Primary: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting affected rows after insert on Primary: %v", err)
	}
	fmt.Println("Insert Status:", rowsAffected)

	// Perform another SELECT query to see the inserted data
	tableRows, err := db.Query(`SELECT * FROM my_table`)
	if err != nil {
		log.Fatalf("Error executing second SELECT query on Primary: %v", err)
	}
	defer tableRows.Close()

	// Get column names again (they should be the same)
	tableColumnNames, err := tableRows.Columns()
	if err != nil {
		log.Fatalf("Error getting column names for second select: %v", err)
	}
	fmt.Println("Column Names (Second Select):", tableColumnNames)

	fmt.Println("Retrieving db:")
	var tableColumns = make([]interface{}, len(tableColumnNames))
	tableColumnPointers := make([]interface{}, len(tableColumnNames))
	for i := range tableColumns {
		tableColumnPointers[i] = &tableColumns[i]
	}
	for tableRows.Next() {
		err := tableRows.Scan(tableColumnPointers...)
		if err != nil {
			log.Fatalf("Error scanning row during second select from Primary: %v", err)
		}
		fmt.Println(tableColumns)
	}
	if err = tableRows.Err(); err != nil {
		log.Fatalf("Error iterating through rows during second select from Primary: %v", err)
	}
}

// connectToReplica connects to the replica PostgreSQL database and performs a select query.
func connectToReplica(connectionConfig string) {
	db, err := sql.Open("postgres", connectionConfig)
	if err != nil {
		log.Fatalf("Error opening connection to Replica PostgreSQL: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to Replica PostgreSQL: %v", err)
	}
	fmt.Println("Connected to Replica PostgreSQL!")

	// Perform a SELECT query
	rows, err := db.Query(`SELECT * FROM my_table`)
	if err != nil {
		log.Fatalf("Error executing SELECT query on Replica: %v", err)
	}
	defer rows.Close()

	// Get column names
	columnNames, err := rows.Columns()
	if err != nil {
		log.Fatalf("Error getting column names from Replica: %v", err)
	}
	fmt.Println("Column Names (Replica):", columnNames)

	var columns = make([]interface{}, len(columnNames))
	columnPointers := make([]interface{}, len(columnNames))
	for i := range columns {
		columnPointers[i] = &columns[i]
	}

	fmt.Println("Data from Replica:")
	for rows.Next() {
		err := rows.Scan(columnPointers...)
		if err != nil {
			log.Fatalf("Error scanning row from Replica: %v", err)
		}
		fmt.Println(columns)
	}
	if err = rows.Err(); err != nil {
		log.Fatalf("Error iterating through rows from Replica: %v", err)
	}
}

func main() {
	// Connection string for the primary server
	primaryConfig := "host=localhost port=5432 user=postgres password='' dbname=test_replication sslmode=disable" // Replace with your password if set

	// Connection string for the standby server
	standbyConfig := "host=localhost port=5433 user=postgres password='' dbname=test_replication sslmode=disable" // If you used a different port for the standby

	connectToPrimary(primaryConfig)
	connectToReplica(standbyConfig) // You can connect to the standby for read operations
}