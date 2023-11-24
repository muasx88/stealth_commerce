package config

var Config config

type config struct {
	PORT string `json:"PORT"`

	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     int    `mapstructure:"DB_PORT"`
	DBUsername string `mapstructure:"DB_USERNAME"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`

	DBMaxConnections         int `mapstructure:"DB_MAX_CONNECTIONS"`
	DBMaxIdleConnections     int `mapstructure:"DB_IDLE_CONNETIONS"`
	DBMaxLifeTimeConnections int `mapstructure:"MAX_LIFETIME_CONNECTIONS"`

	LOG_LEVEL string `mapstructure:"LOG_LEVEL"`
	JWTKey    string `mapstructure:"JWT_KEY"`
}
