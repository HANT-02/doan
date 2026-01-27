package types

// ConsulConfig cấu hình cho Consul
type ConsulConfig struct {
	Address string   `mapstructure:"address" json:"address" yaml:"address"`
	Token   *string  `mapstructure:"token" json:"token" yaml:"token"`
	Paths   []string `mapstructure:"paths" json:"paths" yaml:"paths"`
}

// RedisConfig cấu hình cho Redis
type RedisConfig struct {
	Host           string `mapstructure:"host"`
	Port           string `mapstructure:"port"`
	Password       string `mapstructure:"password"`
	DB             int    `mapstructure:"db"`
	PoolSize       int    `mapstructure:"pool_size"`
	MinIdleConns   int    `mapstructure:"min_idle_conns"`
	MaxRetries     int    `mapstructure:"max_retries"`
	DialTimeoutMs  int    `mapstructure:"dial_timeout_ms"`
	ReadTimeoutMs  int    `mapstructure:"read_timeout_ms"`
	WriteTimeoutMs int    `mapstructure:"write_timeout_ms"`
}

// PSQLConfig cấu hình cho PostgreSQL
type PSQLConfig struct {
	Driver          string `json:"driver,omitempty" yaml:"driver"`
	User            string `json:"user,omitempty" yaml:"user"`
	Password        string `json:"password,omitempty" yaml:"password"`
	Host            string `json:"host,omitempty" yaml:"host"`
	Port            string `json:"port,omitempty" yaml:"port"`
	Schema          string `json:"schema,omitempty" yaml:"schema"`
	MaxOpenConns    int    `json:"max_open_conns,omitempty" yaml:"max_open_conns" mapstructure:"max_open_conns"`
	MaxIdleConns    int    `json:"max_idle_conns,omitempty" yaml:"max_idle_conns" mapstructure:"max_idle_conns"`
	ConnMaxLifetime string `json:"conn_max_lifetime,omitempty" yaml:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`
}

// JWTConfig cấu hình cho JWT
type JWTConfig struct {
	Secret               string `json:"secret,omitempty" yaml:"secret" mapstructure:"secret"`
	AccessTokenDuration  string `json:"access_token_duration,omitempty" yaml:"access_token_duration" mapstructure:"access_token_duration"`
	RefreshTokenDuration string `json:"refresh_token_duration,omitempty" yaml:"refresh_token_duration" mapstructure:"refresh_token_duration"`
}
