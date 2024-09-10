package storage

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"gourd/internal/common"
)

func ConnectDB() *sql.DB {
	// Formulate the connection string
	cfg := common.GetActiveConfig()
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName)

	// Open a connection to the database
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Error().Err(err).Msg("Error connecting to the database")
	}

	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Error().Err(err).Msg("Error pinging to the database")
	}

	log.Info().Msg("Successfully connected to the database!")

	return db
}

func CreateTable(db *sql.DB) {
	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		id uuid PRIMARY KEY,
		firstname varchar(255) NOT NULL,
		lastname varchar(255) NOT NULL,
		is_admin boolean NOT NULL
	);`

	_, err := db.Exec(createUserTable)
	if err != nil {
		log.Error().Err(err).Msg("Error creating user table: ")
	}

	createSessionTable := `
	CREATE TABLE IF NOT EXISTS sessions (
		id uuid PRIMARY KEY,
		user_id uuid NOT NULL,
	    current_step integer DEFAULT 1,
		max_progress integer DEFAULT 1,
		repo varchar(255) NOT NULL,
		started timestamp,
		submitted timestamp,
		time_limit integer NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	_, err = db.Exec(createSessionTable)
	if err != nil {
		log.Error().Err(err).Msg("Error creating session table: ")
	}

	var count int
	query := "SELECT COUNT(*) FROM users"
	err = db.QueryRow(query).Scan(&count)
	if err != nil {
		log.Error().Err(err).Msg("Error checking user count: ")
	}

	if count == 0 {
		id := CreateUser(db, "Mr.", "Robot", true)
		log.Info().Msgf("Created admin token: %v\n", id)
	}

	log.Info().Msg("Table(s) created successfully or already existed.")
}

func CreateUser(db *sql.DB, firstname, lastname string, isAdmin bool) User {
	insertUser := `INSERT INTO users(id, firstname, lastname, is_admin) VALUES ($1, $2, $3, $4) RETURNING id;`
	id := uuid.New()
	_, err := db.Exec(insertUser, id, firstname, lastname, isAdmin)
	if err != nil {
		log.Error().Err(err).Msg("Error inserting user")
	}
	return User{
		ID:        id,
		Firstname: firstname,
		Lastname:  lastname,
		IsAdmin:   isAdmin,
	}
}

func CreateSession(db *sql.DB, userId uuid.UUID, repo string, timelimit int64) uuid.UUID {
	// SQL statement to insert a new person
	insertSQL := `
	INSERT INTO sessions (id, user_id, repo, time_limit)
	VALUES ($1, $2, $3, $4)
	RETURNING id;`

	id := uuid.New()
	_, err := db.Exec(insertSQL, id, userId, repo, timelimit)
	if err != nil {
		log.Error().Err(err).Msg("Error inserting session")
	}

	log.Info().Msgf("Session inserted successfully with ID %d\n", id)
	return id
}

func GetSessions(db *sql.DB) ([]HydratedSession, error) {
	selectSQL := `
		SELECT 
			s.id, s.current_step, s.max_progress, s.repo, s.started, s.submitted, s.time_limit,
			u.id, u.firstname, u.lastname, u.is_admin
		FROM sessions s
		JOIN users u ON s.user_id = u.id;
	`

	rows, err := db.Query(selectSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []HydratedSession

	for rows.Next() {
		var hs HydratedSession
		var user User

		err := rows.Scan(
			&hs.ID,
			&hs.CurrentStep,
			&hs.MaxProgress,
			&hs.Repo,
			&hs.Started,
			&hs.Submitted,
			&hs.Timelimit,
			&user.ID,
			&user.Firstname,
			&user.Lastname,
			&user.IsAdmin,
		)

		hs.User = user
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, hs)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}

func GetSession(db *sql.DB, token string) (HydratedSession, error) {
	selectSQL := `
		SELECT 
			s.id, s.current_step, s.max_progress, s.repo, s.started, s.submitted, s.time_limit,
			u.id, u.firstname, u.lastname, u.is_admin
		FROM sessions s
		JOIN users u ON s.user_id = u.id
		WHERE s.user_id = $1;
	`

	var session HydratedSession
	var user User
	row := db.QueryRow(selectSQL, token)

	err := row.Scan(
		&session.ID,
		&session.CurrentStep,
		&session.MaxProgress,
		&session.Repo,
		&session.Started,
		&session.Submitted,
		&session.Timelimit,
		&user.ID,
		&user.Firstname,
		&user.Lastname,
		&user.IsAdmin,
	)
	if err != nil {
		return session, err
	}
	session.User = user

	return session, nil
}

func UpdateSessionProgress(db *sql.DB, session HydratedSession) error {
	updateSQL := `UPDATE sessions SET current_step = $1, max_progress = $2 WHERE user_id = $3`

	_, err := db.Exec(updateSQL, session.CurrentStep, session.MaxProgress, session.User.ID)
	return err
}

func CheckUserExists(db *sql.DB, token string, shouldBeAdmin bool) bool {
	var query string
	var exists bool
	if !shouldBeAdmin {
		query = `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	} else {
		query = `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND is_admin = true)`
	}
	row := db.QueryRow(query, token)
	err := row.Scan(&exists)

	if err != nil {
		fmt.Printf("%v\n", err)
		return false
	}

	return exists
}

func GetCurrentStep(db *sql.DB, token string) int {
	query := `SELECT step FROM sessions WHERE user_id = $1`

	var step int
	row := db.QueryRow(query, token)
	err := row.Scan(&step)

	if err != nil {
		fmt.Printf("%v\n", err)
		return 0
	}

	return step
}
