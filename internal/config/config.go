package config

import "github.com/spf13/viper"

type Config struct {
	DB      string `mapstructure:"DB"`
	TEST_DB string `mapstructure:"TEST_DB"`
}

func LoadConfig(path string) (Config, error) {
	var config Config
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	return config, nil
}
