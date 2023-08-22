package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

type FileStorage struct {
	filename string
	// mu       sync.Mutex
	// data     []InMemoryStorage
}

func NewFileStorage(filename string) *FileStorage {
	return &FileStorage{
		filename: os.Getenv("FILE_STORAGE_PATH"),
	}
}


func (fs *FileStorage) GetLongURL(id string) (string, error) {
	InMemoryCollection.Mutex.Lock()
	defer InMemoryCollection.Mutex.Unlock()

	for _, v := range InMemoryCollection.ObjectURL {
		if strings.EqualFold(v.ID, id) {
			return v.LongURL, nil
		}
	}

	return "", fmt.Errorf("URL not found")
}

var addData sync.Mutex

func (fs *FileStorage) SaveURL(longURL *InMemoryStorage) error {
	addData.Lock()
	defer addData.Unlock()

	file, err := os.OpenFile(fs.filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	if fileInfo.Size() == 0 {
		_, err = file.Write([]byte("{}"))
		if err != nil {
			return err
		}
	}

	jsonData, err := os.ReadFile(fs.filename)
	if err != nil {
		return err
	}

	var obj JSON
	if err = json.Unmarshal(jsonData, &obj); err != nil {
		return err
	}

	obj.ObjectURL = append(obj.ObjectURL, *longURL)

	if jsonData, err = json.Marshal(&obj); err != nil {
		return err
	}

	if _, err = file.WriteAt(jsonData, 0); err != nil {
		return err
	}

	return nil
}

func  ReadJSONFile() error {

	jsonFile, err := os.ReadFile(os.Getenv("FILE_STORAGE_PATH"))
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonFile, &InMemoryCollection)
	if err != nil {
		return err
	}

	return nil
}
