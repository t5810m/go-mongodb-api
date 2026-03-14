package services

import (
	"context"
	"fmt"
	"time"

	"go-mongodb-api/interfaces"
	"go-mongodb-api/middleware"
	"go-mongodb-api/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userService interfaces.UserService
	jwtSecret   []byte
}

func NewAuthService(userService interfaces.UserService, jwtSecret string) *AuthService {
	return &AuthService{
		userService: userService,
		jwtSecret:   []byte(jwtSecret),
	}
}

// Login authenticates any user by email/password and returns a JWT with their role.
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}
	return s.generateToken(user.ID.Hex(), user.Email, user.Role)
}

// Register creates a new user account with a hashed password.
func (s *AuthService) Register(ctx context.Context, user *models.User) error {
	user.Active = true
	user.CreatedTime = time.Now()
	user.UpdatedTime = time.Now()
	return s.userService.CreateUser(ctx, user)
}

func (s *AuthService) generateToken(id, email, role string) (string, error) {
	claims := &middleware.Claims{
		UserID: id,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
