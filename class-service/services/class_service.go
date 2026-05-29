package services

import (
	"math/rand"
	"phase3/finalproject/class_service-melvin/models"
	"phase3/finalproject/class_service-melvin/repositories"
	"time"
)

func GenerateJoinCode() string {
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

	rand.Seed(time.Now().UnixNano())

	code := make([]rune, 6)

	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}

	return string(code)
}

func CreateClass(class models.Class) error {
	class.JoinCode = GenerateJoinCode()
	class.CreatedAt = time.Now()

	class.Students = []models.Student{}

	return repositories.CreateClass(class)
}

func JoinClass(joinCode string, studentID string, rawToken string) error {
	class, err := repositories.FindClassByJoinCode(joinCode)
	if err != nil {
		return err
	}

	student := models.Student{
		StudentID: studentID,
		JoinedAt:  time.Now(),
	}

	if rawToken != "" {
		if profile, err := GetUserProfile(rawToken); err == nil {
			if name, ok := profile["name"].(string); ok {
				student.Name = name
			}

			if email, ok := profile["email"].(string); ok {
				student.Email = email
			}
		}
	}

	return repositories.JoinClass(class.ID.Hex(), student)
}
