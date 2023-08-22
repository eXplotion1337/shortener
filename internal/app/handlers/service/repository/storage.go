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
	Delete   bool   `json:"DeleteFlag"`
}

type JSON struct {
	sync.Mutex
	ObjectURL []InMemoryStorage
}

var InMemoryCollection JSON

type Longbatch struct {
	LongURL    string
	DeleteFlag bool
}

type Storage interface {
	SaveURL(longURL *InMemoryStorage) (sortURL string, err error)
	GetLongURL(id string) (longURL Longbatch, err error)
	Delete(ids []string, userID string)
}

func (in *JSON) SaveURL(longURL *InMemoryStorage) (sortURL string, err error) {
	in.Lock()
	defer in.Unlock()
	InMemoryCollection.ObjectURL = append(InMemoryCollection.ObjectURL, *longURL)
	return "", nil
}

func (in *JSON) GetLongURL(id string) (Longbatch, error) {
	in.Lock()
	defer in.Unlock()

	var batch Longbatch

	if len(InMemoryCollection.ObjectURL) > 0 {
		for _, v := range InMemoryCollection.ObjectURL {
			if strings.EqualFold(v.ID, id) {
				batch.LongURL = v.LongURL
				batch.DeleteFlag = v.Delete
				return batch, nil
			}
		}
	}

	return batch, nil
}

func (in *JSON) Delete(ids []string, userID string) {
	in.Lock()
	defer in.Unlock()

	if len(InMemoryCollection.ObjectURL) > 0 {
		for _, k := range ids {
			for _, v := range InMemoryCollection.ObjectURL {
				if k == v.ID {
					if v.UserID == userID {
						v.Delete = true
						fmt.Println(v.Delete)
					}
				}
			}
		}
	}
}
