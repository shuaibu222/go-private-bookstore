package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shuaibu222/go-bookstore/config"
	"github.com/shuaibu222/go-bookstore/models"
	"golang.org/x/crypto/bcrypt"
)

type LoginCredentials struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	UserID   string `json:"user_id"`
	jwt.RegisteredClaims
}

type Cookie struct {
	Name string `json:"cookie"`
}

var jwtSecretKey []byte

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	config, err := config.LoadConfig()
	if err != nil {
		log.Println("Error while loading envs: ", err)
	}

	jwtSecretKey = []byte(config.JWTSecret)

	var cred LoginCredentials
	var user models.UserProfile

	json.NewDecoder(r.Body).Decode(&cred)
	enteredPassword := cred.Password

	user = models.GetUserByUsername(cred.UserName)
	if user.Username == "" {
		http.Error(w, "user is does not exist", http.StatusUnauthorized)
	} else {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(enteredPassword))
		if err != nil {
			log.Println("Failed to compare: ", err)
			http.Error(w, "You entered the wrong password!. Please try again", http.StatusUnauthorized)
			return
		}

		// start generating token journey

		// token payload
		expirationTime := time.Now().Add(5 * time.Minute)
		claims := &Claims{
			Username: cred.UserName,
			UserID:   fmt.Sprint(user.ID),
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}

		// token header
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Lastly generate the token - token signature
		tokenString, err := token.SignedString(jwtSecretKey)
		if err != nil {
			log.Println("Error creating JWT string: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    tokenString,
			Expires:  expirationTime,
			HttpOnly: true,
		})

		cok := Cookie{
			Name: tokenString,
		}

		tkn, _ := json.Marshal(cok)
		w.Write(tkn)
	}

}

func Refresh(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tknStr := c.Value
	claims := &Claims{}

	// token validation
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
		return jwtSecretKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// set the new token as the users token
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// immediately clear the token cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Expires: time.Now(),
	})

	json.NewEncoder(w).Encode("Successfully logged out")
}
