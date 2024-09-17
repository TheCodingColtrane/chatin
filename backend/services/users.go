package services

import (
	"chatin/auth"
	"chatin/database"
	"chatin/models"
	usersQueries "chatin/queries/users"
)

type usersService struct{}

func NewUserService() *usersService {
	return &usersService{}
}

func (serv usersService) Get() ([]models.Users, error) {
	var db = database.OpenConnection()
	rows, err := db.Query("SELECT * FROM users ORDER BY ID DESC")
	if err != nil {
		return []models.Users{}, err
	}
	var user models.Users
	var users = make([]models.Users, 0)
	for rows.Next() {
		var err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Username, &user.Password)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (serv usersService) Create(user models.Users) (models.Authentication, error) {
	var db = database.OpenConnection()
	var sqlResult struct {
		id  chan int64
		err chan error
	}
	sqlResult.id = make(chan int64)
	sqlResult.err = make(chan error)

	go func() {
		defer close(sqlResult.id)
		defer close(sqlResult.err)
		hashedPassword, err := auth.HashPassword(user.Password)
		if err != nil {
			sqlResult.err <- err
			return
		}
		statement, err := db.Prepare("INSERT INTO users values (null, ?, ?, ?, ?, ? )")
		if err != nil {
			sqlResult.err <- err
			return
		}
		defer statement.Close()
		result, err := statement.Exec(&user.FirstName, &user.LastName, &user.Email, &user.Username, &hashedPassword)
		if err != nil {
			sqlResult.err <- err
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			sqlResult.err <- err
			return
		}
		sqlResult.id <- id
		sqlResult.err <- err
	}()
	select {
	case id := <-sqlResult.id:
		var authorization, err = auth.GenerateJWT(int(id))
		if err != nil {
			return authorization, err
		}
		return authorization, nil
	case err := <-sqlResult.err:
		return models.Authentication{}, err
	}

}

func (serv usersService) Find(id int) (usersQueries.Find, error) {
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

func (serv usersService) Update(user models.Users) (bool, error) {
	var db = database.OpenConnection()
	result, err := db.Exec("UPDATE users SET first_name = ?, last_name = ? email = ?, username = ?, password = ? WHERE id = ? ", &user.FirstName, &user.LastName, &user.Email, &user.Username, &user.Password, &user.ID)
	if err != nil {
		return false, err
	}
	affectedRecords, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if affectedRecords > 0 {
		return true, nil
	}

	return false, nil
}

func (serv usersService) Delete(user models.Users) (bool, error) {
	var db = database.OpenConnection()
	result, err := db.Exec("DELETE users WHERE id = ?", &user.ID)
	if err != nil {
		return false, err
	}
	affectedRecords, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if affectedRecords > 0 {
		return true, nil
	}

	return false, nil
}
