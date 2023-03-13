package config

type Config interface {
	Server() string
	Protocol() string
	Role() string
	Password() string
}
