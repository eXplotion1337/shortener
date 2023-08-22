package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"shortener/internal/app/handlers"
	"shortener/internal/app/handlers/service/repository"
	"shortener/internal/app/middleware"
	"shortener/internal/config"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
)

func Run(config *config.Config, storage repository.Storage) error {
	r := chi.NewRouter()
	r.Use(middleware.GZipMiddleware)
	r.Use(middleware.SetUserIDCookie)

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.PostAddURL(w, r, config, storage)
	})

	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetByID(w, r, config, storage)
	})

	r.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		handlers.PostAPIShorten(w, r, config, storage)
	})

	r.Get("/api/user/urls", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUrlsHandler(w, r)
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		handlers.PingDB(w, r, config, storage)
	})

	log.Printf("Сервер запущен на %s", config.ServerAddr)
	log.Printf("Base URL  %s", config.BaseURL)
	log.Printf("Файл для сохранения данных расположен %s", config.StoragePath)
	log.Printf("База данных  %s", config.DataBaseDSN)
	log.Printf("Хранение данных реализовано через  %s", config.TypeStorage)

	err := http.ListenAndServe(config.ServerAddr, r)
	if err != nil {
		log.Fatal("Ошибка старта сервера", err)
		return err
	}

	return nil
}

func InitStorage(conf *config.Config) repository.Storage {
	var storage repository.Storage
	if conf.TypeStorage == "In-memoryStorage" {
		storage = &repository.JSON{}
		fmt.Println("in-memory")
	} else if conf.TypeStorage == "FileStorage" {
		storage = repository.NewFileStorage(os.Getenv("FILE_STORAGE_PATH"))
		err := repository.CreateFileIfNotExists(conf.StoragePath)
		if err != nil {
			// Обрабатываем ошибку
			fmt.Println("Ошибка создания файла", err)
		}
		err = repository.ReadJSONFile()
		if err != nil {
			log.Println("Ошибка чтения файла", err)
		}
		fmt.Println("file")
	} else if conf.TypeStorage == "DataBaseStorage" {
		db, err := sql.Open("postgres", conf.DataBaseDSN)
		if err != nil {
			fmt.Println(err)
		}
		// defer db.Close()
		storage = repository.NewDatabaseStorage(db)
		repository.CheckBD(conf.DataBaseDSN)

		fmt.Println("db")
	} else {
		log.Fatal("Не удалось инициализировать хранилище")
	}

	return storage
}
