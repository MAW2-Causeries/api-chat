package databases

import (
	"MessagesService/utils"
	"fmt"

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

	session, err := cluster.CreateSession()
	if err != nil {
		return err
	}
	defer session.Close()

	query := fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {'class': 'NetworkTopologyStrategy', 'replication_factor': 1}`, keyspace)
	if err := session.Query(query).Exec(); err != nil {
		return err
	}

	cluster.Keyspace = keyspace
	Session, err = cluster.CreateSession()
	if err != nil {
		return err
	}

	err = Migrate()
	if err != nil {
		return err
	}
	return nil
}