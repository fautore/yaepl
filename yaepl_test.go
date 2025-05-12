package yaepl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Read(t *testing.T) {
	t.Setenv("TEST_USER", "user")
	t.Setenv("TEST_PASSWORD", "pass")
	type A struct {
		User     string `yaepl:"key:TEST_USER"`
		Password string `yaepl:"key:TEST_PASSWORD"`
	}
	var a A
	Read(&a)
	assert.Equal(t, A{User: "user", Password: "pass"}, a)
}
func Test_ReadRequiredSet(t *testing.T) {
	t.Setenv("TEST_USER", "user")
	type A struct {
		User string `yaepl:"key:TEST_USER;required"`
	}
	var a A
	Read(&a)
	assert.Equal(t, A{User: "user"}, a)
}
func Test_ReadRequiredUnset(t *testing.T) {
	type A struct {
		User string `yaepl:"key:TEST_USER;required"`
	}
	var a A
	err := Read(&a)
	assert.Error(t, err)
}
func Test_ReadExportedField(t *testing.T) {
	t.Setenv("EXPORTED", "exported")
	type A struct {
		Exported    string `yaepl:"key:EXPORTED"`
		notExported string `yaepl:"key:NOT_EXPORTED"`
	}
	var a A
	err := Read(&a)
	assert.NoError(t, err)
	assert.Equal(t, A{Exported: "exported"}, a)
}

// Test for types
func Test_String(t *testing.T) {
	t.Setenv("TEST_USER", "user")
	type A struct {
		User string `yaepl:"key:TEST_USER;required"`
	}
	var a A
	err := Read(&a)
	assert.NoError(t, err)
	assert.Equal(t, A{User: "user"}, a)
}
