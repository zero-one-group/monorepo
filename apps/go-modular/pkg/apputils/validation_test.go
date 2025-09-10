package apputils

import (
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidationErrorsToMap(t *testing.T) {
	v := validator.New()

	t.Run("RequiredField", func(t *testing.T) {
		type S struct {
			Name string `json:"name" validate:"required"`
		}
		s := S{}
		err := v.Struct(s)
		require.Error(t, err)
		m := ValidationErrorsToMap(err, s)
		assert.Equal(t, "The name field is required", m["name"])
	})

	t.Run("RequiredField_PointerObj", func(t *testing.T) {
		type S struct {
			Name string `json:"name" validate:"required"`
		}
		s := &S{}
		err := v.Struct(s)
		require.Error(t, err)
		m := ValidationErrorsToMap(err, s)
		assert.Equal(t, "The name field is required", m["name"])
	})

	t.Run("UUID", func(t *testing.T) {
		type U struct {
			ID string `json:"id" validate:"uuid"`
		}
		u := U{ID: "not-a-uuid"}
		err := v.Struct(u)
		require.Error(t, err)
		m := ValidationErrorsToMap(err, u)
		assert.Equal(t, "Must be a valid UUID", m["id"])
	})

	t.Run("MinConstraint", func(t *testing.T) {
		type M struct {
			Tags []string `json:"tags" validate:"min=2"`
		}
		mobj := M{Tags: []string{"one"}}
		err := v.Struct(mobj)
		require.Error(t, err)
		m := ValidationErrorsToMap(err, mobj)
		assert.Equal(t, "Minimum length is 2", m["tags"])
	})

	t.Run("EqField", func(t *testing.T) {
		type P struct {
			Password string `json:"password" validate:"required"`
			Confirm  string `json:"confirm" validate:"eqfield=Password"`
		}
		p := P{Password: "secret", Confirm: "different"}
		err := v.Struct(p)
		require.Error(t, err)
		m := ValidationErrorsToMap(err, p)
		// fe.Param() is expected to return the struct field name ("Password")
		assert.Equal(t, "Must match Password", m["confirm"])
	})

	t.Run("NonValidatorError", func(t *testing.T) {
		err := fmt.Errorf("unexpected failure")
		m := ValidationErrorsToMap(err, struct{}{})
		assert.Equal(t, "unexpected failure", m["error"])
	})
}
