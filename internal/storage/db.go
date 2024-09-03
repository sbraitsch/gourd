package storage

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"log"
)

// Database configuration constants
const (
	DB_USER     = "local"
	DB_PASSWORD = "pwd"
	DB_NAME     = "gourd_db"
	DB_HOST     = "localhost"
	DB_PORT     = 5432
)

func ConnectDB() *sql.DB {
	// Formulate the connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)

	// Open a connection to the database
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging the database: ", err)
	}

	fmt.Println("Successfully connected to the database!")

	return db
}

func CreateTable(db *sql.DB) {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS sessions (
		id uuid PRIMARY KEY,
		firstname varchar(255) NOT NULL,
		lastname varchar(255) NOT NULL,
	    step integer DEFAULT 0,
		started timestamp DEFAULT CURRENT_TIMESTAMP,
		submitted timestamp,
		time_limit integer NOT NULL
	);`

	// Execute the SQL statement
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Error creating session table: ", err)
	}

	createAdminTable := `
	CREATE TABLE IF NOT EXISTS admins (
		id uuid PRIMARY KEY
	);`

	// Execute the SQL statement
	_, err = db.Exec(createAdminTable)
	if err != nil {
		log.Fatal("Error creating admin table: ", err)
	}

	var count int
	query := "SELECT COUNT(*) FROM admins"
	err = db.QueryRow(query).Scan(&count)
	if err != nil {
		fmt.Printf("error checking user count: %v\n", err)
	}
	if count == 0 {
		insertAdmin := `INSERT INTO admins(id) VALUES ($1) RETURNING id;`
		id := uuid.New()
		err = db.QueryRow(insertAdmin, id).Scan(&id)
		if err != nil {
			log.Fatal("Error inserting admin: ", err)
		}
		fmt.Printf("Created admin token: %v\n", id)
	}

	fmt.Println("Table(s) created successfully or already existed.")
}

func CreateSession(db *sql.DB, firstname string, lastname string, timelimit int64) uuid.UUID {
	// SQL statement to insert a new person
	insertSQL := `
	INSERT INTO sessions (id, firstname, lastname, time_limit)
	VALUES ($1, $2, $3, $4)
	RETURNING id;`

	id := uuid.New()
	err := db.QueryRow(insertSQL, id, firstname, lastname, timelimit).Scan(&id)
	if err != nil {
		log.Fatal("Error inserting session: ", err)
	}

	fmt.Printf("Session inserted successfully with ID %d\n", id)
	return id
}

func GetSessions(db *sql.DB) ([]Session, error) {
	selectSQL := `SELECT * FROM sessions;`

	// Execute the SQL statement
	rows, err := db.Query(selectSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []Session

	for rows.Next() {
		var session Session
		err := rows.Scan(&session.ID, &session.Firstname, &session.Lastname, &session.Step, &session.Started, &session.Submitted, &session.Timelimit)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}

func GetSession(db *sql.DB, token string) (Session, error) {
	selectSQL := `SELECT * FROM sessions WHERE id = $1;`

	var session Session
	row := db.QueryRow(selectSQL, token)

	err := row.Scan(&session.ID, &session.Firstname, &session.Lastname, &session.Step, &session.Started, &session.Submitted, &session.Timelimit)
	if err != nil {
		return session, err
	}

	return session, nil
}

func CheckAdminStatus(db *sql.DB, token string) (bool, error) {
	selectSQL := `SELECT * FROM admins WHERE id = $1;`

	var admin string
	row := db.QueryRow(selectSQL, token)

	err := row.Scan(&admin)
	if err != nil {
		return false, err
	}

	return true, nil
}
