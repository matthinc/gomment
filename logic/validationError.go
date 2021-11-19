package logic

import "fmt"

type validationTypeT string

const (
	validationTypeRequired validationTypeT = "required"
	validationTypeLength   validationTypeT = "length"
	validationTypeSymbol   validationTypeT = "symbol"
)

type ValidationError struct {
	FieldName      string
	ValidationType validationTypeT
	Info           string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("'%s' validation error for field '%s': %s", e.ValidationType, e.FieldName, e.Info)
}

func ValidationErrorRequired(fieldName string) ValidationError {
	return ValidationError{
		FieldName:      fieldName,
		ValidationType: validationTypeRequired,
		Info:           "",
	}
}

func ValidationErrorLength(fieldName string, min uint, max uint) ValidationError {
	return ValidationError{
		FieldName:      fieldName,
		ValidationType: validationTypeLength,
		Info:           fmt.Sprintf("%d,%d", min, max),
	}
}

func ValidationErrorSymbol(fieldName string, symbols string) ValidationError {
	return ValidationError{
		FieldName:      fieldName,
		ValidationType: validationTypeSymbol,
		Info:           symbols,
	}
}
