// @title Database Monitoring API
// @version 0.0.1
// @description PostgreSQL database monitoring API.
// @host localhost:8080
// @BasePath /
// @contact.name Andranik Sarkisyan
// @contact.url https://www.linkedin.com/in/a-sarkisyan/
// @contact.email andrewsarkisyan@gmail.com
package main

import (
	"context"
	"encoding/json"
	"fmt"
	_ "go_pg_mon/docs"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"
)

type QueryRequest struct {
	DSN   string `json:"dsn"`
	Query string `json:"query"`
}
type LongRunningQueriesRequest struct {
	DSN           string `json:"dsn"`
	QueryDuration string `json:"query_duration"`
	PgDatabase    string `json:"pg_database"`
	QueryState    string `json:"query_state"`
}
type LargestTablesRequest struct {
	DSN        string `json:"dsn"`
	QueryLimit int    `json:"query_limit"`
}

// @Summary Execute Custom Query
// @Description Execute a custom SQL SELECT query on the specified database.
// @Tags Database
// @Accept json
// @Produce json
// @Param request body QueryRequest true "Query Request"
// @Success 200 {object} []map[string]interface{} "Query results"
// @Failure 400 {object} string "Invalid request body or parameters"
// @Failure 500 {object} string "Server error"
// @Router /db/custom-query [post]
func executeQuery(w http.ResponseWriter, r *http.Request) {
	// Set the response header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Parse the JSON body from the request
	var req QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Validate input parameters
	if req.DSN == "" || req.Query == "" {
		http.Error(w, "dsn, query are required", http.StatusBadRequest)
		return
	}

	// Validate that the query is a SELECT statement
	if !strings.HasPrefix(strings.ToUpper(req.Query), "SELECT") {
		http.Error(w, "Only SELECT queries are allowed.", http.StatusBadRequest)
		return
	}

	// Database connection
	pool, err := pgxpool.New(context.Background(), req.DSN)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to connect to database: %v", err), http.StatusInternalServerError)
		return
	}
	defer pool.Close()

	// Execute the query
	rows, err := pool.Query(context.Background(), req.Query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to execute query: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Fetch column names
	columnNames := rows.FieldDescriptions()
	columnNamesStr := make([]string, len(columnNames))
	for i, col := range columnNames {
		columnNamesStr[i] = string(col.Name)
	}

	// Prepare to store results dynamically
	var results []map[string]interface{}
	for rows.Next() {
		// Dynamically scan each row
		values := make([]interface{}, len(columnNames))
		pointers := make([]interface{}, len(columnNames))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			http.Error(w, fmt.Sprintf("Failed to scan row: %v", err), http.StatusInternalServerError)
			return
		}

		// Build a map for the row
		row := make(map[string]interface{})
		for i, colName := range columnNamesStr {
			row[colName] = values[i]
		}

		results = append(results, row)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf("Row iteration error: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert results to JSON
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal JSON: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the JSON data
	w.Write(jsonData)
}

// @Summary Retrieve Long-Running Queries
// @Description Get details of long-running queries in a PostgreSQL database.
// @Tags Database
// @Accept json
// @Produce json
// @Param request body LongRunningQueriesRequest true "Long-Running Queries Request"
// @Success 200 {object} []map[string]interface{} "Long-running queries"
// @Failure 400 {object} string "Invalid request body or parameters"
// @Failure 500 {object} string "Server error"
// @Router /db/long-running-queries [post]
func executeLongRunningQueries(w http.ResponseWriter, r *http.Request) {
	// Set the response header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Parse the JSON body from the request
	var req LongRunningQueriesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Validate input parameters
	if req.DSN == "" || req.QueryDuration == "" || req.PgDatabase == "" || req.QueryState == "" {
		http.Error(w, "dsn, query_duration, pg_database are required", http.StatusBadRequest)
		return
	}

	// Read the SQL query from the file
	queryBytes, err := os.ReadFile("./sql/long_running_queries.sql")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read SQL file: %v", err), http.StatusInternalServerError)
		return
	}
	query := string(queryBytes)

	// Replace placeholders with actual values
	query = strings.ReplaceAll(query, "$query_duration", req.QueryDuration)
	query = strings.ReplaceAll(query, "$pg_database", req.PgDatabase)
	query = strings.ReplaceAll(query, "$query_state", req.QueryState)

	// Database connection
	pool, err := pgxpool.New(context.Background(), req.DSN)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to connect to database: %v", err), http.StatusInternalServerError)
		return
	}
	defer pool.Close()

	// Execute the query
	rows, err := pool.Query(context.Background(), query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to execute query: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Fetch column names
	columnNames := rows.FieldDescriptions()
	columnNamesStr := make([]string, len(columnNames))
	for i, col := range columnNames {
		columnNamesStr[i] = string(col.Name)
	}

	// Prepare to store results dynamically
	var results []map[string]interface{}
	for rows.Next() {
		// Dynamically scan each row
		values := make([]interface{}, len(columnNames))
		pointers := make([]interface{}, len(columnNames))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			http.Error(w, fmt.Sprintf("Failed to scan row: %v", err), http.StatusInternalServerError)
			return
		}

		// Build a map for the row
		row := make(map[string]interface{})
		for i, colName := range columnNamesStr {
			row[colName] = values[i]
		}

		results = append(results, row)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf("Row iteration error: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert results to JSON
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal JSON: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the JSON data
	w.Write(jsonData)
}

// @Summary Get Largest Tables
// @Description Retrieve a list of the largest tables in the database by size.
// @Tags Database
// @Accept json
// @Produce json
// @Param request body LargestTablesRequest true "Largest Tables Request"
// @Success 200 {object} []map[string]interface{} "Largest tables information"
// @Failure 400 {object} string "Invalid request body or parameters"
// @Failure 500 {object} string "Server error"
// @Router /db/largest-tables [post]
func executeLargestTables(w http.ResponseWriter, r *http.Request) {
	// Set the response header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Parse the JSON body from the request
	var req LargestTablesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Validate input parameters
	if req.DSN == "" {
		http.Error(w, "dsn is required", http.StatusBadRequest)
		return
	}
	if req.QueryLimit < 1 || req.QueryLimit > 100 {
		http.Error(w, "query_limit must be an integer from 1 to 100", http.StatusBadRequest)
		return
	}

	// Read the SQL query from the file
	queryBytes, err := os.ReadFile("./sql/largest_tables.sql")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read SQL file: %v", err), http.StatusInternalServerError)
		return
	}
	query := string(queryBytes)

	// Replace placeholders with actual values
	query = strings.ReplaceAll(query, "$query_limit", strconv.Itoa(req.QueryLimit))

	// Database connection
	pool, err := pgxpool.New(context.Background(), req.DSN)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to connect to database: %v", err), http.StatusInternalServerError)
		return
	}
	defer pool.Close()

	// Execute the query
	rows, err := pool.Query(context.Background(), query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to execute query: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Fetch column names
	columnNames := rows.FieldDescriptions()
	columnNamesStr := make([]string, len(columnNames))
	for i, col := range columnNames {
		columnNamesStr[i] = string(col.Name)
	}

	// Prepare to store results dynamically
	var results []map[string]interface{}
	for rows.Next() {
		// Dynamically scan each row
		values := make([]interface{}, len(columnNames))
		pointers := make([]interface{}, len(columnNames))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			http.Error(w, fmt.Sprintf("Failed to scan row: %v", err), http.StatusInternalServerError)
			return
		}

		// Build a map for the row
		row := make(map[string]interface{})
		for i, colName := range columnNamesStr {
			row[colName] = values[i]
		}

		results = append(results, row)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf("Row iteration error: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert results to JSON
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal JSON: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the JSON data
	w.Write(jsonData)
}

func main() {
	// API endpoints
	http.HandleFunc("/db/custom-query", executeQuery)
	http.HandleFunc("/db/long-running-queries", executeLongRunningQueries)
	http.HandleFunc("/db/largest-tables", executeLargestTables)

	// Serve Swagger documentation
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	// Start the server
	port := ":8080"
	fmt.Printf("[INFO] - Server listening on http://localhost%s...\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("[ERROR] - Error starting server: %v", err)
	}
}
