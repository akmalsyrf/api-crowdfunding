package helper

import (
	"errors"
	"os"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

func APIResponse(message string, code int, status string, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	jsonResponse := Response{
		Meta: meta,
		Data: data,
	}

	return jsonResponse
}

func FormatValidationError(err error) []string {
	var errors []string

	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}

	return errors
}

func ValidateFolderExist(path string) error {
	_, err := os.Stat(path)
	isNotExist := errors.Is(err, os.ErrNotExist)

	if isNotExist {
		err := os.Mkdir(path, os.ModePerm)

		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteFile(path string) error {
	if err := os.Remove(path); err != nil {
		return err
	}
	return nil
}
