package domain

type NetServerConfig struct {
	Host string
	Port int
}

type RadioLogicConfig struct {
	ResponseRate int `mapstructure:"response_rate"`
}

type Config struct {
	RadioTcpServer NetServerConfig
	RadioLogic     RadioLogicConfig
}
