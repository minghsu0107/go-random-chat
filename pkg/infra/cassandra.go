package infra

import (
	"github.com/gocql/gocql"
	"github.com/minghsu0107/go-random-chat/pkg/config"
)

var CassandraSession *gocql.Session

func NewCassandraSession(config *config.Config) (*gocql.Session, error) {
	cluster := gocql.NewCluster(config.Cassandra.Host)
	cluster.Port = config.Cassandra.Port
	cluster.Keyspace = config.Cassandra.Keyspace
	cluster.Consistency = gocql.Quorum
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{
		NumRetries: 3,
	}
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: config.Cassandra.User,
		Password: config.Cassandra.Password,
	}
	cluster.DefaultIdempotence = false
	// number of connections per host
	cluster.NumConns = 3
	CassandraSession, err := cluster.CreateSession()
	return CassandraSession, err
}
