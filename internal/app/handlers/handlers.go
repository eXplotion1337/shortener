package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"shortener/internal"
	"shortener/internal/app/handlers/service/repository"
	"shortener/internal/config"
	"strings"

	"github.com/go-chi/chi/v5"
)

func PostAddURL(w http.ResponseWriter, r *http.Request, config *config.Config, storage repository.Storage) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// str, _ := middleware.UngzipData(b)

	sit, err := url.ParseRequestURI(string(body))
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
		UserID:   "1",
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
