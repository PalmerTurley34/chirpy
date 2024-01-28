package database

import "errors"

type User struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

func (db *DB) CreateUser(email, password string) (User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	if db.checkDuplicateEmail(email) {
		return User{}, errors.New("email is already is database")
	}
	newID := len(dbStruct.Users) + 1
	newUser := User{ID: newID, Email: email, Password: password}
	dbStruct.Users[newUser.ID] = newUser
	err = db.writeDB(dbStruct)
	if err != nil {
		return User{}, err
	}
	return newUser, nil
}

func (db *DB) GetUsers() (map[int]User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	return dbStruct.Users, nil
}

func (db *DB) UpdateUser(userID int, email, password string) (User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	_, ok := dbStruct.Users[userID]
	if !ok {
		return User{}, errors.New("users does not exist")
	}
	newUser := User{ID: userID, Email: email, Password: password}
	dbStruct.Users[userID] = newUser
	err = db.writeDB(dbStruct)
	if err != nil {
		return User{}, err
	}
	return newUser, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	for _, user := range dbStruct.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return User{}, errors.New("User does not exist")
}

func (db *DB) checkDuplicateEmail(email string) bool {
	users, err := db.GetUsers()
	if err != nil {
		return false
	}
	for _, user := range users {
		if user.Email == email {
			return true
		}
	}
	return false
}

func (db *DB) UpgradeUserToPremium(userID int) error {
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}
	user, ok := dbStruct.Users[userID]
	if !ok {
		return errors.New("users does not exist")
	}
	user.IsChirpyRed = true
	dbStruct.Users[userID] = user
	return db.writeDB(dbStruct)
}
