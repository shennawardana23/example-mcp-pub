package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shennawardana23/example-mcp-pub/internal/app/converter"
	"github.com/shennawardana23/example-mcp-pub/internal/app/entity"
	"github.com/shennawardana23/example-mcp-pub/internal/app/model"
	"github.com/shennawardana23/example-mcp-pub/internal/app/repository"
	"github.com/shennawardana23/example-mcp-pub/internal/config"
	"golang.org/x/crypto/bcrypt"
)

// AuthService defines the interface for authentication business logic
type AuthService interface{
	Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error)
	ValidateToken(tokenString string) (*Claims, error)
	GetUserByUsername(ctx context.Context, username string) (*model.UserResponse, error)
}

// Claims represents JWT claims
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// authService implements AuthService
type authService struct {
	userRepo  repository.UserRepository
	converter *converter.UserConverter
	jwtConfig *config.JWTConfig
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, jwtConfig *config.JWTConfig) AuthService {
	return &authService{
		userRepo:  userRepo,
		converter: converter.NewUserConverter(),
		jwtConfig: jwtConfig,
	}
}

// Login authenticates a user and returns JWT tokens
func (s *authService) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	// Get user by username
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.Active {
		return nil, fmt.Errorf("user account is inactive")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate access token
	accessToken, err := s.generateToken(user, s.jwtConfig.AccessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := s.generateToken(user, s.jwtConfig.RefreshTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.jwtConfig.AccessTokenTTL.Seconds()),
	}, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *authService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtConfig.SecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// GetUserByUsername retrieves a user by username
func (s *authService) GetUserByUsername(ctx context.Context, username string) (*model.UserResponse, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return s.converter.ToResponse(user), nil
}

// generateToken generates a JWT token for the user
func (s *authService) generateToken(user *entity.User, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtConfig.SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}
