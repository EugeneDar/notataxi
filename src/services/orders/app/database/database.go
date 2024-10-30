package database

import (
	"errors"
	"sync"
)

// todo(eugenedar): implement a real database
type Database struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

func NewDatabase() *Database {
	return &Database{
		data: make(map[string]interface{}),
	}
}

func (db *Database) AddItem(key string, value interface{}) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[key] = value
}

func (db *Database) GetItem(key string) (interface{}, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	value, exists := db.data[key]
	if !exists {
		return nil, errors.New("item not found")
	}
	return value, nil
}

func (db *Database) DeleteItem(key string) (interface{}, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	value, exists := db.data[key]
	if !exists {
		return nil, errors.New("item not found")
	}
	delete(db.data, key)
	return value, nil
}
