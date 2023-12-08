package config

import "github.com/spf13/viper"

type ConfigInfo struct {
	Salt   string
	Secret string

	PGHost     string
	PGPort     string
	PGUser     string
	PGDBName   string
	PGSSLMode  string
	PGPassword string

	MinioHost     string
	MinioPort     string
	MinioUser     string
	MinioPassword string

	Rabbit string

	StripeKey string
}

func NewConfig(name, fileType, path string) (*ConfigInfo, error) {
	viper.SetConfigName(name)
	viper.SetConfigType(fileType)
	viper.AddConfigPath(path)
	viper.ReadInConfig()

	var config ConfigInfo

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
