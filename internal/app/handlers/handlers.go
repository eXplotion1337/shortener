package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"shortener/internal"
	"shortener/internal/app/handlers/service/repository"
	"shortener/internal/app/middleware"
	"shortener/internal/config"
	"strings"

	"github.com/go-chi/chi/v5"
)

func PostAddURL(w http.ResponseWriter, r *http.Request, config *config.Config, storage repository.Storage) {

	var userID string

	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		fmt.Println("userID not found in context")

	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	str, _ := middleware.UngzipData(body)

	sit, err := url.ParseRequestURI(str)
	if err != nil {
		fmt.Println("URL is not valid", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s := sit.String()
	sitr, err := url.PathUnescape(s)
	if err != nil {
		fmt.Println(err, sitr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := internal.GenerateRandomString(10)
	shortURL := config.BaseURL + "/" + id

	newItem := repository.InMemoryStorage{
		ID:       id,
		LongURL:  sitr,
		ShortURL: shortURL,
		UserID:   userID,
	}

	storage.SaveURL(&newItem)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))

}

func GetByID(w http.ResponseWriter, r *http.Request, config *config.Config, storage repository.Storage) {
	id := chi.URLParam(r, "id")

	long, _ := storage.GetLongURL(id)
	Location := strings.TrimSpace(long)
	w.Header().Set("Location", long)

	if Location != "" {
		http.Redirect(w, r, Location, http.StatusTemporaryRedirect)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func PostAPIShorten(w http.ResponseWriter, r *http.Request, config *config.Config, storage repository.Storage) {
	var requestData struct {
		URL string `json:"url"`
	}

	var userID string

	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		fmt.Println("userID not found in context")

	}

	// Декодируем тело запроса
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil || requestData.URL == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	sit, err := url.ParseRequestURI(requestData.URL)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s := sit.String()
	sitr, err := url.PathUnescape(s)
	if err != nil {
		fmt.Println(err, sitr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Назначаем случайный id
	id := internal.GenerateRandomString(10)
	respoID := config.BaseURL + "/" + id

	newItem := repository.InMemoryStorage{
		ID:       id,
		LongURL:  requestData.URL,
		ShortURL: respoID,
		UserID:   userID,
	}

	storage.SaveURL(&newItem)
	response := map[string]string{"result": respoID}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

}

func GetUrlsHandler(w http.ResponseWriter, r *http.Request) {
	var userID string

	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		fmt.Println("userID not found in context")

	}

	urls, err := repository.FindURL(userID)
	if err != nil {
		http.Error(w, "Не удалось получить список URL пользователя", http.StatusInternalServerError)
		return
	}

	if len(urls) == 0 {
		http.Error(w, "Нет данных", http.StatusNoContent)
		return
	}

	jsonResult, err := json.Marshal(urls)
	if err != nil {
		http.Error(w, "Ошибка сериализации JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResult)
}