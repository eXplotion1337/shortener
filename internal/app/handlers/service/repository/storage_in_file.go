package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

func (fs *FileStorage) GetLongURL(id string) (Longbatch, error) {
	InMemoryCollection.Mutex.Lock()
	defer InMemoryCollection.Mutex.Unlock()
	var batch Longbatch

	for _, v := range InMemoryCollection.ObjectURL {
		if strings.EqualFold(v.ID, id) {
			batch.LongURL = v.LongURL
			batch.DeleteFlag = v.Delete
			return batch, nil
		}
	}

	return batch, fmt.Errorf("URL not found")
}

var addData sync.Mutex

func (fs *FileStorage) SaveURL(longURL *InMemoryStorage) (shortURL string, err error) {
	addData.Lock()
	defer addData.Unlock()

	file, err := os.OpenFile(fs.filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}

	if fileInfo.Size() == 0 {
		_, err = file.Write([]byte("{}"))
		if err != nil {
			return "", err
		}
	}

	jsonData, err := os.ReadFile(fs.filename)
	if err != nil {
		return "", err
	}

	var obj JSON
	if err = json.Unmarshal(jsonData, &obj); err != nil {
		return "", err
	}

	obj.ObjectURL = append(obj.ObjectURL, *longURL)
	InMemoryCollection.ObjectURL = append(InMemoryCollection.ObjectURL, *longURL)

	if jsonData, err = json.Marshal(&obj); err != nil {
		return "", err
	}
	if _, err = file.WriteAt(jsonData, 0); err != nil {
		return "", err
	}

	return "", nil
}

func ReadJSONFile() error {

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

func CreateFileIfNotExists(path string) error {
	// path := os.Getenv("FILE_STORAGE_PATH")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// получаем директорию, где должен быть файл
		dir := filepath.Dir(path)

		// создаем все директории в пути к файлу
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// создаем сам файл
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

func (fs *FileStorage) Delete(ids []string, userID string) {
	InMemoryCollection.Mutex.Lock()
	defer InMemoryCollection.Mutex.Unlock()

	for _, k := range ids {
		for _, v := range InMemoryCollection.ObjectURL {
			if k == v.ID {
				if userID == v.UserID {
					v.Delete = true
					fmt.Println(v.Delete)
				}
			}
		}
	}
}
