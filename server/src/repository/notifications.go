package repository

import (
	"github.com/jmoiron/sqlx"
	"log"
	"os"
)

type DevicesDatabaseInterface interface {
	AddDevice(id string, token string) error
	GetDevicesTokens(id string) ([]string, error)
}

type DevicesPersistentDatabase struct {
	db *sqlx.DB
}

func NewDevicesPersistentDatabase(db *sqlx.DB) (DevicesDatabaseInterface, error) {
	if err := prepareDatabase(db); err != nil {
		return nil, err
	}
	return &DevicesPersistentDatabase{db: db}, nil
}

func prepareDatabase(db *sqlx.DB) error {
	//create notifications_db if it doesn't exist
	query := `
		CREATE TABLE IF NOT EXISTS devices_db`
	if os.Getenv("ENVIRONMENT") != "HEROKU" {
		query += `_test`
	}

	query += `(
		    			user_id TEXT,
		    			token TEXT,
		    			PRIMARY KEY (user_id, token)
		);
	`
	log.Println("Executing query: ", query)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	log.Println("Database created successfully")
	return nil
}

func (db *DevicesPersistentDatabase) AddDevice(id string, token string) error {
	_, err := db.db.Exec(`INSERT INTO devices_db (user_id, token) VALUES ($1, $2)`, id, token)
	return err
}

func (db *DevicesPersistentDatabase) GetDevicesTokens(id string) ([]string, error) {
	rows, err := db.db.Query(`SELECT token FROM devices_db WHERE user_id=$1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tokens := []string{}
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}
