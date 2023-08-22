package repository

import "sync"

type Rez struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"original_url"`
}

var mu sync.Mutex

func FindURL(input string) ([]Rez, error) {
	var result []Rez

	mu.Lock()
	defer mu.Unlock()

	if len(InMemoryCollection.ObjectURL) > 0 {
		for _, record := range InMemoryCollection.ObjectURL {
			if record.UserID == input {
				prom := Rez{LongURL: record.LongURL, ShortURL: record.ShortURL}
				result = append(result, prom)
			}
		}
	}

	return result, nil
}
