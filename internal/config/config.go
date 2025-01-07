package config

import "github.com/spf13/viper"

type Config struct {
	DB               string `mapstructure:"DB"`
	TEST_DB          string `mapstructure:"TEST_DB"`
	EMAIL_FROM       string `mapstructure:"EMAIL_FROM"`
	EMAIL_PASSWORD   string `mapstructure:"EMAIL_PASSWORD"`
	EMAIL_REPORTS    string `mapstructure:"EMAIL_REPORTS"`
	EMAIL_LOGO       string `mapstructure:"EMAIL_LOGO"`
	EMAIL_LINK       string `mapstructure:"EMAIL_LINK"`
	EMAIL_LINK_LABEL string `mapstructure:"EMAIL_LINK_LABEL"`
	STORE_URL        string `mapstructure:"STORE_URL"`
	SSL_CERT         string `mapstructure:"SSL_CERT"`
	SSL_KEY          string `mapstructure:"SSL_KEY"`
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
