package main

import (
"log"
"net/smtp"
)

func main() {
	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		"@qq.com",
		"",
		"smtp.qq.com:465",
	)
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		"smtp.qq.com:465",
		auth,
		"@qq.com",
		[]string{"@staff.weibo.com"},
		[]byte("This is the email body."),
	)
	if err != nil {
		log.Fatal(err)
	}
}
