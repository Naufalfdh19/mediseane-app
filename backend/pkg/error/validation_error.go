package apperror

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func ExtractValidationError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "gte":
		return "Should be greater than or equal with " + fe.Param()
	case "lte":
		return "Should be less than or equal with " + fe.Param()
	case "max":
		return "Should be less than or equal with " + fe.Param()
	case "min":
		return "Should be greater than or equal with " + fe.Param()
	case "email":
		return fe.Field() + " must be in email format"
	case "datetime":
		return "date is not valid"
	}
	return "Unknown error"
}

func FormatValidatedField() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}

func JsonValidator(c *gin.Context) error {
	jsonData, _ := io.ReadAll(c.Request.Body)
	if !json.Valid(jsonData) {
		return ErrInvalidJSON
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonData))
	return nil
}

func ValidateUsernameAndPassword(username, password string) error {
	isUsernameValid := isUsernameValid(username)
	if !isUsernameValid {
		return ErrUsernameNotValid
	}
	isPasswordValid := IsPasswordValid(password)
	if !isPasswordValid {
		return ErrPasswordNotValid
	}
	return nil
}

func IsPasswordValid(password string) bool {
	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return false
	}
	return true
}

func isUsernameValid(username string) bool {
	isValid := regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(username)
	return isValid
}

func IsPhoneNumberValid(phoneNumber string) bool {
	isValid := regexp.MustCompile(`^08\d{5,13}$`).MatchString(phoneNumber)
	return isValid
}

func CheckActiveDays(days string) error {
	daysMap := map[string]int{
		"monday": 0, "tuesday": 0, "wednesday": 0,
		"thursday": 0, "friday": 0, "saturday": 0,
		"sunday": 0,
	}

	activeDays := strings.Split(days, ",")

	for _, day := range activeDays {
		_, exists := daysMap[day]
		if !exists {
			return ErrInvalidDay
		} else {
			daysMap[day]++
		}

		if daysMap[day] > 1 {
			return ErrDuplicateDay
		}
	}

	return nil
}

func CheckPartnerTime(startHour, endHour string) error {
	isValid := IsTimeValid(startHour)
	if !isValid {
		return ErrInvalidHour
	}
	isValid = IsTimeValid(endHour)
	if !isValid {
		return ErrInvalidHour
	}
	sh, err := time.Parse("15:04", startHour)
	if err != nil {
		return ErrInvalidHour
	}
	eh, err := time.Parse("15:04", endHour)
	if err != nil {
		return ErrInvalidHour
	}

	if sh.After(eh) || sh.Equal(eh){
		return ErrInvalidHour
	}

	return nil
}

func IsTimeValid(hour string) bool {
	_, err := time.Parse(`15:04`, hour)
	if err != nil {
		return false
	}
	activeHour := strings.Split(hour, ".")
	return len(activeHour[0]) >= 2
}

func IsYearFoundedValid(yearFounded string) bool {
	yearFoundedInt, err := strconv.Atoi(yearFounded)
	if err != nil || yearFoundedInt < 0 || yearFoundedInt > time.Now().Year() {
		return false
	}
	return true
}

func CheckCurrentPage(currPage, totalPage, limit int) int {
	if currPage <= 0 {
		currPage = 1
	} else if currPage > totalPage {
		currPage = totalPage
	}
	return currPage
}

func IsOrderCanBeCanceled(orderStatus string) bool {
	StatusThatCanBeCanceled := []string{"Pending", "Processing"}

	for _, status := range StatusThatCanBeCanceled {
		if orderStatus == status {
			return true
		}
	}

	return false
}
