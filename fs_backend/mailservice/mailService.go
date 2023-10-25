package mailservice

import (
	"fs_backend/models"
	"log"
	"net/smtp"
	"os"
)

type MailService struct {
	smtpHost string
	smtpPort string
}

func (ms *MailService) Initialize() {
	ms.smtpHost = os.Getenv("SMTP_HOST")
	ms.smtpPort = os.Getenv("SMTP_PORT")
	if len(ms.smtpHost) == 0 || len(ms.smtpPort) == 0 {
		log.Fatal("Needed a SMTP Host and Port for mail")
	}
}

func (ms MailService) SendAccountCreationMail(account models.OwnerAccount) {
	from := "accounts@fs_rbac.io"
	to := []string{account.Email}
	message := []byte("From: accounts@fs_rbac.io\r\n" +
		"To: " + account.Email + "\r\n" +
		"Subject: Account Created\r\n\r\n" +
		"Hello " + account.Name + "\nYour account has been created successfully.\r\n")
	auth := smtp.PlainAuth("", from, "", ms.smtpHost)
	err := smtp.SendMail(ms.smtpHost+":"+ms.smtpPort, auth, from, to, message)
	if err != nil {
		log.Default().Println(err)
		return
	}
}

func (ms MailService) SendLoginMail(account models.OwnerAccount) {
	from := "accounts@fs_rbac.io"
	to := []string{account.Email}
	message := []byte("From: accounts@fs_rbac.io\r\n" +
		"To: " + account.Email + "\r\n" +
		"Subject: Login Successful\r\n\r\n" +
		"Hello " + account.Name + "\nYour account has been logged in.\r\n")
	auth := smtp.PlainAuth("", from, "", ms.smtpHost)
	err := smtp.SendMail(ms.smtpHost+":"+ms.smtpPort, auth, from, to, message)
	if err != nil {
		log.Default().Println(err)
		return
	}
}
