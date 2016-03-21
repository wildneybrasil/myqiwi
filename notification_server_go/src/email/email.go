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
		"qiwi",
		"vacaloca69",
		"postfix",
	)
	err = SendMail(
		"postfix:25",
		auth,
		"qiwi@qiwi.com.br",
		[]string{rcpt},
		[]byte(text),
	)

	return err
}
