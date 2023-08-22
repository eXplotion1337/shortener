package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type DatabaseStorage struct {
	db *sql.DB
}

func NewDatabaseStorage(db *sql.DB) *DatabaseStorage {
	return &DatabaseStorage{
		db: db,
	}
}

func (ds *DatabaseStorage) SaveURL(item *InMemoryStorage) (string, error) {
	insertQuery := `
		INSERT INTO urls (id, long_url, short_url, user_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (long_url) DO UPDATE SET long_url = EXCLUDED.long_url
		RETURNING short_url
	`
	var shortURL string

	err := ds.db.QueryRow(insertQuery, item.ID, item.LongURL, item.ShortURL, item.UserID).Scan(&shortURL)
	if err != nil {

		return "", err
	}

	if shortURL != item.ShortURL {
		return shortURL, err
	}

	return "", nil
}

func (ds *DatabaseStorage) GetLongURL(id string) (string, error) {

	var longURL string

	selectQuery := `
		SELECT long_url FROM urls WHERE id = $1
	`

	err := ds.db.QueryRow(selectQuery, id).Scan(&longURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("URL not found")
		}
		return "", err
	}

	return longURL, nil
}

func tableExists(db *sql.DB, tableName string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)", tableName).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func createURLsTable(db *sql.DB) error {
	createTableQuery := `
		CREATE TABLE urls (
			id VARCHAR(36) PRIMARY KEY,
			long_url TEXT UNIQUE NOT NULL,
			short_url VARCHAR(100) NOT NULL,
			user_id VARCHAR(36) NOT NULL
		)
	`

	_, err := db.Exec(createTableQuery)
	return err
}

func CheckBD(databaseDSN string) {
	if databaseDSN == "" {
		log.Println("DATABASE_DSN environment variable is not set")
	}

	// Открытие соединения с базой данных
	db, err := sql.Open("postgres", databaseDSN)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	tableName := "urls"

	// Проверка наличия таблицы
	exists, err := tableExists(db, tableName)
	if err != nil {
		log.Println("Нет доступа к БД")
	}

	if !exists {
		err := createURLsTable(db)
		if err != nil {
			log.Println("Нет доступа к БД")
		}
		log.Println("Table 'urls' created successfully.")
	} else {
		log.Println("Table 'urls' already exists.")
	}
}
