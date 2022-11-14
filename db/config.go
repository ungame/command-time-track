package db

import "fmt"

const (
	defaultUser           = "root"
	defaultPassword       = "root"
	defaultHost           = "localhost"
	defaultPort           = 3306
	defaultDatabase       = "command_time_track"
	mysqlStringConnection = "%s:%s@tcp(%s:%d)/%s?parseTime=true"
)

type config struct {
	user string
	pass string
	host string
	port int
}

func (c *config) Source() string {
	return fmt.Sprintf(mysqlStringConnection, c.user, c.pass, c.host, c.port, defaultDatabase)
}

type Option func(c *config)

func WithUser(user string) Option {
	return func(c *config) {
		c.user = user
	}
}

func WithPassword(pass string) Option {
	return func(c *config) {
		c.pass = pass
	}
}

func WithHost(host string) Option {
	return func(c *config) {
		c.host = host
	}
}

func WithPort(port int) Option {
	return func(c *config) {
		c.port = port
	}
}

var defaultConfig config

func init() {
	defaultConfig = config{
		user: defaultUser,
		pass: defaultPassword,
		host: defaultHost,
		port: defaultPort,
	}
}

func newConfig(opts ...Option) *config {
	cfg := &defaultConfig
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
