package main

import (
  _ "github.com/lib/pq"
)

// Database configuration constants
const (
	DB_USER     = "yourusername"
	DB_PASSWORD = "yourpassword"
	DB_NAME     = "yourdbname"
	DB_HOST     = "localhost"
	DB_PORT     = 5432
)

// Person struct represents a person in the database
type Person struct {
	ID    int
	Name  string
	Email string
}

// ConnectDB initializes a connection to the PostgreSQL database
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

// CreateTable creates the "people" table if it does not exist
func CreateTable(db *sql.DB) {
	// SQL statement to create the "people" table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS people (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT UNIQUE NOT NULL
	);`

	// Execute the SQL statement
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Error creating table: ", err)
	}

	fmt.Println("Table created successfully or already exists.")
}

// InsertPerson inserts a new person into the "people" table
func InsertPerson(db *sql.DB, name, email string) {
	// SQL statement to insert a new person
	insertSQL := `
	INSERT INTO people (name, email)
	VALUES ($1, $2)
	RETURNING id;`

	// Execute the SQL statement and get the new ID
	var id int
	err := db.QueryRow(insertSQL, name, email).Scan(&id)
	if err != nil {
		log.Fatal("Error inserting person: ", err)
	}

	fmt.Printf("Person inserted successfully with ID %d\n", id)
}

// GetPeople retrieves all people from the "people" table
func GetPeople(db *sql.DB) ([]Person, error) {
	// SQL statement to select all people
	selectSQL := `
	SELECT id, name, email FROM people;`

	// Execute the SQL statement
	rows, err := db.Query(selectSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Prepare a slice to hold the results
	var people []Person

	// Iterate through the result set
	for rows.Next() {
		var person Person
		err := rows.Scan(&person.ID, &person.Name, &person.Email)
		if err != nil {
			return nil, err
		}
		people = append(people, person)
	}

	// Check for errors from the row iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return people, nil
}

