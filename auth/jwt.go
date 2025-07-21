package auth

import (
	"fmt"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func IsValidPassword(password string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9]{6,}$`)
	return re.MatchString(password)
}
func GenerateJWT(id int, username string, roleSlice []string) string {

	type Claims struct {
		Username string
		Role     []string
		Id       int
		jwt.StandardClaims
	}
	var jwtkey = []byte("secret-key")
	if len(roleSlice) == 0 {
		roleSlice = []string{}
	}

	expireTime := time.Now().Add(time.Minute * 5)
	claims := &Claims{
		Username: username,
		Role:     roleSlice,
		Id:       id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtkey)
	if err != nil {
		fmt.Printf("reding error: %v", err)

	}
	return tokenString
}
