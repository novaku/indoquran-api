package config

// Struct : configuration structure
type Struct struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		URI    string `yaml:"uri"`
		DBName string `yaml:"dbname"`
	} `yaml:"database"`
}
