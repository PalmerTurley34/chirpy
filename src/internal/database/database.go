package database

import (
	"encoding/json"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps        map[int]Chirp           `json:"chirps"`
	Users         map[int]User            `json:"users"`
	RefreshTokens map[string]RefreshToken `json:"refresh_tokens"`
}

func NewDB(path string) (*DB, error) {
	newDB := DB{path, &sync.RWMutex{}}
	err := newDB.ensureDB()
	return &newDB, err
}

func (db *DB) ensureDB() error {
	data, err := os.ReadFile(db.path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if len(data) != 0 {
		return nil
	}
	err = db.writeDB(db.createDBStruct())
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) createDBStruct() DBStructure {
	emptyDBStruct := DBStructure{
		make(map[int]Chirp),
		make(map[int]User),
		make(map[string]RefreshToken),
	}
	return emptyDBStruct
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	data, err := os.ReadFile(db.path)
	if err != nil {
		if os.IsNotExist(err) {
			return db.createDBStruct(), nil
		}
		return DBStructure{}, err
	}
	dbStruct := DBStructure{}
	err = json.Unmarshal(data, &dbStruct)
	if err != nil {
		return DBStructure{}, err
	}
	return dbStruct, nil
}

func (db *DB) writeDB(dbStruct DBStructure) error {
	toWrite, err := json.MarshalIndent(dbStruct, "", "	")
	if err != nil {
		return err
	}
	db.mux.Lock()
	defer db.mux.Unlock()
	err = os.WriteFile(db.path, toWrite, 0666)
	if err != nil {
		return err
	}
	return nil
}
