package mail

import (
	"crypto/tls"
	"fmt"
	"gohub/pkg/logger"
	"net/smtp"
	"strings"

	configPKG "gohub/pkg/config"

	emailPKG "github.com/jordan-wright/email"
)

// SMTP 实现 email.Driver interface
type SMTP struct{}

// Send 实现 email.Driver interface 的 Send 方法
// 原本返回的是bool，现在是error
func (s *SMTP) Send(email Email, config map[string]string) error {

	e := emailPKG.NewEmail()

	e.From = fmt.Sprintf("%v <%v>", email.From.Name, email.From.Address)
	e.To = email.To
	e.Bcc = email.Bcc
	e.Cc = email.Cc
	e.Subject = email.Subject
	e.Text = email.Text
	e.HTML = email.HTML

	logger.DebugJSON("发送邮件", "发件详情", e)
	if configPKG.GetBool("mail.smtp.tls") {
		err := e.SendWithTLS(
			fmt.Sprintf("%v:%v", config["host"], config["port"]),
			smtp.PlainAuth(
				"",
				config["username"],
				config["password"],
				config["host"],
			),
			&tls.Config{
				ServerName: config["host"],
			},
		)
		if err != nil && !strings.Contains(err.Error(), "short response: \u0000\u0000\u0000\u001a\u0000\u0000\u0000") {
			logger.ErrorString("发送邮件", "发件出错", err.Error())
			return fmt.Errorf(err.Error())
		}
	} else {
		err := e.Send(
			fmt.Sprintf("%v:%v", config["host"], config["port"]),
			smtp.PlainAuth(
				"",
				config["username"],
				config["password"],
				config["host"],
			),
		)
		if err != nil {
			logger.ErrorString("发送邮件", "发件出错", err.Error())
			return fmt.Errorf(err.Error())
		}
	}
	logger.DebugString("发送邮件", "发件成功", "")
	return nil
}
