package service

import (
	"errors"
	"fmt"
	"os"
	"time"
	"user-service/helper"
	model "user-service/models"
	repository "user-service/repositories"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo         *repository.UserRepository
	EmailService *EmailService
}

func NewUserService(repo *repository.UserRepository, emailService *EmailService) *UserService {
	return &UserService{
		Repo:         repo,
		EmailService: emailService,
	}
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
		ID:         bson.NewObjectID(),
		Name:       req.Name,
		Email:      req.Email,
		Password:   string(hashedPassword),
		Role:       req.Role,
		IsVerified: false,
		CreatedAt:  time.Now(),
	}

	err = s.Repo.Create(user)
	if err != nil {
		return err
	}

	verificationToken, err := helper.GenerateVerificationToken(
		user.ID.Hex(),
		user.Email,
	)

	if err != nil {
		return err
	}

	err = s.EmailService.SendVerificationEmail(
		user.Email,
		user.Name,
		verificationToken,
	)

	if err != nil {
		fmt.Println("failed to send verification email:", err)
	}

	return nil
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

	if !user.IsVerified {
		return "", errors.New("please verify your email first")
	}

	return token, nil
}

func (s *UserService) GetProfile(userID string) (*model.User, error) {
	return s.Repo.FindByID(userID)
}

func (s *UserService) VerifyEmail(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return errors.New("invalid or expired token")
	}

	claims := token.Claims.(jwt.MapClaims)

	userID := claims["user_id"].(string)

	return s.Repo.VerifyUser(userID)
}

func (s *UserService) UpdateUser(
	targetUserID string,
	requesterID string,
	requesterRole string,
	req model.UpdateUserRequest,
) error {
	targetUser, err := s.Repo.FindByID(targetUserID)
	if err != nil {
		return errors.New("user not found")
	}

	if requesterRole != "admin" && requesterID != targetUserID {
		return errors.New("you can only update your own account")
	}

	if requesterRole == "admin" && targetUser.Role == "admin" && requesterID != targetUserID {
		return errors.New("admin cannot update another admin")
	}

	updateData := bson.M{}

	if req.Name != "" {
		updateData["name"] = req.Name
	}

	if req.Email != "" {
		updateData["email"] = req.Email
	}

	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		updateData["password"] = string(hashedPassword)
	}

	if req.Role != "" {
		if requesterRole != "admin" {
			return errors.New("only admin can update role")
		}

		if req.Role != "student" && req.Role != "teacher" {
			return errors.New("admin can only set role to student or teacher")
		}

		updateData["role"] = req.Role
	}

	if len(updateData) == 0 {
		return errors.New("no data to update")
	}

	return s.Repo.UpdateByID(targetUserID, updateData)
}

func (s *UserService) DeleteUser(
	targetUserID string,
	requesterID string,
	requesterRole string,
) error {
	targetUser, err := s.Repo.FindByID(targetUserID)
	if err != nil {
		return errors.New("user not found")
	}

	if requesterRole != "admin" && requesterID != targetUserID {
		return errors.New("you can only delete your own account")
	}

	if requesterRole == "admin" && targetUser.Role == "admin" && requesterID != targetUserID {
		return errors.New("admin cannot delete another admin")
	}

	return s.Repo.DeleteByID(targetUserID)
}
