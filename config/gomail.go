package config

import (
	g "github.com/TrinityKnights/Backend/pkg/gomail"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

func NewGomail(viper *viper.Viper, log *logrus.Logger) *g.ImplGomail {
	dialer := gomail.NewDialer(
		viper.GetString("SMTP_HOST"),
		viper.GetInt("SMTP_PORT"),
		viper.GetString("SMTP_USERNAME"), // fromEmail
		viper.GetString("SMTP_PASSWORD"),
	)

	if _, err := dialer.Dial(); err != nil {
		log.Errorf("failed to connect to SMTP server: %v", err)
	} else {
		log.Info("Successfully connected to SMTP server")
	}

	return g.NewGomail(dialer, viper.GetString("SMTP_USERNAME"))
}
