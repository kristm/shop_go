package config

import "github.com/spf13/viper"

type Config struct {
	DB      string `mapstructure:"DB"`
	TEST_DB string `mapstructure:"TEST_DB"`
}

func LoadConfig(path string) (Config, error) {
	var cfg Config
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return cfg, err
	}

	err = viper.Unmarshal(&cfg)
	return cfg, nil
}
