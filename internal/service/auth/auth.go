package auth

import (
	"errors"
	"os"
	"time"

	"godance/internal/domain"
	"godance/internal/dto"
	"godance/internal/models"
	types "godance/internal/type"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo domain.UserRepository
}

func NewService(repo domain.UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(req dto.RegistrationRequest) (*dto.AuthResponse, error) {
	if _, err := s.repo.FindByEmail(req.Email); err == nil {
		return nil, domain.ErrEmailTaken
	} else if !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, err
	}

	bio := ""
	if req.Bio != nil {
		bio = *req.Bio
	}
	newUser := models.User{
		Email:        req.Email,
		PasswordHash: string(passHash),
		Phone:        req.Phone,
		Role:         string(req.Role),
		Profile: &models.Profile{
			FullName: req.Fullname,
			Bio:      bio,
		},
	}

	if err := s.repo.Create(&newUser); err != nil {
		return nil, err
	}

	access, refresh, err := generateTokenPair(newUser.ID, newUser.Role, os.Getenv("JWT_SECRET"))
	if err != nil {
		return nil, err
	}
	return toAuthResponse(access, refresh, *toUserShort(newUser)), nil
}

func (s *Service) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if errors.Is(err, domain.ErrUserNotFound) {
		return nil, domain.ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return nil, domain.ErrInvalidCredentials
	}

	access, refresh, err := generateTokenPair(user.ID, user.Role, os.Getenv("JWT_SECRET"))
	if err != nil {
		return nil, err
	}
	return toAuthResponse(access, refresh, *toUserShort(*user)), nil
}

func (s *Service) Refresh(refreshToken string) (*dto.AuthResponse, error) {
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrInvalidToken
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "refresh" {
		return nil, domain.ErrInvalidToken
	}
	sub, ok := claims["sub"].(float64)
	if !ok {
		return nil, domain.ErrInvalidToken
	}

	user, err := s.repo.FindByID(uint(sub))
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	access, refresh, err := generateTokenPair(user.ID, user.Role, os.Getenv("JWT_SECRET"))
	if err != nil {
		return nil, err
	}
	return toAuthResponse(access, refresh, *toUserShort(*user)), nil
}

func toAuthResponse(access string, refresh string, req dto.UserShort) *dto.AuthResponse {
	return &dto.AuthResponse{
		Access_token:  access,
		Refresh_token: refresh,
		User:          req,
	}
}

func toUserShort(user models.User) *dto.UserShort {
	return &dto.UserShort{
		ID:       user.ID,
		Email:    user.Email,
		Role:     types.UserRole(user.Role),
		Fullname: user.Profile.FullName,
	}
}

func generateTokenPair(userID uint, role string, secret string) (access string, refresh string, err error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"type": "access",
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	access, err = accessToken.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  userID,
		"type": "refresh",
		"exp":  time.Now().Add(time.Hour * 24 * 35).Unix(),
	})
	refresh, err = refreshToken.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}
