package config

type Redis struct {
	Address  string `yaml:"Address"`
	Password string `yaml:"Password"`
	DB       int    `yaml:"Port"`
}