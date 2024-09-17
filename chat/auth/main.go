package auth

import (
	"fmt"
	"messenger/config"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/sqids/sqids-go"
)

var jwtKey = []byte(config.NewConfig().SecretKey)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func VerifyToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtKey), nil
	})

	if err != nil {
		fmt.Print(tokenStr)
		return claims, err
	}

	var id float64
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id = claims["user_id"].(float64)
	}

	return &Claims{UserID: strconv.Itoa(int(id))}, nil
}

func DecodeID(id string) (int, error) {
	ids, err := sqids.New(sqids.Options{
		MinLength: 10,
	})

	if err != nil {
		return 0, err
	}

	var decodedId = ids.Decode(id)
	return int(decodedId[0]), nil

}

func EncodeID(id uint64) (string, error) {
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
