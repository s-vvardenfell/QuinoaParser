package cmd

type Config struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Telegram `mapstructure:"telegram"`
}

type Telegram struct {
	Token  string `mapstructure:"token"`
	ChatId string `mapstructure:"chat_id"`
}
