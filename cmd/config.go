package cmd

type Config struct {
	Host      string `mapstructure:"host"`
	Port      string `mapstructure:"port"`
	Localhost bool
	Logrus
	Platforms []string `mapstructure:"platforms"`
	Proxy     []string `mapstructure:"proxy"`
}

type Logrus struct {
	LogLvl int    `mapstructure:"log_level"`
	ToFile bool   `mapstructure:"to_file"`
	ToJson bool   `mapstructure:"to_json"`
	LogDir string `mapstructure:"log_dir"`
}
