package email

import (
	"fmt"
)

func Send( rcpt string , text string ) ( err error ) {
		defer func() {
	        if r := recover(); r != nil {
	            fmt.Println("Panic", r)			
				err = fmt.Errorf("pkg: %v", r)
	        }
	    }()	

        auth := PlainAuth(
                "",
                "irisescola",
                "iris9920",
                "od.hostname.org",
        )
        err = SendMail(
                "od.hostname.org:25",
                auth,
                "noreply@inventt.com.br",
                []string{rcpt},
                []byte(text),
        )
		
        return err;
}