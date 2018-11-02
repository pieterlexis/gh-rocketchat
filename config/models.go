package config

type Hook struct {
	Name        string `yaml:"name"`
	Destination string `yaml:"destination"`
	Endpoint    string `yaml:"endpoint"`
	Secret      string `yaml:"secret"`
}

type Config struct {
	ListenAddress string `yaml:"listenaddress"`
	Hooks []Hook `yaml:"hooks"`
}
