package mail

import (
	"gopkg.in/gomail.v2"
)

var dialer *gomail.Dialer

func InitDialer() {
	dialer = gomail.NewDialer("zeus.bendi.cf", 1025, "admin", "")
}

func SendMessage(m *gomail.Message) error {
	return dialer.DialAndSend(m)
}
