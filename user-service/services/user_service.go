package service

import (
	"errors"
	"time"
	"user-service/helper"
	model "user-service/models"
	repository "user-service/repositories"

	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) Register(req model.RegisterRequest) error {
	if req.Name == "" || req.Email == "" || req.Password == "" || req.Role == "" {
		return errors.New("all fields are required")
	}

	if req.Role != "student" && req.Role != "teacher" && req.Role != "admin" {
		return errors.New("invalid role")
	}

	existingUser, _ := s.Repo.FindByEmail(req.Email)
	if existingUser != nil {
		return errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := model.User{
		ID:        bson.NewObjectID(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Role:      req.Role,
		CreatedAt: time.Now(),
	}

	return s.Repo.Create(user)
}

func (s *UserService) Login(req model.LoginRequest) (string, error) {
	user, err := s.Repo.FindByEmail(req.Email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := helper.GenerateToken(user.ID.Hex(), user.Email, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) GetProfile(userID string) (*model.User, error) {
	return s.Repo.FindByID(userID)
}
