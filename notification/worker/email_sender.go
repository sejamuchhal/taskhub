package worker

import (
	"context"
	"time"

	"github.com/mailersend/mailersend-go"
	"github.com/sejamuchhal/taskhub/notification/common"
)

//go:generate mockgen -destination=mocks/mock_email_sender.go -package=mocks . EmailSenderInterface
type EmailSenderInterface interface {
	SendEmail(to, subject, body string) (*mailersend.Response, error)
}

type EmailSender struct {
	MailerSendServer *mailersend.Mailersend
	SenderName       string
	SenderEmail      string
}

func NewEmailSender(conf *common.Config) *EmailSender {
	return &EmailSender{
		MailerSendServer: mailersend.NewMailersend(conf.MailersendAPIKey),
		SenderName:       conf.MailersendSenderName,
		SenderEmail:      conf.MailersendSenderEmail,
	}
}

func (e *EmailSender) SendEmail(to, subject, body string) (*mailersend.Response, error) {

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	from := mailersend.From{
		Name:  e.SenderName,
		Email: e.SenderEmail,
	}

	recipients := []mailersend.Recipient{
		{
			Email: to,
		},
	}

	message := e.MailerSendServer.Email.NewMessage()

	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(subject)
	message.SetText(body)
	message.SetInReplyTo("client-id")

	res, err := e.MailerSendServer.Email.Send(ctx, message)
	if err != nil {
		return nil, err
	}

	return res, nil
}
