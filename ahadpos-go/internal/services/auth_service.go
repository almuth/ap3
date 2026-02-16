package services

import (
	"ahadpos-go/internal/config"
	"ahadpos-go/internal/models"
	"ahadpos-go/internal/repository"
	"ahadpos-go/pkg/utils"
	"errors"
	"time"
)

type AuthService struct {
	cfg            *config.Config
	userRepository *repository.UserRepository
}

func NewAuthService(cfg *config.Config, userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		cfg:            cfg,
		userRepository: userRepo,
	}
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Name     string `json:"name"`
}

func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	// Find user by username
	user, err := s.userRepository.GetByUsernameAndPassword(req.Username)
	if err != nil {
		return nil, utils.ErrUnauthorized
	}

	// Check password
	if !utils.CheckPassword(user.Password, req.Password) {
		return nil, utils.ErrUnauthorized
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Username, user.Role, s.cfg.JWT.Expiration)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &LoginResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		Name:     user.Name,
	}, nil
}

func (s *AuthService) GetUserByID(userID uint) (*models.User, error) {
	user, err := s.userRepository.GetByID(userID)
	if err != nil {
		return nil, utils.ErrNotFound
	}
	return user, nil
}

func (s *AuthService) CreateInitialAdmin(username, password, name string) error {
	// Check if admin already exists
	if s.userRepository.ExistsByUsername(username) {
		return errors.New("admin user already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	// Create admin user
	user := &models.User{
		Username: username,
		Password: hashedPassword,
		Name:     name,
		Role:     "admin",
		Status:   1,
		UpdatedBy: 0, // System created
	}

	return s.userRepository.Create(user)
}
