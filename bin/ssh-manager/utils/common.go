package utils

type Server struct {
	HostName string `mapstructure:"hostname"`
	IP       string `mapstructure:"ip"`
	KeyAuth  bool   `mapstructure:"keyAuth"`
	Alias    string `mapstructure:"alias"`
}

type Env struct {
	Name    string   `mapstructure:"name"`
	Servers []Server `mapstructure:"servers"`
}

type Group struct {
	Name        string `mapstructure:"name"`
	User        string `mapstructure:"user"`
	Environment []Env  `mapstructure:"environment"`
}

type Config struct {
	Groups []Group `mapstructure:"groups"`
}
