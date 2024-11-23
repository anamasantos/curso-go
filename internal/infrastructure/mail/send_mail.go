package mail

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

func SendMail() error {
	fmt.Println("Sending email...")

	d := gomail.NewDialer(os.Getenv("EMAIL_SMTP"), 587, os.Getenv("EMAIL_USER"), os.Getenv("EMAIL_PASSWORD"))

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_USER"))
	m.SetHeader("To", "asantosgo24@gmail.com")
	m.SetHeader("Subject", "Hello")
	m.SetBody("text/html", "Hello <b>Ana Marina</b>!")
	return d.DialAndSend(m)
}
