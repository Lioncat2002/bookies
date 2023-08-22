package token

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// GenerateToken generates a jwt token and assign a email to it's claims and return it
func GenerateToken(user_id uint) (string, error) {

	token_lifespan, err := strconv.Atoi(os.Getenv("TOKEN_LIFESPAN")) //how long token is valid

	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = strconv.Itoa(int(user_id))
	fmt.Println(user_id)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix() //setting a expiry time
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("API_KEY")))

}

// if the token is valid
func TokenValid(c *gin.Context) error {
	tokenString := ExtractToken(c)
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_KEY")), nil
	})
	if err != nil {
		return err
	}
	return nil
}

// extract token from params of auth bearer token
func ExtractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

// get the id from the token
func ExtractID(c *gin.Context) (uint, error) {

	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			//wrong signing method
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_KEY")), nil //API_KEY has the secret string used for signing up the user
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		fmt.Println(claims["id"])
		id, _ := strconv.ParseUint(claims["id"].(string), 10, 64) //get the id

		return uint(id), nil
	}
	return 0, nil
}
