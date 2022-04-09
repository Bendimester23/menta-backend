package mail

import (
	"os"

	"gopkg.in/gomail.v2"
)

var dialer *gomail.Dialer

func InitDialer() {
	dialer = gomail.NewDialer(os.Getenv(`MAIL_HOST`), 25, "admin", "")
}

func SendMessage(m *gomail.Message) error {
	return dialer.DialAndSend(m)
}
