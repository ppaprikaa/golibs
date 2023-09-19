package e_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ppaprikaa/golibs/e"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrapErr(t *testing.T) {
	t.Run("returns error", func(t *testing.T) {
		t.Run("from two non-empty errors", func(t *testing.T) {
			var (
				outer    = errors.New("outer")
				inner    = errors.New("inner")
				expected = fmt.Errorf("%w: %w", outer, inner)
			)

			actual := e.WrapErr(outer, inner)

			assert.Equal(t, expected.Error(), actual.Error())
		})

		t.Run("error chaining", func(t *testing.T) {
			var (
				outer    = e.WrapErr(errors.New("first"), errors.New("second"))
				inner    = e.WrapErr(errors.New("third"), errors.New("fourth"))
				expected = fmt.Errorf("%w: %w", outer, inner)
			)

			require.Error(t, outer)
			require.Error(t, inner)

			actual := e.WrapErr(outer, inner)

			assert.Equal(t, expected.Error(), actual.Error())
		})
	})

	t.Run("returns nil", func(t *testing.T) {
		t.Run("empty error strings", func(t *testing.T) {
			type testCase struct {
				outer string
				inner string
			}

			testCases := []testCase{
				{
					outer: "   ",
					inner: "inner",
				},
				{
					outer: "",
					inner: "inner",
				},
				{
					outer: "outer",
					inner: "       ",
				},
				{
					outer: "outer",
					inner: "",
				},
			}

			for _, tc := range testCases {
				result := e.WrapErr(errors.New(tc.outer), errors.New(tc.inner))

				assert.Nil(t, result)
			}
		})

		t.Run("nil errors", func(t *testing.T) {
			type testCase struct {
				outer error
				inner error
			}

			testCases := []testCase{
				{
					outer: nil,
					inner: errors.New("err"),
				},
				{
					outer: errors.New("err"),
					inner: nil,
				},
				{
					outer: nil,
					inner: nil,
				},
			}

			for _, tc := range testCases {
				result := e.WrapErr(tc.outer, tc.inner)

				assert.Nil(t, result)
			}
		})
	})
}

func TestErrNotEmpty(t *testing.T) {
	t.Run("returns false", func(t *testing.T) {
		type testCase struct {
			name string
			err  error
		}

		testCases := []testCase{
			{
				name: "for nil error",
				err:  nil,
			},
			{
				name: "for errors with empty error strings",
				err:  errors.New(""),
			},
			{
				name: "for errors with space error strings",
				err:  errors.New("   "),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				assert.False(t, e.ErrNotEmpty(tc.err))
			})
		}
	})

	t.Run("returns true", func(t *testing.T) {
		type testCase struct {
			name string
			err  error
		}

		testCases := []testCase{
			{
				name: "for non-empty error",
				err:  errors.New("ERROR"),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				assert.True(t, e.ErrNotEmpty(tc.err))
			})
		}
	})
}

func TestWalk(t *testing.T) {
	t.Run("Walks through", func(t *testing.T) {
		t.Run("chained errors", func(t *testing.T) {
			var (
				errs []error
				err  = e.WrapErr(
					errors.New("outer"),
					fmt.Errorf("%w: %w", errors.New("inner"), errors.New("chained error")),
				)

				expectedErrorStrings = []string{"outer", "inner", "chained error"}
			)

			e.Walk(err, func(err error) {
				if err != nil {
					errs = append(errs, err)
				}
			})

			require.Equal(t, len(expectedErrorStrings), len(errs))

			for i := 0; i < len(expectedErrorStrings); i++ {
				assert.Equal(t, expectedErrorStrings[i], errs[i].Error())
			}
		})
	})
}

func TestGetAllTargetErrs(t *testing.T) {
	t.Run("returns nil", func(t *testing.T) {
		t.Run("for different target and chained errors", func(t *testing.T) {
			var (
				err = fmt.Errorf("%w: %w: %w",
					errors.New("first"),
					errors.New("second"),
					errors.New("third"),
				)
				target           = errors.New("fourth")
				expected []error = nil
			)

			actual := e.GetAllTargetErrs(err, target)

			assert.Equal(t, expected, actual)
		})
	})
}
