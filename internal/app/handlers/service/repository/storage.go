package repository

import (
	"fmt"
	"strings"
	"sync"
)

type InMemoryStorage struct {
	ID       string `json:"id"`
	LongURL  string `json:"longURL"`
	ShortURL string `json:"short_url"`
	UserID   string `json:"userID"`
}

type JSON struct {
	sync.Mutex
	ObjectURL []InMemoryStorage
}

var InMemoryCollection JSON

type Storage interface {
	SaveURL(longURL *InMemoryStorage) (sortURL string, err error)
	GetLongURL(id string) (longURL string, err error)
}

func (in *JSON) SaveURL(longURL *InMemoryStorage) (sortURL string, err error) {
	in.Lock()
	defer in.Unlock()
	InMemoryCollection.ObjectURL = append(InMemoryCollection.ObjectURL, *longURL)
	return "", nil
}

func (in *JSON) GetLongURL(id string) (longURL string, err error) {
	in.Lock()
	defer in.Unlock()
	fmt.Println(InMemoryCollection.ObjectURL, id)
	if len(InMemoryCollection.ObjectURL) > 0 {
		for _, v := range InMemoryCollection.ObjectURL {
			if strings.EqualFold(v.ID, id) {
				return v.LongURL, nil
			}
		}
	}

	return "", nil
}
