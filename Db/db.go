package Db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
)

var db *sql.DB
var err error

func Connect(dbname string) error {

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err = sql.Open("postgres", psqlconn)
	if err != nil {
		fmt.Println("[DB] Error connecting to postgres server.")
		return err
	}

	err = db.Ping()

	return err

}

func CreateTableDG() {
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS dailygrowth (
		id SERIAL PRIMARY KEY,
		date_time VARCHAR,
		node_count INT,
		total_users INT,
		total_comments INT,
		total_posts INT,
		exec_time VARCHAR
	);
	
	`)
	if err != nil {
		fmt.Println("[CREATE DG] Cant create.", err.Error())
		return
	}
	fmt.Println("[DB] Created Table dailygrowth")
}

func CreateTableNodes() string {
	tableName := "Nodes"

	_, err := db.Exec(fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		id SERIAL PRIMARY KEY,
		date_time VARCHAR,
		name VARCHAR,
		version VARCHAR,
		tld VARCHAR,
		total_users INT,
		total_comments INT,
		total_posts INT,
		metadata JSONB
	);
	`, tableName))
	if err != nil {
		fmt.Println("[CREATE NODES] Cant create.", err.Error())
		return ""
	}
	fmt.Println("[DB] Created Table Nodes")
	return tableName
}

func CreateTableError() string {
	errorTableName := "Error"
	_, err = db.Exec(fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		id SERIAL PRIMARY KEY,
		date_time TIMESTAMP,
		url VARCHAR,
		reason VARCHAR
	);
	
	`, errorTableName))
	if err != nil {
		fmt.Println("[CREATE Error] Cant create.", err.Error())
		return ""
	}
	fmt.Println("[DB] Created Table error")
	return errorTableName
}

func InsertDG(node_count int, total_users int, total_comments int, total_posts int, tableName string, errorTableName string, execTime string) {
	dateTime := time.Now()
	_, err = db.Exec(`
		INSERT INTO dailygrowth (date_time, node_count, total_users, total_comments, total_posts, table_name, error_table_name, exec_time)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, dateTime.Format("2006_01_02_15_04_05"), node_count, total_users, total_comments, total_posts, tableName, errorTableName, execTime)
	if err != nil {
		fmt.Println("[INSERT DG] Cant execute query.", err.Error())
		return
	}
	fmt.Println("[DB] Inserted daily_growth_stats")
}

func InsertNodes(tableName string, name string, version string, tld string, users int, comments int, posts int, metadata map[string]interface{}) {
	dateTime := time.Now()

	jsonBytes, err := json.Marshal(metadata)
	if err != nil {
		fmt.Println("[INSERT NODES] Cant marshal metadata.")
		return
	}
	_, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %s (date_time, name, version, tld, total_users, total_comments, total_posts, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, tableName), dateTime.Format("2006-01-02 15:04:05"), name, version, tld, users, comments, posts, jsonBytes)
	if err != nil {
		fmt.Println("[INSERT NODES] Cant execute query.", err.Error())
		return
	}
	fmt.Println("[DB] Inserted Nodes")
}

func Readpsql(query string) ([]map[string]interface{}, error) {

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return nil, err
		}

		rowMap := make(map[string]interface{})
		for i, colName := range columns {
			rowMap[colName] = values[i]
		}

		results = append(results, rowMap)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func InsertError(errorTableName string, url string, reason string) {
	dateTime := time.Now()
	_, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %s (date_time, url, reason)
		VALUES ($1, $2, $3)
	`, errorTableName), dateTime.Format("2006-01-02 15:04:05"), url, reason)
	if err != nil {
		fmt.Println("[INSERT Error] Cant execute query.", err.Error())
		return
	}
	fmt.Println("[DB] Inserted error")
}
