package api

import (
	"github.com/Viczdera/bank-gorm/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}

// The function takes a single argument, fl of type validator.FieldLevel, which provides access to the field being validated. Inside the function, it attempts to extract the field's value and assert that it is a string. If the assertion succeeds (meaning the field is indeed a string), it calls util.IsSupportedCurrency(currency), which presumably checks whether the string represents a valid or supported currency code (like "USD", "EUR", etc.). If the field is not a string, or if the currency is not supported, the function returns false, indicating that the validation failed.
