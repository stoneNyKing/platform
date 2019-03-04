package sender

import (
	"encoding/base64"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"

	"github.com/libra9z/log4go"
)

func encodeRFC2047(String string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{String, ""}
	return strings.Trim(addr.String(), " <>")
}

func SendMail(tomail string, title string, body string) {
	log := log4go.Global

	smtpServer := "smtp.163.com"
	auth := smtp.PlainAuth(
		"",
		"alrammsg@163.com",
		"alrammsg12345",
		smtpServer,
	)

	from := mail.Address{"ihealth", "alrammsg@163.com"}
	to := mail.Address{tomail, tomail}

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = encodeRFC2047(title)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		smtpServer+":25",
		auth,
		from.Address,
		[]string{to.Address},
		[]byte(message),
	)
	log.Info("SendMail:", tomail, title, body)
	if err != nil {
		log.Error("SendMail:", tomail, title, body, err)
	}
}
