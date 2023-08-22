package main

import (
	"log"
	"shortener/internal/config"
	"shortener/internal/app"
)

func main() {
	config, err := config.InitConfig()
	if err != nil {
		log.Fatal("Ошибка загрузки конфига", err)
	}

	storage := app.InitStorage(config)

	if err := app.Run(config, storage); err != nil {
		log.Fatal("Ошибка старта сервера", err)
	}

}
