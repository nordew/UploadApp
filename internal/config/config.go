package config

import "github.com/spf13/viper"

type ConfigInfo struct {
	MongoDBName   string
	Salt          string
	Secret        string
	MinioHost     string
	MinioUser     string
	MinioPassword string
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
