package util

import (
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application
// The Values are read by viper from a config file or environment
type Config struct{
	DBDriver  				string 			`mapstructure:"DB_DRIVER"`
	DBSource  				string 			`mapstructure:"DB_SOURCE"`
	ServerAddress  			string 			`mapstructure:"SERVER_ADDRESS"`
	TOKEN_SYMMETRIC_KEY  	string 			`mapstructure:"TOKEN_SYMMETRIC_KEY"`
	ACCESS_TOKEN_DURATION   time.Duration   `mapstructure:"ACCESS_TOKEN_DURATION"`
}

// LoadConfig reads configurations from file or environment
func LoadConfig(path string)(config Config,err error){
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err!=nil{
		return 
	}
	err = viper.Unmarshal(&config)
	return 
	
}