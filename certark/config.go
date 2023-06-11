package certark

const (
	MODE_DEV  = "dev"
	MODE_PROD = "prod"
)

type Config struct {
	Mode string `yaml:"mode"`
}

var DefaultConfig = Config{
	Mode: "dev",
}

var CurrentConfig = Config{}
