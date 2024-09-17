package data

import (
	"chatin/auth"
	"chatin/database"
	usersQueries "chatin/queries/users"
)

type userData struct{}

func NewUserData() *userData {
	return &userData{}
}

func (dt *userData) FindById(id int) {
	var db = database.OpenConnection()
	var result struct {
		rows chan usersQueries.Find
		err  chan error
	}
	result.rows = make(chan usersQueries.Find)
	result.err = make(chan error)
	go func() {
		defer close(result.rows)
		defer close(result.err)
		var statement, err = db.Prepare("SELECT first_name, last_name, email, username FROM users WHERE id = ?")
		if err != nil {
			result.err <- err
			return
		}
		var user usersQueries.Find
		var row = statement.QueryRow(&id)
		err = row.Scan(&user.FirstName, &user.LastName, &user.Email, &user.Username)
		if err != nil {
			result.err <- err
			return
		}
		user.Code, err = auth.EncodeUserID(uint64(id))
		if err != nil {
			result.err <- err
			return
		}

		result.rows <- user
	}()

	select {
	case user := <-result.rows:
		return user, nil
	case err := <-result.err:
		return usersQueries.Find{}, err
	}
}
