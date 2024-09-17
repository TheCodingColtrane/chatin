package auth

import (
	"chatin/config"
	"chatin/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sqids/sqids-go"
	"golang.org/x/crypto/bcrypt"
)

var configs = config.NewConfig()

func GenerateJWT(id int) (models.Authentication, error) {
	var sk = []byte(configs.SecretKey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	var exp = time.Now().Add(time.Hour * 10).Unix()
	claims["exp"] = exp
	claims["user_id"] = id
	tokenString, err := token.SignedString(sk)
	if err != nil {
		return models.Authentication{}, err
	}
	var authentication = models.Authentication{Token: tokenString, Expiry: exp}
	return authentication, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 16)
	return string(bytes), err
}

func DecodeUserID(id string) (int, error) {
	ids, err := sqids.New(sqids.Options{
		MinLength: 10,
	})

	if err != nil {
		return 0, err
	}

	var decodedId = ids.Decode(id)
	return int(decodedId[0]), nil

}

func EncodeUserID(id uint64) (string, error) {
	ids, err := sqids.New(sqids.Options{
		MinLength: 10,
	})

	if err != nil {
		return "", err
	}

	encodedId, err := ids.Encode([]uint64{id})
	if err != nil {
		return "", err
	}

	return encodedId, nil
}
