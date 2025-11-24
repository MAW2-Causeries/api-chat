package databases

import (
	"os"

	"github.com/gocql/gocql"
)

// Session is the global ScyllaDB session
var Session *gocql.Session

// InitDatabases initializes the database connections
func InitDatabases() error {
		getEnv := func(key, def string) string {
		if v, ok := os.LookupEnv(key); ok && v != "" {
			return v
		}
		return def
	}

	host := getEnv("SCYLLA_HOST", "localhost")
	username := getEnv("SCYLLA_USER", "scylla")
	password := getEnv("SCYLLA_PASS", "your-awesome-password")
	keyspace := getEnv("SCYLLA_KEYSPACE", "messages_service")

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