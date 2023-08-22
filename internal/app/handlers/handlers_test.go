package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"shortener/internal/app/handlers/service/repository"
	"shortener/internal/config"
	"testing"

	"github.com/go-chi/chi"
)

func TestPostAddURL(t *testing.T) {
	// Создаем JSON-тело запроса
	requestBody := []byte(`https://example.com`)

	// Создаем фейковый объект реквеста
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	// req.Header.Set("Content-Type", "application/json")

	// Создаем фейковый объект ResponseWriter
	rr := httptest.NewRecorder()

	// Создаем фейковый маршрутизатор
	r := chi.NewRouter()

	// Подготавливаем конфигурацию и хранилище
	// Замените это на создание фейковых объектов конфига и хранилища
	config := &config.Config{}
	storage := &repository.JSON{}

	// Вызываем тестируемую функцию
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		PostAddURL(w, r, config, storage)
	})
	r.Post("/", handler)

	// Выполняем запрос к маршруту
	r.ServeHTTP(rr, req)

	// Проверяем код ответа
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Ожидался статус %d, но получили %d", http.StatusCreated, status)
	}


}
