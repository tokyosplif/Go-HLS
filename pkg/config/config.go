package config

type Config struct {
	MySQLDSN  string
	RedisAddr string
}

func LoadConfig() *Config {
	return &Config{
		MySQLDSN:  "root:root@tcp(localhost:3306)/test_task",
		RedisAddr: "localhost:6379",
	}
}
