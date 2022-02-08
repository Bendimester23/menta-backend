package mail

import (
	"os"

	"gopkg.in/gomail.v2"
)

var dialer *gomail.Dialer

func InitDialer() {
	dialer = gomail.NewDialer(os.Getenv("EMAIL_HOST"), 1025, os.Getenv("EMAIL_USER"), os.Getenv("EMAIL_PASS"))
}

func SendMessage(m *gomail.Message) error {
	return dialer.DialAndSend(m)
}
