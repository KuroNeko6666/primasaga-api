package helper

import (
	"log"

	"github.com/KuroNeko6666/prima-api/email"
)

func SendEmailVerification(target string, token string) error {

	toMails := []string{target}

	subject := "Subject: Verifikasi Akun Superapp!\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := "<html><body>Verifikasi Akun Anda <a href='http://192.168.100.107:8000/api/auth/verified?token=" + token + "' >Klik Disini</a></body></html>"
	body := []byte(subject + mime + msg)

	r := email.Send(body, toMails)
	// handling the errors
	if r != nil {
		log.Fatal(r)
	}

	return nil
}
