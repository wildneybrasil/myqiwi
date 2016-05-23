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
		"nao-responda@qiwibr.com",
		"ypMgYnPi9FNpKZ9zwJ-dVg",
		"smtp.mandrillapp.com",
	)
	err = SendMail(
		"smtp.mandrillapp.com:25",
		auth,
		"nao-responda@qiwibr.com",
		[]string{rcpt},
		[]byte(text),
	)
        fmt.Println(err)
	return err
}
