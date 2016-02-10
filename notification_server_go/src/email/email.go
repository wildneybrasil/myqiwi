package email

import (
	"fmt"
)

func Send(rcpt string, text string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic", r)
			err = fmt.Errorf("pkg: %v", r)
		}
	}()

	fmt.Println("Sending email to" + rcpt)

	auth := PlainAuth(
		"",
		"irisescola",
		"iris9920",
		"10.30.2.1",
	)
	err = SendMail(
		"10.30.2.1:25",
		auth,
		"noreply@inventt.com.br",
		[]string{rcpt},
		[]byte(text),
	)

	return err
}
