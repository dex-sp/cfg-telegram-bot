package config

import (
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	TelegramToken string

	Owner Owner

	UserDBPath string `mapstructure:"user_db_file"`

	ButtonTemplates  ButtonTemplates
	CommandResponses CommandResponses
	QueryResponses   QueryResponses
	Errors           Errors

	MainChat    string
	GuestChat   string
	LocationURL string
}

type Owner struct {
	Name       string
	CreditCard string
	Phone      string
	TelegramID int64
}

type ButtonTemplates struct {
	Registration string `mapstructure:"registration_btn"`
	Cancel       string `mapstructure:"cancel_btn"`
	Location     string `mapstructure:"location_btn"`
	Price        string `mapstructure:"price_btn"`
	Pay          string `mapstructure:"pay_btn"`
	Call         string `mapstructure:"call_btn"`
	MainChat     string `mapstructure:"main_chat_btn"`
	GameRules    string `mapstructure:"rules"`

	ChangePhone string `mapstructure:"change_phone_btn"`
	GetPhone    string `mapstructure:"get_phone_btn"`
	GetLocation string `mapstructure:"get_location_btn"`

	GetPaymentDoc    string `mapstructure:"payment_doc_button"`
	PaymentConfirmed string `mapstructure:"payment_confirmed_button"`
	PaymentDeclined  string `mapstructure:"payment_declined_button"`
}

type CommandResponses struct {
	Start    string `mapstructure:"start"`
	Gameover string `mapstructure:"gameover"`
	Default  string `mapstructure:"default"`
}

type QueryResponses struct {
	FirstRegistration                  string `mapstructure:"first_registration"`
	FirstRegistrationOwnerNotification string `mapstructure:"first_registration_owner_notification"`
	Registration                       string `mapstructure:"registration"`
	Price                              string `mapstructure:"price"`
	Pay                                string `mapstructure:"pay"`
	NewPhone                           string `mapstructure:"new_phone"`
	ChangePhone                        string `mapstructure:"change_phone"`
	SetPhone                           string `mapstructure:"set_phone"`
	Location                           string `mapstructure:"location"`
	Thanks                             string `mapstructure:"thanks"`
	CheckPayment                       string `mapstructure:"check_payment"`
	GetRuleBook                        string `mapstructure:"get_rule_book"`

	OwnerСonfirmedPayment string `mapstructure:"owner_confirmed_payment"`
	OwnerDeclinedPayment  string `mapstructure:"owner_declined_payment"`

	PlayerNotification string `mapstructure:"player_notification"`
	PlayerCallRequest  string `mapstructure:"player_call_request"`
}

type Errors struct {
	InvalidPaymentDocument string `mapstructure:"invalid_payment_doc"`
	DeclinedPayment        string `mapstructure:"declined_payment"`
	NotEnoughRights        string `mapstructure:"not_enough_rights"`
	AlreadyConfirmed       string `mapstructure:"already_confirmed"`
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

	if err := viper.UnmarshalKey("errors", &cfg.Errors); err != nil {
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

	if err := viper.BindEnv("OWNER_TELEGRAM_ID"); err != nil {
		return err
	}

	if err := viper.BindEnv("MAIN_CHAT"); err != nil {
		return err
	}

	if err := viper.BindEnv("GUEST_CHAT"); err != nil {
		return err
	}

	if err := viper.BindEnv("LOCATION"); err != nil {
		return err
	}

	if err := viper.BindEnv("LOCATION_URL"); err != nil {
		return err
	}

	telegramIDstr := viper.GetString("OWNER_TELEGRAM_ID")
	telegramID, err := strconv.ParseInt(telegramIDstr, 10, 64)
	if err != nil {
		return err
	}
	cfg.Owner.TelegramID = telegramID

	cfg.Owner.Name = viper.GetString("OWNER")
	cfg.Owner.CreditCard = viper.GetString("OWNER_CREDIT_CARD")
	cfg.Owner.Phone = viper.GetString("OWNER_PHONE")

	cfg.TelegramToken = viper.GetString("TELEGRAM_APITOKEN")

	cfg.MainChat = viper.GetString("MAIN_CHAT")
	cfg.GuestChat = viper.GetString("GUEST_CHAT")
	cfg.LocationURL = viper.GetString("LOCATION_URL")
	return nil
}
