package auth

import (
	"godance/internal/dto"
	"godance/internal/models"
	"godance/internal/repository/user"
	types "godance/internal/type"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *user.Repository
}

func NewService(repo *user.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(req dto.RegistrationRequest) (*dto.AuthResponse, error) {
	foundUser, err := s.repo.FindByEmail(req.Email)

	if foundUser != nil {
		return nil, err
	}
	passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)

	if err != nil {
		return nil, err
	}

	newUser := models.User{
		Email:        req.Email,
		PasswordHash: string(passHash),
		Phone:        req.Phone,
		Role:         string(req.Role),
	}

	err = s.repo.Create(&newUser)

	if err != nil {
		return nil, err
	}

	access, refresh, err := generateTokenPair(newUser.ID, newUser.Role, os.Getenv("JWT_SECRET"))

	return toAuthResponse(access, refresh, *toUserShort(newUser)), nil
}

func (s *Service) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	foundUser, err := s.repo.FindByEmail(req.Email)

	if foundUser == nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, err
	}

	access, refresh, err := generateTokenPair(foundUser.ID, foundUser.Role, os.Getenv("JWT_SECRET"))

	return toAuthResponse(access, refresh, *toUserShort(*foundUser)), nil
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
