package database

type Chirp struct {
	ID     int    `json:"id"`
	UserID int    `json:"author_id"`
	Body   string `json:"body"`
}

func (db *DB) CreateChirp(userID int, body string) (Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	newID := len(dbStruct.Chirps) + 1
	newChirp := Chirp{ID: newID, UserID: userID, Body: body}
	dbStruct.Chirps[newChirp.ID] = newChirp
	err = db.writeDB(dbStruct)
	if err != nil {
		return Chirp{}, err
	}
	return newChirp, nil
}

func (db *DB) GetChirps() (map[int]Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	return dbStruct.Chirps, nil
}

func (db *DB) DeleteChirp(chirpID int) error {
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}
	delete(dbStruct.Chirps, chirpID)
	return nil
}
