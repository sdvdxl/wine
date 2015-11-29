package mail

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/sdvdxl/wine/sources/util/log"
	"gopkg.in/gomail.v1"
	"html/template"
)

type mail struct {
	SmtpServer     string
	Port           int
	SenderEmail    string
	SenderName     string
	SenderPassword string
	ToEmail        string
	Subject        string
	Message        string
}

var (
	m      = mail{SmtpServer: "smtp.exmail.qq.com", Port: 465, SenderEmail: "support@uke.life", SenderName: "support", SenderPassword: "abc1235566"}
	mailer *gomail.Mailer
)

func SendEmail(toPersons []string, subject string, mailTemplate string, data interface{}) error {
	log.Logger.Debug("will send mail...")

	errPerson := make([]string, 0, 10)

	for _, v := range toPersons {
		msg := gomail.NewMessage()
		msg.SetAddressHeader("From", m.SenderEmail, m.SenderName)
		msg.SetAddressHeader("To", v, "")
		msg.SetHeader("Subject", subject)
		tmpl, err := template.ParseFiles("mailtemplates/" + mailTemplate)
		if err != nil {
			return err
		}

		var contents bytes.Buffer
		tmpl.Execute(&contents, data)
		msg.SetBody("text/html", contents.String())

		err = mailer.Send(msg)
		if err != nil {
			errPerson = append(errPerson, v)
			log.Logger.Error("error occured when send email to :%v, err:%v", v, err)
		}
	}

	if len(errPerson) == 0 {
		return nil
	}

	return errors.New(fmt.Sprintf("发送失败的收件人:%v", errPerson))
}

func init() {

	mailer = gomail.NewMailer(m.SmtpServer, m.SenderEmail, m.SenderPassword, m.Port)

}
