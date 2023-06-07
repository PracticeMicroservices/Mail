package service

import (
	"bytes"
	"fmt"
	"html/template"
	entities2 "mail/cmd/entities"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

type MailService interface {
	SendSMTPMessage(msg entities2.Message) error
}

type mailService struct {
	Mail *entities2.Mail
}

func NewMailService() MailService {
	return &mailService{
		Mail: entities2.NewMail(),
	}
}

func (m *mailService) SendSMTPMessage(msg entities2.Message) error {
	fmt.Println("Sending SMTP message")
	if msg.From == "" {
		msg.From = m.Mail.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = m.Mail.FromName
	}

	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return err
	}

	plainMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	server := mail.NewSMTPClient()
	server.Host = m.Mail.Host
	server.Port = m.Mail.Port
	server.Username = m.Mail.Username
	server.Password = m.Mail.Password
	server.Encryption = m.getEncryption(m.Mail.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject)

	email.SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextHTML, formattedMessage)

	if len(msg.Attachments) > 0 {
		for _, x := range msg.Attachments {
			email.AddAttachment(x)
		}
	}

	fmt.Println("Sending email")
	err = email.Send(smtpClient)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (m *mailService) buildHTMLMessage(msg entities2.Message) (string, error) {
	fmt.Println("Building HTML message")
	templateToRender := "./templates/mail.html.gohtml"

	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		fmt.Println("Error parsing template", err)
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	formattedMessage := tpl.String()
	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		fmt.Println("Error inlining CSS")
		return "", err
	}

	return formattedMessage, nil
}

func (m *mailService) buildPlainTextMessage(msg entities2.Message) (string, error) {
	fmt.Println("Building plain text message")
	templateToRender := "./templates/mail.plain.gohtml"

	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	plainMessage := tpl.String()

	return plainMessage, nil
}

func (m *mailService) inlineCSS(s string) (string, error) {
	fmt.Println("Inlining CSS")
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}

func (m *mailService) getEncryption(s string) mail.Encryption {
	fmt.Println("Getting encryption")
	switch s {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none", "":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}
