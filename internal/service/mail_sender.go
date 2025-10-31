package service

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type MailSenderService struct {
	sendMail  string
	port      int
	protocol  string
	secretKey string
}

func NewMailSenderService(data SendMailDep) *MailSenderService {
	return &MailSenderService{
		sendMail:  data.SendMail,
		port:      data.Port,
		protocol:  data.Protocol,
		secretKey: data.SecretKey,
	}
}

func (s *MailSenderService) SendMessage(to string, title string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.sendMail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", title)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(s.protocol, s.port, s.sendMail, s.secretKey)

	if err := d.DialAndSend(m); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
