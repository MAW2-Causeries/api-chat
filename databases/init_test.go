package databases

import (
	"errors"
	"testing"

	"github.com/gocql/gocql"
	"github.com/stretchr/testify/assert"
)

func TestInitDatabases(t *testing.T) {
	firstSession := &gocql.Session{}
	secondSession := &gocql.Session{}
	var createSessionCalls int
	oldCreateCluster := createCluster
	oldCreateClusterSession := createClusterSession
	oldCloseSession := closeSession
	oldSessionQuery := sessionQuery
	oldExecQuery := execQuery
	oldRunMigrations := runMigrations
	t.Cleanup(func() {
		createCluster = oldCreateCluster
		createClusterSession = oldCreateClusterSession
		closeSession = oldCloseSession
		sessionQuery = oldSessionQuery
		execQuery = oldExecQuery
		runMigrations = oldRunMigrations
	})

	createCluster = func(hosts ...string) *gocql.ClusterConfig {
		return &gocql.ClusterConfig{Hosts: hosts}
	}
	createClusterSession = func(_ *gocql.ClusterConfig) (*gocql.Session, error) {
		createSessionCalls++
		if createSessionCalls == 1 {
			return firstSession, nil
		}
		return secondSession, nil
	}
	closeSession = func(_ *gocql.Session) {}
	sessionQuery = func(_ *gocql.Session, _ string, _ ...interface{}) *gocql.Query {
		return &gocql.Query{}
	}
	execQuery = func(query *gocql.Query) error {
		assert.NotNil(t, query)
		return nil
	}
	runMigrations = func() error {
		return nil
	}

	Session = nil
	err := InitDatabases()

	assert.NoError(t, err)
	assert.Equal(t, 2, createSessionCalls)
	assert.Same(t, secondSession, Session)
}

func TestInitDatabasesReturnsFirstSessionError(t *testing.T) {
	expectedErr := errors.New("create session failed")
	oldCreateCluster := createCluster
	oldCreateClusterSession := createClusterSession
	t.Cleanup(func() {
		createCluster = oldCreateCluster
		createClusterSession = oldCreateClusterSession
	})

	createCluster = func(hosts ...string) *gocql.ClusterConfig {
		return &gocql.ClusterConfig{Hosts: hosts}
	}
	createClusterSession = func(_ *gocql.ClusterConfig) (*gocql.Session, error) {
		return nil, expectedErr
	}

	err := InitDatabases()

	assert.ErrorIs(t, err, expectedErr)
}

func TestInitDatabasesReturnsKeyspaceCreationError(t *testing.T) {
	expectedErr := errors.New("create keyspace failed")
	oldCreateCluster := createCluster
	oldCreateClusterSession := createClusterSession
	oldCloseSession := closeSession
	oldSessionQuery := sessionQuery
	oldExecQuery := execQuery
	t.Cleanup(func() {
		createCluster = oldCreateCluster
		createClusterSession = oldCreateClusterSession
		closeSession = oldCloseSession
		sessionQuery = oldSessionQuery
		execQuery = oldExecQuery
	})

	createCluster = func(hosts ...string) *gocql.ClusterConfig {
		return &gocql.ClusterConfig{Hosts: hosts}
	}
	createClusterSession = func(_ *gocql.ClusterConfig) (*gocql.Session, error) {
		return &gocql.Session{}, nil
	}
	closeSession = func(_ *gocql.Session) {}
	sessionQuery = func(_ *gocql.Session, _ string, _ ...interface{}) *gocql.Query {
		return &gocql.Query{}
	}
	execQuery = func(_ *gocql.Query) error { return expectedErr }

	err := InitDatabases()

	assert.ErrorIs(t, err, expectedErr)
}

func TestInitDatabasesReturnsSecondSessionError(t *testing.T) {
	expectedErr := errors.New("second session failed")
	var createSessionCalls int
	oldCreateCluster := createCluster
	oldCreateClusterSession := createClusterSession
	oldCloseSession := closeSession
	oldSessionQuery := sessionQuery
	oldExecQuery := execQuery
	t.Cleanup(func() {
		createCluster = oldCreateCluster
		createClusterSession = oldCreateClusterSession
		closeSession = oldCloseSession
		sessionQuery = oldSessionQuery
		execQuery = oldExecQuery
	})

	createCluster = func(hosts ...string) *gocql.ClusterConfig {
		return &gocql.ClusterConfig{Hosts: hosts}
	}
	createClusterSession = func(_ *gocql.ClusterConfig) (*gocql.Session, error) {
		createSessionCalls++
		if createSessionCalls == 1 {
			return &gocql.Session{}, nil
		}
		return nil, expectedErr
	}
	closeSession = func(_ *gocql.Session) {}
	sessionQuery = func(_ *gocql.Session, _ string, _ ...interface{}) *gocql.Query {
		return &gocql.Query{}
	}
	execQuery = func(_ *gocql.Query) error {
		return nil
	}

	err := InitDatabases()

	assert.ErrorIs(t, err, expectedErr)
}

func TestInitDatabasesReturnsMigrationError(t *testing.T) {
	expectedErr := errors.New("migration failed")
	var createSessionCalls int
	oldCreateCluster := createCluster
	oldCreateClusterSession := createClusterSession
	oldCloseSession := closeSession
	oldSessionQuery := sessionQuery
	oldExecQuery := execQuery
	oldRunMigrations := runMigrations
	t.Cleanup(func() {
		createCluster = oldCreateCluster
		createClusterSession = oldCreateClusterSession
		closeSession = oldCloseSession
		sessionQuery = oldSessionQuery
		execQuery = oldExecQuery
		runMigrations = oldRunMigrations
	})

	createCluster = func(hosts ...string) *gocql.ClusterConfig {
		return &gocql.ClusterConfig{Hosts: hosts}
	}
	createClusterSession = func(_ *gocql.ClusterConfig) (*gocql.Session, error) {
		createSessionCalls++
		return &gocql.Session{}, nil
	}
	closeSession = func(_ *gocql.Session) {}
	sessionQuery = func(_ *gocql.Session, _ string, _ ...interface{}) *gocql.Query {
		return &gocql.Query{}
	}
	execQuery = func(_ *gocql.Query) error {
		return nil
	}
	runMigrations = func() error {
		return expectedErr
	}

	err := InitDatabases()

	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 2, createSessionCalls)
}
