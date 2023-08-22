package app

import (
	"database/sql"
	"log"
	"net/http"

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

	r.Post("/api/shorten/batch", func(w http.ResponseWriter, r *http.Request) {
		handlers.PostBatch(w, r, config, storage)
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

	} else if conf.TypeStorage == "FileStorage" {
		storage = repository.NewFileStorage(conf.StoragePath)

		err := repository.CreateFileIfNotExists(conf.StoragePath)
		if err != nil {

			log.Println("Ошибка создания файла", err)
		}

		err = repository.ReadJSONFile()
		if err != nil {
			log.Println("Ошибка чтения файла", err)
		}

	} else if conf.TypeStorage == "DataBaseStorage" {
		db, err := sql.Open("postgres", conf.DataBaseDSN)
		if err != nil {
			log.Println(err)
		}

		storage = repository.NewDatabaseStorage(db)
		repository.CheckBD(conf.DataBaseDSN)

	} else {
		log.Fatal("Не удалось инициализировать хранилище")
	}

	return storage
}
