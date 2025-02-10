package envconf

import (
	"time"

	"github.com/spf13/viper"
)

type Env struct {
	Port        string `mapstructure:"PORT"`
	Environment string `mapstructure:"ENVIRONMENT"`
	Name        string `mapstructure:"NAME"`
	Debug       bool   `mapstructure:"DEBUG"`

	DSN    string `mapstructure:"DB_POSTGRES_DSN"`
	DBName string `mapstructure:"DB_NAME"`
	DBUrl  string `mapstructure:"DB_URL"`

	JWTSecret           string        `mapstructure:"JWT_SECRET"`
	JwtExpiresIn        int           `mapstructure:"JWT_EXPIRES"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func New() (*Env, error) {
	var env Env

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		return nil, err
	}

	return &env, nil
}
