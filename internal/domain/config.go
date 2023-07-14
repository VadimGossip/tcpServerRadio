package domain

type NetServerConfig struct {
	Host string
	Port int
}

type Config struct {
	RadioTcpServer NetServerConfig
}
