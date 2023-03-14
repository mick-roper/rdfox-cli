package config

type simpleConfig struct {
	server   string
	protocol string
	role     string
	password string
	logLevel string
}

func (f simpleConfig) Server() string {
	return f.server
}

func (f simpleConfig) Protocol() string {
	return f.protocol
}

func (f simpleConfig) Role() string {
	return f.role
}

func (f simpleConfig) Password() string {
	return f.password
}

func (f simpleConfig) LogLevel() string {
	return f.logLevel
}

type loader func() Config

func Load(loaders ...loader) Config {
	var cfg simpleConfig

	for _, loader := range loaders {
		x := loader()

		if s := x.LogLevel(); s != "" {
			cfg.logLevel = s
		}

		if s := x.Password(); s != "" {
			cfg.password = s
		}

		if s := x.Protocol(); s != "" {
			cfg.protocol = s
		}

		if s := x.Role(); s != "" {
			cfg.role = s
		}

		if s := x.Server(); s != "" {
			cfg.server = s
		}
	}

	return &cfg
}
