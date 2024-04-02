package mailx

import (
	"context"
	"fmt"
	"net/smtp"
	"onlineCLoud/internel/app/config"
	"strings"
	"sync"
	"time"

	"github.com/jordan-wright/email"
)

type Emailx struct {
	p *email.Pool
}

var (
	Email Emailx
	once  sync.Once
	DEBUG = false
)

func Init() func() {
	once.Do(func() {
		var cfg = config.C.Email
		DNS := fmt.Sprintf("%v:%v", cfg.Host, cfg.Port)

		p, err := email.NewPool(
			DNS,
			10,
			smtp.PlainAuth("", cfg.UserName, cfg.Password, cfg.Host),
		)
		if err != nil {
			panic(err)
		}

		Email.p = p
	})
	return Email.Clear
}

func (mail Emailx) Clear() {
	mail.p.Close()
	mail.p = nil
}

func (mail Emailx) SendMsgWithText(ctx context.Context, dest string, subject string, msg string) error {
	e := email.NewEmail()
	build := strings.Builder{}
	build.WriteString("online cloud<")
	build.WriteString(config.C.Email.UserName + ">")
	e.From = build.String()

	e.To = []string{dest}
	e.Subject = subject
	e.Text = []byte(msg)
	return mail.p.Send(e, 10*time.Second)
}
func (mail Emailx) SendMsgwithHtml(ctx context.Context, dest string, subject string, msg string) error {
	fmt.Printf("%v", mail.p)
	e := email.NewEmail()
	build := strings.Builder{}
	build.WriteString("online cloud<")
	build.WriteString(config.C.Email.UserName + ">")
	e.From = build.String()
	e.To = []string{dest}
	e.Subject = subject
	e.HTML = []byte(msg)
	return mail.p.Send(e, 10*time.Second)
}
