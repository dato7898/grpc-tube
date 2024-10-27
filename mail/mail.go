package mail

import (
	"log"
	"net/smtp"
)

func Send(to []string, body string) {
	from := "grpctube@gmail.com"
	pass := "Password123_"

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, to, []byte(body))

	if err != nil {
		log.Fatalf("smtp error: %s", err)
	}
}
