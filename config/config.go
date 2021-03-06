package config

type Config struct {
	ServerHost     string   `mapstructure:"server_host"`
	ServerPort     string   `mapstructure:"server_port"`
	WithReflection bool     `mapstructure:"with_reflection"`
	Localhost      bool     `mapstructure:"enable_localhost"`
	Urls           Urls     `mapstructure:"urls"`
	Logrus         Logrus   `mapstructure:"logrus"`
	Proxy          []string `mapstructure:"proxy"`
}

type Logrus struct {
	LogLvl int    `mapstructure:"log_level"`
	ToFile bool   `mapstructure:"to_file"`
	ToJson bool   `mapstructure:"to_json"`
	LogDir string `mapstructure:"log_dir"`
}

type Urls struct {
	MainUrl     string `mapstructure:"main_url"`
	QueryUrl    string `mapstructure:"query_url"`
	SearchUrl   string `mapstructure:"search_url"`
	ImgUrlTempl string `mapstructure:"img_url_temp"`
}
