package config

import "os"

type Config struct {
	Port                    string
	ConnectionDBHost        string
	ConnectionDBPort        string
	SigningJwtKey           string
	NatsHost                string
	NatsPort                string
	NatsUser                string
	NatsPass                string
	BlockUserCommandSubject string
	BlockUserReplySubject   string
}

func NewConfig() *Config {
	return &Config{
		Port:                    os.Getenv("CONNECTION_SERVICE_PORT"),
		ConnectionDBHost:        os.Getenv("CONNECTION_DB_HOST"),
		ConnectionDBPort:        os.Getenv("CONNECTION_DB_PORT"),
		SigningJwtKey:           os.Getenv("SIGNING_JWT_KEY"),
		NatsHost:                os.Getenv("NATS_HOST"),
		NatsPort:                os.Getenv("NATS_PORT"),
		NatsUser:                os.Getenv("NATS_USER"),
		NatsPass:                os.Getenv("NATS_PASS"),
		BlockUserCommandSubject: os.Getenv("CREATE_ORDER_COMMAND_SUBJECT"),
		BlockUserReplySubject:   os.Getenv("CREATE_ORDER_REPLY_SUBJECT"),
	}
}
