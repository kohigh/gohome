package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("field %v is invalid with error %v", v.Field, v.Err)
}

var programError = errors.New("program err has occurred")
var invalidLen = errors.New("not equal to len")
var invalidDueToRegexp = errors.New("regexp failed")
var invalidIn = errors.New("not in slice")
var invalidMin = errors.New("less than min")
var invalidMax = errors.New("more than max")

func validateLen(lenSymbols string, val string) (bool, error) {
	lenInt, err := strconv.Atoi(lenSymbols)
	if err != nil {
		return false, fmt.Errorf("string to int conversion err: %v; %v", err, programError)
	}

	if lenInt != len(val) {
		return false, invalidLen
	}

	return true, nil
}
func validateRegexp(reg, val string) (bool, error) {
	matched, err := regexp.MatchString(reg, val)
	if err != nil {
		return false, fmt.Errorf("regexp program err %v; %v", err, programError)
	}
	if !matched {
		return false, invalidDueToRegexp
	}

	return true, nil
}
func validateIn(inStr string, val string) (bool, error) {
	inSlice := strings.Split(inStr, ",")

	for _, elem := range inSlice {
		if elem == val {
			return true, nil
		}
	}
	return false, invalidIn
}
func validateMin(minSymbols string, val int) (bool, error) {
	minInt, err := strconv.Atoi(minSymbols)
	if err != nil {
		return false, fmt.Errorf("string to int conversion err: %v; %v", err, programError)
	}

	if val <= minInt {
		return false, invalidMin
	}
	return true, nil
}
func validateMax(maxSymbols string, val int) (bool, error) {
	maxInt, err := strconv.Atoi(maxSymbols)
	if err != nil {
		return false, fmt.Errorf("string to int conversion err: %v; %v", err, programError)
	}

	if val >= maxInt {
		return false, invalidMax
	}

	return true, nil
}

type ValidationErrors []ValidationError

var errChan chan error

func (v ValidationErrors) Error() string {
	var errStr strings.Builder

	for _, err := range v {
		errStr.WriteString(err.Error())
	}

	return errStr.String()
}

func Validate(v interface{}) error {
	errChan = make(chan error, 1)

	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Struct {
		return nil
	}

	t := val.Type()

	var errs ValidationErrors

	go func() {
		defer close(errChan)

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			tag, ok := field.Tag.Lookup("validate")
			if !ok {
				continue
			}

			fv := val.Field(i)

			switch fv.Kind() {
			case reflect.Slice:
				for j := 0; j < fv.Len()-1; j++ {
					validateField(tag, fv.Index(j), field)
				}
			default:
				validateField(tag, fv, field)
			}
		}
	}()

	for err := range errChan {
		if errors.Is(err, programError) {
			return err
		}

		errs = append(errs, err.(ValidationError))
	}

	return errs
}

func validateField(tag string, val reflect.Value, field reflect.StructField) {
	validations := strings.Split(tag, "|")

	for _, validation := range validations {
		validationPair := strings.Split(validation, ":")

		var (
			valid bool
			err   error
		)
		switch validationPair[0] {
		case "len":
			valid, err = validateLen(validationPair[1], val.String())
		case "regexp":
			valid, err = validateRegexp(validationPair[1], val.String())
		case "in":
			var checkVal string
			switch val.Kind() {
			case reflect.Int:
				checkVal = strconv.Itoa(int(val.Int()))
			default:
				checkVal = val.String()
			}
			valid, err = validateIn(validationPair[1], checkVal)
		case "min":
			valid, err = validateMin(validationPair[1], int(val.Int()))
		case "max":
			valid, err = validateMax(validationPair[1], int(val.Int()))
		}

		if !valid && err != nil {
			err = ValidationError{Err: err, Field: field.Name}
		}

		if err != nil {
			errChan <- err
		}
	}
}
