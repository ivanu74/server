package config

import "github.com/jinzhu/configor"

// Config - struct for wrap configure
type Config struct {
	Tusd struct {
		File_path string
		URL_path  string
		URL_addr  string
	}

	DB struct {
		File_path string
		File_name string `default:"data.db"`
		//	Password  string `required:"true" env:"DBPassword"`
		//	Port      uint   `default:"3306"`
	}

	SYR struct {
		URL_path     string `required:"true"`
		File_path    string `required:"true"`
		Field_form   string
		File_ext     string
		URL_auth     string
		Token_field  string
		Token_header string
	}

	SYR_login    string
	SYR_password string
}

// NewConfig - create access to config
func NewConfig() (*Config, error) {
	cfg := Config{}
	err := configor.Load(&cfg, "./config/config.yaml")
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
