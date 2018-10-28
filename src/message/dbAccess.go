package message

import (
	"Utils"
	"github.com/gocql/gocql"
	"strconv"
	"time"
)

var Session *gocql.Session

func init() {
	var err error
	cluster := gocql.NewCluster(Utils.GetConfig(Utils.DATABASE_HOST_KEY))
	cluster.Port, err = strconv.Atoi(Utils.GetConfig(Utils.DATABASE_PORT_KEY))
	if err != nil {
		panic(err)
	}
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: Utils.GetConfig(Utils.DATABASE_USERNAME_KEY),
		Password: Utils.GetConfig(Utils.DATABASE_PASSWORD_KEY)}

	cluster.Timeout = 20 * time.Second
	Session, err = cluster.CreateSession()
	if err != nil {
		panic(err)
	}
}
