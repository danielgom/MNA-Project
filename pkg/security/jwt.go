package security

import (
	"fmt"
	"strconv"
	"time"

	"MNA-project/pkg/config"
	"github.com/golang-jwt/jwt"
)

// GenerateTokenWithExp generates a JWT with an expiration of 1 hour (exp time comes from the config).
func GenerateTokenWithExp(userID int64) (string, *time.Time, error) {
	jwtConfig := config.LoadConfig().JWT

	currentTime := time.Now().Local()
	expirationDate := currentTime.Add(time.Second * time.Duration(jwtConfig.Expiration))

	claims := jwt.StandardClaims{
		ExpiresAt: expirationDate.Unix(),
		IssuedAt:  currentTime.Unix(),
		Issuer:    "SysPet-API",
		Subject:   strconv.Itoa(int(userID)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedToken, err := token.SignedString([]byte(jwtConfig.Key))

	if err != nil {
		return "", nil, fmt.Errorf("could not generate JWT %w please try again", err)
	}

	return signedToken, &expirationDate, nil
}
