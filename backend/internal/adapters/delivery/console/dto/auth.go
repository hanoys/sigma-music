package dto

import (
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strings"
)

type SignUpUserDTO struct {
	Name     string
	Email    string
	Phone    string
	Password string
	County   string
}

func InputSignUpUserDTO(dto *SignUpUserDTO) error {
	fmt.Print("Name: ")
	fmt.Scan(&dto.Name)

	fmt.Print("Password: ")
	fmt.Scan(&dto.Password)

	fmt.Print("Email: ")
	fmt.Scan(&dto.Email)

	_, err := mail.ParseAddress(dto.Email)
	if err != nil {
		return errors.New("incorrect email")
	}

	fmt.Print("Phone: ")
	fmt.Scan(&dto.Phone)
	e164Regex := `^\+[1-9]\d{1,14}$`
	re := regexp.MustCompile(e164Regex)
	dto.Phone = strings.ReplaceAll(dto.Phone, " ", "")

	if re.Find([]byte(dto.Phone)) == nil {
		return errors.New("invalid phone number format")
	}

	fmt.Print("Country: ")
	fmt.Scan(&dto.County)

	return nil
}

type SignUpMusicianDTO struct {
	Name        string
	Email       string
	Password    string
	County      string
	Description string
}

func InputSignUpMusicianDTO(dto *SignUpMusicianDTO) error {
	fmt.Print("Name: ")
	fmt.Scan(&dto.Name)

	fmt.Print("Password: ")
	fmt.Scan(&dto.Password)

	fmt.Print("Email: ")
	fmt.Scan(&dto.Email)

	_, err := mail.ParseAddress(dto.Email)
	if err != nil {
		return errors.New("incorrect email")
	}

	fmt.Print("Country: ")
	fmt.Scan(&dto.County)

	fmt.Print("Description: ")
	fmt.Scan(&dto.Description)

	return nil
}

type LogInDTO struct {
	Name     string
	Password string
}

func InputLogInDTO(dto *LogInDTO) {
	fmt.Print("Name: ")
	fmt.Scan(&dto.Name)

	fmt.Print("Password: ")
	fmt.Scan(&dto.Password)
}
