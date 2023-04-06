package email

import (
	"fmt"
	"net/smtp"

	"github.com/KuroNeko6666/prima-api/config"
)

func Send(body []byte, toList []string) error {

	auth := smtp.PlainAuth("", config.SMTP_AUTH_EMAIL, config.SMTP_AUTH_PASSWORD, config.SMTP_HOST)
	err := smtp.SendMail(config.SMTP_HOST+":"+config.SMTP_PORT, auth, config.SMTP_SENDER_NAME, toList, body)

	// handling the errors
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
