package databases

import (
	"errors"
	"strings"
	"testing"

	"github.com/bouk/monkey"
	"github.com/gocql/gocql"
	"github.com/stretchr/testify/assert"
)

func TestMigrate(t *testing.T) {
	defer monkey.UnpatchAll()

	Session = &gocql.Session{}
	queries := map[*gocql.Query]string{}
	var executed []string

	monkey.Patch((*gocql.Session).Query, func(_ *gocql.Session, stmt string, _ ...interface{}) *gocql.Query {
		query := &gocql.Query{}
		queries[query] = stmt
		return query
	})
	monkey.Patch((*gocql.Query).Exec, func(query *gocql.Query) error {
		executed = append(executed, queries[query])
		return nil
	})
	monkey.Patch((*gocql.Query).Scan, func(_ *gocql.Query, dest ...any) error {
		*(dest[0].(*int)) = 0
		return nil
	})

	err := Migrate()

	assert.NoError(t, err)
	assert.Len(t, executed, 3)
	assert.True(t, strings.Contains(executed[0], "CREATE TABLE IF NOT EXISTS migrations"))
	assert.True(t, strings.Contains(executed[1], "CREATE TABLE IF NOT EXISTS messages"))
	assert.True(t, strings.Contains(executed[2], "INSERT INTO migrations"))
}

func TestMigrateSkipsAppliedMigration(t *testing.T) {
	defer monkey.UnpatchAll()

	Session = &gocql.Session{}
	queries := map[*gocql.Query]string{}
	var executed []string

	monkey.Patch((*gocql.Session).Query, func(_ *gocql.Session, stmt string, _ ...interface{}) *gocql.Query {
		query := &gocql.Query{}
		queries[query] = stmt
		return query
	})
	monkey.Patch((*gocql.Query).Exec, func(query *gocql.Query) error {
		executed = append(executed, queries[query])
		return nil
	})
	monkey.Patch((*gocql.Query).Scan, func(_ *gocql.Query, dest ...any) error {
		*(dest[0].(*int)) = 1
		return nil
	})

	err := Migrate()

	assert.NoError(t, err)
	assert.Len(t, executed, 1)
	assert.True(t, strings.Contains(executed[0], "CREATE TABLE IF NOT EXISTS migrations"))
}

func TestMigrateReturnsCountQueryError(t *testing.T) {
	defer monkey.UnpatchAll()

	expectedErr := errors.New("count failed")
	Session = &gocql.Session{}

	monkey.Patch((*gocql.Session).Query, func(_ *gocql.Session, _ string, _ ...interface{}) *gocql.Query {
		return &gocql.Query{}
	})
	monkey.Patch((*gocql.Query).Exec, func(_ *gocql.Query) error {
		return nil
	})
	monkey.Patch((*gocql.Query).Scan, func(_ *gocql.Query, _ ...any) error {
		return expectedErr
	})

	err := Migrate()

	assert.ErrorIs(t, err, expectedErr)
}

func TestMigrateReturnsCreateMigrationsTableError(t *testing.T) {
	defer monkey.UnpatchAll()

	expectedErr := errors.New("create migrations failed")
	Session = &gocql.Session{}

	monkey.Patch((*gocql.Session).Query, func(_ *gocql.Session, _ string, _ ...interface{}) *gocql.Query {
		return &gocql.Query{}
	})
	monkey.Patch((*gocql.Query).Exec, func(_ *gocql.Query) error {
		return expectedErr
	})

	err := Migrate()

	assert.ErrorIs(t, err, expectedErr)
}

func TestMigrateReturnsMigrationExecError(t *testing.T) {
	defer monkey.UnpatchAll()

	expectedErr := errors.New("up failed")
	Session = &gocql.Session{}
	queries := map[*gocql.Query]string{}

	monkey.Patch((*gocql.Session).Query, func(_ *gocql.Session, stmt string, _ ...interface{}) *gocql.Query {
		query := &gocql.Query{}
		queries[query] = stmt
		return query
	})
	monkey.Patch((*gocql.Query).Exec, func(query *gocql.Query) error {
		if strings.Contains(queries[query], "CREATE TABLE IF NOT EXISTS messages") {
			return expectedErr
		}
		return nil
	})
	monkey.Patch((*gocql.Query).Scan, func(_ *gocql.Query, dest ...any) error {
		*(dest[0].(*int)) = 0
		return nil
	})

	err := Migrate()

	assert.ErrorIs(t, err, expectedErr)
}

func TestMigrateReturnsInsertError(t *testing.T) {
	defer monkey.UnpatchAll()

	expectedErr := errors.New("insert failed")
	Session = &gocql.Session{}
	queries := map[*gocql.Query]string{}

	monkey.Patch((*gocql.Session).Query, func(_ *gocql.Session, stmt string, _ ...interface{}) *gocql.Query {
		query := &gocql.Query{}
		queries[query] = stmt
		return query
	})
	monkey.Patch((*gocql.Query).Exec, func(query *gocql.Query) error {
		if strings.Contains(queries[query], "INSERT INTO migrations") {
			return expectedErr
		}
		return nil
	})
	monkey.Patch((*gocql.Query).Scan, func(_ *gocql.Query, dest ...any) error {
		*(dest[0].(*int)) = 0
		return nil
	})

	err := Migrate()

	assert.ErrorIs(t, err, expectedErr)
}

func TestCreateMigrationsTableReturnsExecError(t *testing.T) {
	defer monkey.UnpatchAll()

	expectedErr := errors.New("create table failed")
	Session = &gocql.Session{}

	monkey.Patch((*gocql.Session).Query, func(_ *gocql.Session, _ string, _ ...interface{}) *gocql.Query {
		return &gocql.Query{}
	})
	monkey.Patch((*gocql.Query).Exec, func(_ *gocql.Query) error {
		return expectedErr
	})

	err := createMigrationsTable()

	assert.ErrorIs(t, err, expectedErr)
}
