package tools

import (
	"MedodsTestTask/domain"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

var secret = "my_super_secret"

func GenerateAccessToken(userIp string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)
	expTime := time.Now().Add(10 * time.Minute)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = expTime.Unix()
	claims["client_ip"] = userIp

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ExtractData(tokenString string) (bool, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("cannot parse token")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return false, "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userIP := claims["client_ip"].(string)
		return true, userIP, nil
	}
	return false, "", nil

}

func CheckClientIpHash(ip, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(ip))
	return err == nil
}

func GenerateHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func GenerateRefreshRandHash() (string, error) {
	rawToken := make([]byte, 32)
	if _, err := rand.Read(rawToken); err != nil {
		return "", err
	}
	token := base64.StdEncoding.EncodeToString(rawToken)

	bytes, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	return string(bytes), err
}

func HashClientIP(user *domain.User) error {

	ipHash, err := GenerateHash(user.ClientIp)
	if err != nil {
		return err
	}

	user.ClientIp = ipHash
	return nil
}

func SendEmailWarning(userEmail, oldIP, newIP string) {
	log.Printf("Email sent to %s: IP address changed from %s to %s", userEmail, oldIP, newIP)
}
