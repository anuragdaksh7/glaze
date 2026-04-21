package config

import "log"

type EmailConfig struct {
	Host        string
	Port        int
	Username    string
	Password    string
	SenderEmail string
	SenderName  string
}

var MailerConfig EmailConfig

func LoadMailerConfig() {
	config, err := LoadConfig(".")
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	MailerConfig = EmailConfig{
		Host:        config.SMTPHost,
		Port:        config.SMTPPort,
		Username:    config.SMTPUsername,
		Password:    config.SMTPPassword,
		SenderEmail: config.SMTPUsername,
		SenderName:  "Team Resourcify",
	}
}
