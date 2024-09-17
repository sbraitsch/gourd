package storage

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"gourd/internal/common"
)

// ConnectDB connects to the local postgres database using the active config. /*
func ConnectDB() *sql.DB {
	cfg := common.GetActiveConfig()
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to the database")
	}
	err = db.Ping()
	if err != nil {
		log.Fatal().Err(err).Msg("Error pinging to the database")
	}
	log.Info().Msg("Successfully connected to the database!")
	return db
}

// InitAdminUser creates an initial admin user.
func InitAdminUser(db *sql.DB) {
	var count int
	query := "SELECT COUNT(*) FROM users"
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		log.Fatal().Err(err).Msg("Error checking user count: ")
	}

	if count == 0 {
		id := CreateUser(db, "Mr.", "Robot", true)
		log.Info().Msgf("Created admin token: %v\n", id)
	}
}

// CreateUser inserts a new row into the users table.
func CreateUser(db *sql.DB, firstname, lastname string, isAdmin bool) User {
	insertUser := `INSERT INTO users(id, firstname, lastname, is_admin) VALUES ($1, $2, $3, $4) RETURNING id;`
	id := uuid.New()
	_, err := db.Exec(insertUser, id, firstname, lastname, isAdmin)
	if err != nil {
		log.Fatal().Err(err).Msg("Error inserting user")
	}
	return User{
		ID:        id,
		Firstname: firstname,
		Lastname:  lastname,
		IsAdmin:   isAdmin,
	}
}

// CreateSession inserts a new row into the sessions table.
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

// GetSessions returns all sessions from the database and hydrates the user data.
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

// GetSession returns the hydrated session of the given user/token.
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

// UpdateSessionProgress updates the current_step and max_progress fields of a session.
func UpdateSessionProgress(db *sql.DB, session HydratedSession) error {
	updateSQL := `UPDATE sessions SET current_step = $1, max_progress = $2 WHERE user_id = $3`

	_, err := db.Exec(updateSQL, session.CurrentStep, session.MaxProgress, session.User.ID)
	return err
}

// CheckUserExists returns whether a given token maps to an existing user or not.
func CheckUserExists(db *sql.DB, token string) (exists, isAdmin bool) {
	query := `SELECT is_admin FROM users WHERE id = $1`
	row := db.QueryRow(query, token)
	err := row.Scan(&isAdmin)

	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	exists = true
	return
}
