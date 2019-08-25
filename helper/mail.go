package helper

import (
	"bytes"
	"html/template"
	"fmt"
	"net/smtp"
	"strings"
	"crypto/tls"
	"strconv"
)

const (
	MailServerAddress 	= "smtp.privateemail.com"
	MailServerPort 		= 465
	MailVerifyAddress 	= "verify@istatribe.org"
	MailVerifyPassword 	= "s6GttK336y8C"
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

type (
	MailRequest struct {
		from 		string
		pw			string
		to		    []string
		subject  	string
		body     	string
	}
)

func NewMailRequest(to []string, subject string) *MailRequest {
	return &MailRequest{
		from: MailVerifyAddress,
		pw: MailVerifyPassword,
		to: to,
		subject: subject,
	}
}

func (r *MailRequest) parseTemplate(fileName string, data interface{}) error {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return err
	}
	r.body = buffer.String()
	return nil
}

func (r *MailRequest) sendMail() error {
	body := fmt.Sprintf("From: %s\r\n", r.from)
	body += "To: " + strings.Join(r.to,";") + "\r\nSubject: " + r.subject + "\r\n" + MIME + "\r\n" + r.body

	SMTP := MailServerAddress + ":" + strconv.Itoa(MailServerPort)
	//build an auth
	auth := smtp.PlainAuth("", MailVerifyAddress, MailVerifyPassword, MailServerAddress)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         MailServerAddress,
	}
	conn, err_mail := tls.Dial("tcp", SMTP, tlsconfig)
	if err_mail != nil {
		return err_mail
	}

	client, err_mail := smtp.NewClient(conn, MailServerAddress)
	if err_mail != nil {
		return err_mail
	}

	// step 1: Use Auth
	if err_mail = client.Auth(auth); err_mail != nil {
		return err_mail
	}

	// step 2: add all from and to
	if err_mail = client.Mail(MailVerifyAddress); err_mail != nil {
		return err_mail
	}
	for _, k := range r.to {
		if err_mail = client.Rcpt(k); err_mail != nil {
			return err_mail
		}
	}

	// Data
	w, err_mail := client.Data()
	if err_mail != nil {
		return err_mail
	}
	_, err_mail = w.Write([]byte(body))
	if err_mail != nil {
		return err_mail
	}

	err_mail = w.Close()
	if err_mail != nil {
		return err_mail
	}

	return client.Quit()
}

func (r *MailRequest) Send(templateName string, items interface{}) error {
	err := r.parseTemplate("./email_tpl/" + templateName + ".html", items)
	if err != nil {
		return err
	}
	return r.sendMail()
}