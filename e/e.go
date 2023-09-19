package e

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type WalkFunc func(error)

// error: "outter: inner"
func WrapErr(outer error, inner error) error {
	if ErrNotEmpty(outer) && ErrNotEmpty(inner) {
		return error(&wrappedError{
			outer: outer,
			inner: inner,
		})

	}

	return nil
}

// error: "prefix: err"
func WrapPrefix(p string, err error) error {
	return WrapErr(errors.New(p), err)
}

func Contains(err error, msg string) bool {
	return len(GetAll(err, msg)) > 0
}

func ContainsType(err error, v error) bool {
	return len(GetAllTypes(err, v)) > 0
}

func ContainsTargetErr(err, target error) bool {
	return len(GetAllTargetErrs(err, target)) > 0
}

func GetAll(errs error, msg string) []error {
	return GetAllTargetErrs(errs, errors.New(msg))
}

func GetAllTypes(err error, v any) []error {
	var res []error

	var vtype string
	if v != nil {
		vtype = reflect.TypeOf(v).String()
	}

	Walk(err, func(err error) {
		var ertype string
		if err != nil {
			ertype = reflect.TypeOf(err).String()
		}

		if ertype == vtype {
			res = append(res, err)
		}
	})

	return res
}

func GetAllTargetErrs(errs, target error) []error {
	var res []error

	Walk(errs, func(err error) {
		if errors.Is(target, err) || target.Error() == err.Error() {
			res = append(res, err)
		}
	})

	return res
}

func ErrNotEmpty(err error) bool {
	var errNotEmpty = err != nil

	if errNotEmpty {
		if len([]rune(strings.TrimSpace(err.Error()))) == 0 {
			errNotEmpty = false
		}
	}

	return errNotEmpty
}

func Walk(err error, cb WalkFunc) {
	if err == nil {
		return
	}

	switch e := err.(type) {
	case *wrappedError:
		cb(e.outer)
		Walk(e.inner, cb)
	case interface{ Unwrap() []error }:
		for _, err := range e.Unwrap() {
			Walk(err, cb)
		}
	case interface{ Unwrap() error }:
		Walk(e.Unwrap(), cb)
	default:
		cb(err)
	}
}

type wrappedError struct {
	outer error
	inner error
}

func (w *wrappedError) Error() string {
	if w.outer != nil && w.inner != nil {
		return fmt.Sprintf("%s: %s", w.outer.Error(), w.inner.Error())
	}

	if w.outer != nil {
		return w.outer.Error()
	}

	if w.inner != nil {
		return w.inner.Error()
	}

	return ""
}

func (w *wrappedError) Unwrap() error {
	if w.inner != nil {
		return w.inner
	}

	return w.outer
}

func New(s string) error {
	return &wrappedError{
		inner: errors.New(s),
	}
}
