package main

import (
	"errors"
	"testing"

	"cpnv.ch/messagesservice/databases"
	"github.com/bouk/monkey"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestMainFunction(t *testing.T) {
	defer monkey.UnpatchAll()

	monkey.Patch(databases.InitDatabases, func() error {
		return nil
	})
	monkey.Patch((*echo.Echo).Start, func(_ *echo.Echo, _ string) error {
		return nil
	})

	assert.NotPanics(t, main)
}

func TestMainFunctionPanicsWhenDatabaseInitFails(t *testing.T) {
	defer monkey.UnpatchAll()

	expectedErr := errors.New("init failed")
	monkey.Patch(databases.InitDatabases, func() error {
		return expectedErr
	})

	assert.PanicsWithValue(t, expectedErr, func() {
		main()
	})
}
