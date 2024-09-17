package services

import (
	"chatin/auth"
	"chatin/database"
	"chatin/models"
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type authService struct{}

func NewAuthService() *authService {
	return &authService{}
}

func (s *authService) Login(email string, password string) (models.Authentication, error) {
	var db = database.OpenConnection()
	var result struct {
		id  chan int
		err chan error
	}
	result.id = make(chan int)
	result.err = make(chan error)
	go func() {
		defer close(result.id)
		defer close(result.err)
		var user models.Users
		err := db.QueryRow("SELECT id, password FROM users WHERE email = ?", email).Scan(&user.ID, &user.Password)
		if err != nil {
			if err == sql.ErrNoRows {
				result.err <- fmt.Errorf("unauthenticated")
				return
			}

			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			result.err <- fmt.Errorf("internal server error")
			return
		}

		result.id <- user.ID
		result.err <- nil
	}()

	select {
	case id := <-result.id:
		var authorization, err = auth.GenerateJWT(id)
		if err != nil {
			return models.Authentication{}, err
		}
		return authorization, nil

	case err := <-result.err:
		return models.Authentication{}, err

	}

}
