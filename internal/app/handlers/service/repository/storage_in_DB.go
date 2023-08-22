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
		INSERT INTO urls (id, long_url, short_url, user_id, delete_flag)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (long_url) DO UPDATE SET long_url = EXCLUDED.long_url
		RETURNING short_url
	`
	var shortURL string

	err := ds.db.QueryRow(insertQuery, item.ID, item.LongURL, item.ShortURL, item.UserID, item.Delete).Scan(&shortURL)
	if err != nil {

		return "", err
	}

	if shortURL != item.ShortURL {
		return shortURL, err
	}

	return "", nil
}


func (ds *DatabaseStorage) GetLongURL(id string) (Longbatch, error) {
	var U_data Longbatch
	selectQuery := `
		SELECT long_url, delete_flag FROM urls WHERE id = $1
	`

	err := ds.db.QueryRow(selectQuery, id).Scan(&U_data.LongURL, &U_data.Delete_flag)
	if err != nil {
		if err == sql.ErrNoRows {
			return U_data, fmt.Errorf("URL not found")
		}
		return U_data, err
	}

	return U_data, nil
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
			user_id VARCHAR(36) NOT NULL,
			delete_flag BOOL NOT NULL
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

func (ds *DatabaseStorage) Delete(ids []string, userID string){
	for _, k := range ids{
		_, err := ds.db.Exec("UPDATE urls SET delete_flag = true WHERE id = $1 AND user_id = $2", k, userID)
		if err != nil {
			log.Println("не удалось заменить флаг")
		}

	}


}
