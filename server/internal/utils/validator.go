package utils

import (
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

func isInteger(fl validator.FieldLevel) bool {
	if fl.Field().Kind() != reflect.Float64 {
		return false
	}
	fieldFloat := fl.Field().Float()
	// 10進で文字列変換し、%d+にマッチするかどうかをチェック
	fieldString := strconv.FormatFloat(fieldFloat, 'g', -1, 64)
	matched, err := regexp.MatchString("^\\d+$", fieldString)
	if err != nil {
		return false
	}
	return matched
}

func isBoolean(fl validator.FieldLevel) bool {
	if fl.Field().Kind() != reflect.Bool {
		return false
	}
	return true
}

func GetValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("is_integer", isInteger)
	validate.RegisterValidation("is_boolean", isBoolean)
	return validate
}

func GetFirstValidationErrorTarget(err error) string {
	validationErrors, convErr := err.(validator.ValidationErrors)
	if !convErr {
		log.Printf("failed to convert error to validator.ValidationErrors: %v", err)
		return ""
	}

	for _, e := range validationErrors {
		fieldName := e.Field()
		// フィールド名を小文字のスネークケースに変換
		// 例: "Active" -> "active", "IntervalMin" -> "interval_min"
		return toSnakeCase(fieldName)
	}
	return ""
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
