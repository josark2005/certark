package certark

const (
	MODE_DEV  = "dev"
	MODE_PROD = "prod"
)

type Config struct {
	Mode string `yaml:"mode"`
	Port int64  `yaml:"port"`
}

var DefaultConfig = Config{
	Mode: "dev",
	Port: 7701,
}

var CurrentConfig = Config{}
