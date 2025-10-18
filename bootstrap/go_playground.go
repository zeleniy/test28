package bootstrap

import (
	"regexp"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func SetUpGoPlayground() {

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("regex", validateRegex)
		v.RegisterValidation("date", validateDate)
	}
}

// Date validator
func validateDate(fl validator.FieldLevel) bool {

	value, ok := fl.Field().Interface().(string)

	if !ok {
		return false
	}

	layout := fl.Param()
	if layout == "" {
		return false
	}

	_, err := time.Parse(layout, value)
	return err == nil
}

// Regex validator
func validateRegex(fl validator.FieldLevel) bool {

	value, ok := fl.Field().Interface().(string)

	if !ok {
		return false
	}

	regexPattern := fl.Param()
	if regexPattern == "" {
		return false
	}

	matched, err := regexp.MatchString(regexPattern, value)
	if err != nil {
		return false
	}

	return matched
}
