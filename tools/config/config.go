package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func InitEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}
func InitValidation() validator.Validate {
	return *validator.New()
}
