package utils

import (
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"

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

func isFloat64(fl validator.FieldLevel) bool {
	if fl.Field().Kind() != reflect.Float64 {
		return false
	}
	return true
}

func isNullableFloat64(fl validator.FieldLevel) bool {
	// 対象値がnullの場合Kindは何故かreflect.Interfaceになる
	return (fl.Field().Kind() == reflect.Interface && fl.Field().IsNil()) || isFloat64(fl)
}

func isBoolean(fl validator.FieldLevel) bool {
	if fl.Field().Kind() != reflect.Bool {
		return false
	}
	return true
}

func isString(fl validator.FieldLevel) bool {
	return fl.Field().Kind() == reflect.String
}

func isNullableString(fl validator.FieldLevel) bool {
	// 対象値がnullの場合Kindは何故かreflect.Interfaceになる
	return (fl.Field().Kind() == reflect.Interface && fl.Field().IsNil()) || isString(fl)
}

func isNotOnlyWhitespaces(fl validator.FieldLevel) bool {
	if fl.Field().Kind() != reflect.String {
		return false
	}
	fieldString := fl.Field().String()
	return strings.TrimSpace(fieldString) != ""
}

var (
	validatorInstance *validator.Validate
	once              sync.Once
)

// 独自バリデーションルールを加えたvalidator.Validateを返す
//
// 独自バリデーションルール:
// - is_integer: 整数かどうかをチェック
// - is_float64: float64かどうかをチェック
// - is_nullable_float64: float64またはnullであるかどうかをチェック
// - is_boolean: 真偽値かどうかをチェック
// - is_string: 文字列かどうかをチェック
// - is_nullable_string: 文字列またはnullであるかどうかをチェック
// - not_consists_of_whitespaces: 空白文字のみで構成されていないかどうかをチェック
func GetValidator() *validator.Validate {
	once.Do(func() {
		validatorInstance = validator.New()
		validatorInstance.RegisterValidation("is_integer", isInteger)
		validatorInstance.RegisterValidation("is_float64", isFloat64)
		validatorInstance.RegisterValidation("is_nullable_float64", isNullableFloat64)
		validatorInstance.RegisterValidation("is_boolean", isBoolean)
		validatorInstance.RegisterValidation("is_string", isString)
		validatorInstance.RegisterValidation("is_nullable_string", isNullableString)
		validatorInstance.RegisterValidation("not_only_whitespaces", isNotOnlyWhitespaces)
	})
	return validatorInstance
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
