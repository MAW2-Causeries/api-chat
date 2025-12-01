package databases

import (
	"MessagesService/utils"

	"github.com/gocql/gocql"
)

// Session is the global ScyllaDB session
var Session *gocql.Session

// InitDatabases initializes the database connections
func InitDatabases() error {
	host := utils.GetEnv("SCYLLA_HOST", "localhost")
	username := utils.GetEnv("SCYLLA_USER", "scylla")
	password := utils.GetEnv("SCYLLA_PASS", "your-awesome-password")
	keyspace := utils.GetEnv("SCYLLA_KEYSPACE", "messages_service")

	cluster := gocql.NewCluster(host)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: username,
		Password: password,
	}
	cluster.Keyspace = keyspace

	session, err := cluster.CreateSession()
	if err != nil {
		return err
	}

	Session = session

	err = Migrate()
	if err != nil {
		return err
	}
	return nil
}