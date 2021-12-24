package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	TelegramToken string

	Owner Owner

	DBPath string `mapstructure:"db_file"`

	ButtonTemplates  ButtonTemplates
	CommandResponses CommandResponses
	QueryResponses   QueryResponses

	MainChat  string
	GuestChat string
}

type Owner struct {
	Name       string
	CreditCard string
	Phone      string
}

type ButtonTemplates struct {
	Registration string `mapstructure:"registration_btn"`
	Cancel       string `mapstructure:"cancel_btn"`
	Location     string `mapstructure:"location_btn"`
	Price        string `mapstructure:"price_btn"`
	Call         string `mapstructure:"call_btn"`
	MainChat     string `mapstructure:"main_chat_btn"`
	GuestChat    string `mapstructure:"guest_chat_btn"`

	ChangePhone string `mapstructure:"change_phone_btn"`
	GetPhone    string `mapstructure:"get_phone_btn"`
	GetLocation string `mapstructure:"get_location_btn"`
}

type CommandResponses struct {
	Start   string `mapstructure:"start"`
	Default string `mapstructure:"default"`
}

type QueryResponses struct {
	Cancel      string `mapstructure:"cancel"`
	NewPhone    string `mapstructure:"new_phone"`
	ChangePhone string `mapstructure:"change_phone"`
	SetPhone    string `mapstructure:"set_phone"`
	Thanks      string `mapstructure:"thanks"`
}

func Init() (*Config, error) {

	viper.AddConfigPath("configs")
	viper.SetConfigName("templates")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("button_templates", &cfg.ButtonTemplates); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("command_responses", &cfg.CommandResponses); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("query_responses", &cfg.QueryResponses); err != nil {
		return nil, err
	}

	if err := parseEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func parseEnv(cfg *Config) error {

	if err := viper.BindEnv("TELEGRAM_APITOKEN"); err != nil {
		return err
	}

	if err := viper.BindEnv("OWNER"); err != nil {
		return err
	}

	if err := viper.BindEnv("OWNER_CREDIT_CARD"); err != nil {
		return err
	}

	if err := viper.BindEnv("OWNER_PHONE"); err != nil {
		return err
	}

	if err := viper.BindEnv("MAIN_CHAT"); err != nil {
		return err
	}

	if err := viper.BindEnv("GUEST_CHAT"); err != nil {
		return err
	}

	cfg.Owner.Name = viper.GetString("OWNER")
	cfg.Owner.CreditCard = viper.GetString("OWNER_CREDIT_CARD")
	cfg.Owner.Phone = viper.GetString("OWNER_PHONE")

	cfg.TelegramToken = viper.GetString("TELEGRAM_APITOKEN")

	cfg.MainChat = viper.GetString("MAIN_CHAT")
	cfg.GuestChat = viper.GetString("GUEST_CHAT")

	return nil
}
