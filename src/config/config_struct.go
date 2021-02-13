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
	Cache struct {
		URI  string `yaml:"uri"`
		Addr string `yaml:"addr"`
		Pass string `yaml:"pass"`
		TTL  int64  `yaml:"ttl"`
	} `yaml:"cache"`
	Session struct {
		Name   string `yaml:"name"`
		Secret string `yaml:"secret"`
		UserID string `yaml:"user_id"`
	} `yaml:"session"`
	Secret struct {
		Geolication string `yaml:"geolication"`
	} `yaml:"secrets"`
	Email struct {
		SMTP string `yaml:"smtp"`
		Port int    `yaml:"port"`
		From string `yaml:"from"`
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
	} `yaml:"email"`
}
