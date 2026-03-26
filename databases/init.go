package databases

import (
	"cpnv.ch/messagesservice/utils"
	"fmt"

	"github.com/gocql/gocql"
)

// Session is the global ScyllaDB session
var Session *gocql.Session
var createCluster = gocql.NewCluster
var createClusterSession = func(cluster *gocql.ClusterConfig) (*gocql.Session, error) {
	return cluster.CreateSession()
}
var closeSession = func(session *gocql.Session) {
	session.Close()
}
var sessionQuery = func(session *gocql.Session, stmt string, values ...interface{}) *gocql.Query {
	return session.Query(stmt, values...)
}
var execQuery = func(query *gocql.Query) error {
	return query.Exec()
}
var runMigrations = Migrate

// InitDatabases initializes the database connections
func InitDatabases() error {
	host := utils.GetEnv("SCYLLA_HOST", "localhost")
	username := utils.GetEnv("SCYLLA_USER", "scylla")
	password := utils.GetEnv("SCYLLA_PASS", "your-awesome-password")
	keyspace := utils.GetEnv("SCYLLA_KEYSPACE", "messages_service")

	cluster := createCluster(host)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: username,
		Password: password,
	}

	session, err := createClusterSession(cluster)
	if err != nil {
		return err
	}
	defer closeSession(session)

	query := fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {'class': 'NetworkTopologyStrategy', 'replication_factor': 1}`, keyspace)
	if err := execQuery(sessionQuery(session, query)); err != nil {
		return err
	}

	cluster.Keyspace = keyspace
	Session, err = createClusterSession(cluster)
	if err != nil {
		return err
	}

	err = runMigrations()
	if err != nil {
		return err
	}
	return nil
}
