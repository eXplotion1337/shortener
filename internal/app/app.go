package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"shortener/internal/app/handlers"
	"shortener/internal/app/handlers/service/repository"
	"shortener/internal/config"

	"github.com/go-chi/chi/v5"
)

func Run(config *config.Config, storage repository.Storage) error {
	r := chi.NewRouter()

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.PostAddURL(w, r, config, storage)
	})

	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetByID(w, r, config, storage)
	})

	r.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		handlers.PostAPIShorten(w, r, config, storage)
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
	if conf.TypeStorage == "in-memory" {
		storage = &repository.JSON{}
		fmt.Println("in-memory")
	} else if conf.TypeStorage == "file" {
		storage = repository.NewFileStorage(os.Getenv("FILE_STORAGE_PATH"))
		err := repository.ReadJSONFile()
		if err != nil{
			log.Println("Ошибка чтения файла", err)
		}
		fmt.Println("file")
	} else {
		fmt.Println("db")
	}
	return storage
}
