package validator

import "regexp"

var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type validator struct {
	Errors map[string]string
}

func New() *validator { return &validator{Errors: make(map[string]string)} }

func (v *validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *validator) Error(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *validator) Check(ok bool, key, message string) {
	if !ok {
		v.Error(key, message)
	}
}

func MatchesRX(s string, rx *regexp.Regexp) bool {
	return rx.MatchString(s)
}

func HasListVal[T comparable](list []T, val T) bool {

	for _, v := range list {
		if v == val {
			return true
		}
	}

	return false
}

func IsSet[T comparable](slice []T) bool {
	UniqueVals := make(map[T]bool)

	for _, v := range slice {
		UniqueVals[v] = true
	}

	return len(UniqueVals) == len(slice)
}
