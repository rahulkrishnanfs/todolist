package auth

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	jwtrequest "github.com/golang-jwt/jwt/v5/request"
)

type Authenticator struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func NewAuthenticator(private *rsa.PrivateKey, public *rsa.PublicKey) *Authenticator {
	return &Authenticator{
		PrivateKey: private,
		PublicKey:  public,
	}
}

// https://pkg.go.dev/github.com/golang-jwt/jwt/v5#example-NewWithClaims-CustomClaimsType
type ApplicationClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func (a *Authenticator) GenerateJWT(username, role string) (string, error) {

	claims := ApplicationClaims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "zaagpro.com",
		}}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(a.PrivateKey)

	if err != nil {
		return "", err
	}
	return ss, nil
}

// https://pkg.go.dev/github.com/golang-jwt/jwt/request#MultiExtractor
func (a *Authenticator) AuthorizeRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := jwtrequest.ParseFromRequest(r, jwtrequest.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (any, error) {
				return a.PublicKey, nil
			}, jwtrequest.WithClaims(&ApplicationClaims{}))
		if err != nil {
			fmt.Println("ERROR In the token ", err)
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return

		}
		if token.Valid {
			ctx := context.WithValue(r.Context(), "user", token.Claims.(*ApplicationClaims).Username)
			slog.Info("TOKEN IS VALID")
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})

}
