package auth

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
)

type Claims struct {
	Username string
	Role     []string
	Id       int
	jwt.StandardClaims
}

var jwtkey = []byte("secret-key")

type contextKey string

const UserIDKey contextKey = "userID"

func IsValidPassword(password string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9]{6,}$`)
	return re.MatchString(password)
}

func GenerateJWT(id int, username string, roleSlice []string) string {

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
func NormalJwtmiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// پاسخ به preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return

		}

		// توکن چک کردن
		authHeader := r.Header.Get("Authorization")
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtkey, nil
		})
		fmt.Println("token is: ", tokenStr)
		if err != nil || !tkn.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		fmt.Println("claim: ", claims)
		next(w, r)
	}
}
func AdminJwtMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// پاسخ به preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return

		}

		// توکن چک کردن
		authHeader := r.Header.Get("Authorization")
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtkey, nil
		})
		fmt.Println("token is: ", tokenStr)
		if err != nil || !tkn.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		fmt.Println("claim: ", claims)
		valid := false
		for _, v := range claims.Role {
			if v == "1" {
				valid = true

			}

		}
		if !valid {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		rdb := redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})

		blocked, err := rdb.Get(tokenStr).Result()
		if err != redis.Nil && blocked == "blocked" {
			http.Error(w, "Token is blocked", http.StatusUnauthorized)
			return
		}

		// اگه همه چیز اوکی بود، بره سراغ هندلر اصلی
		next(w, r)
	}
}
func StudentJwtMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// پاسخ به preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return

		}

		// توکن چک کردن
		authHeader := r.Header.Get("Authorization")
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtkey, nil
		})
		fmt.Println("token is: ", tokenStr)
		if err != nil || !tkn.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		valid := false
		for _, v := range claims.Role {
			if v == "2" {
				valid = true

			}

		}
		if !valid {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		rdb := redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})

		blocked, err := rdb.Get(tokenStr).Result()
		if err != redis.Nil && blocked == "blocked" {
			http.Error(w, "Token is blocked", http.StatusUnauthorized)
			return
		}
		// اگه همه چیز اوکی بود، بره سراغ هندلر اصلی
		ctx := context.WithValue(r.Context(), UserIDKey, claims.Id)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
func ProfessorjwtMiddleware3(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// پاسخ به preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return

		}

		// توکن چک کردن
		authHeader := r.Header.Get("Authorization")
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtkey, nil
		})
		fmt.Println("token is: ", tokenStr)
		if err != nil || !tkn.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		fmt.Println("claimedrole: ", claims.Role)
		valid := false
		for _, v := range claims.Role {
			if v == "3" {
				valid = true

			}

		}
		if !valid {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		rdb := redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})

		blocked, err := rdb.Get(tokenStr).Result()
		if err != redis.Nil && blocked == "blocked" {
			http.Error(w, "Token is blocked", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, claims.Id)
		next.ServeHTTP(w, r.WithContext(ctx))

	}
}
func CheackOriginMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// پاسخ به preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return

		}

		next(w, r)
	}
}
