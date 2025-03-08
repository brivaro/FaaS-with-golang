package auth

import (
	"errors"
	"faas/models"
	"faas/repository"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Username string `json:"username"`
	Key      string `json:"key"`
	jwt.RegisteredClaims
}

type AuthService struct {
	jwtKey      []byte
	consumerKey string
	repo        repository.UserRepository
}

func NewAuthService(
	jwtKey string,
	consumerKey string,
	repo repository.UserRepository,
) *AuthService {
	return &AuthService{
		jwtKey:      []byte(jwtKey),
		consumerKey: consumerKey,
		repo:        repo,
	}
}

func (s *AuthService) Register(username, password string) error {
	if strings.ContainsAny(username, ".*>$") {
    return errors.New("username can't contain: '.', '*', '>', '$'") 
	}
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	now := time.Now()
	formattedDate := now.Format("02/01/2006 15:04:05")

	user := models.User{
		Username:  username,
		Password:  string(hashedPassword),
		Role:      "user",
		CreatedAt: formattedDate,
	}
	if user.Username == os.Getenv("ADMIN_USER") {
		user.Role = "admin"
	}

	_, err = s.repo.GetUserByUsername(username)
	if err == nil {
		return errors.New("user already exists")
	}

	err = nil

	err = s.repo.InsertUser(user)

	if err != nil {
		return fmt.Errorf("error trying to insert User:%s", err)
	}

	return nil
}

func (s *AuthService) Login(username, password string) (string, error) {
	if strings.ContainsAny(username, ".*>$") {
        return "", errors.New("username can't contain: '.', '*', '>', '$'") 
	}

	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)

	if err != nil {
		return "", errors.New("invalid credentials")
	}

	claims := Claims{
		Username: user.Username,
		Key:      "faas_jwt_consumer",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "faas-auth-service",
			Subject:   "faas-users",
			Audience:  jwt.ClaimStrings{"faas-users"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Crear Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar y obtener el token codificado como un string
	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) Validate(tokenString string) (*Claims, error) {
	// Create a new instance of Claims
	claims := &Claims{}

	// Parse the token with explicit claims type
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Check expiration
	if claims.ExpiresAt.Compare(time.Now()) <= 0 {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
