package service

import (
	"context"
	"net/http"
	"time"
	"todolist/backend/config"
	"todolist/backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type AuthService interface {
	Login(ctx context.Context, email, password string) (*LoginResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*RefreshResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, &ServiceError{Code: http.StatusUnauthorized, Message: "invalid credentials"}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, &ServiceError{Code: http.StatusUnauthorized, Message: "invalid credentials"}
	}

	accessToken, err := generateToken(user.ID, "access", config.App.JWTAccessTTL)
	if err != nil {
		return nil, err
	}
	refreshToken, err := generateToken(user.ID, "refresh", config.App.JWTRefreshTTL)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    config.App.JWTAccessTTL,
	}, nil
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*RefreshResponse, error) {
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(config.App.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, &ServiceError{Code: http.StatusUnauthorized, Message: "invalid refresh token"}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "refresh" {
		return nil, &ServiceError{Code: http.StatusUnauthorized, Message: "invalid refresh token"}
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, &ServiceError{Code: http.StatusUnauthorized, Message: "invalid refresh token"}
	}
	userID, err := bson.ObjectIDFromHex(userIDStr)
	if err != nil {
		return nil, &ServiceError{Code: http.StatusUnauthorized, Message: "invalid refresh token"}
	}

	accessToken, err := generateToken(userID, "access", config.App.JWTAccessTTL)
	if err != nil {
		return nil, err
	}

	return &RefreshResponse{
		AccessToken: accessToken,
		ExpiresIn:   config.App.JWTAccessTTL,
	}, nil
}

func generateToken(userID bson.ObjectID, tokenType string, ttlSeconds int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.Hex(),
		"type":    tokenType,
		"exp":     time.Now().Add(time.Duration(ttlSeconds) * time.Second).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.App.JWTSecret))
}

