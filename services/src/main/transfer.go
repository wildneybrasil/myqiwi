// main
package main

import (
	"db"
	b64 "encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"ws"

	"golang.org/x/crypto/scrypt"
)

type s_values_item_hdr struct {
	Id     string `json:"id,attr,omitempty"`
	Amount string `json:"amount,attr,omitempty"`
}
type s_values_response_hdr struct {
	Items []s_values_item_hdr `json:"items,omitempty"`
}
type s_transferCredits_request_hdr struct {
	AuthToken string `json:"authToken"`
	Rcpt      string `json:"rcpt"`
	Id        string `json:"id"`
	Session   string `json:"session"`
	Amount    string `json:"amount"`
	Password  string `json:"password"`
}
type s_transferCredits_response_data_hdr struct {
	Session  string                `json:"session,omitempty"`
	Id       string                `json:"id,omitempty"`
	Nominals s_values_response_hdr `json:"nominals,omitempty"`
}
type s_transferCredits_response_hdr struct {
	s_status
	Data s_transferCredits_response_data_hdr `json:"data,omitempty"`
}
type s_transferCredits2_response_data_hdr struct {
	Voucher string `json:"voucher,omitempty"`
	Amount  string `json:"amount,omitempty"`
}
type s_transferCredits2_response_hdr struct {
	s_status
	Data s_transferCredits2_response_data_hdr `json:"data,omitempty"`
}

func transferCredits1(s_transferCredits_request s_transferCredits_request_hdr) (s_transferCredits_response s_transferCredits_response_hdr, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("PANIC - ", r)

			err = fmt.Errorf("panic")
		}
	}()
	s_transferCredits_response = s_transferCredits_response_hdr{}

	dbConn := db.Connect()

	defer dbConn.Close()

	s_rcpt_info, err := db.GetLoginInfoByCel(dbConn, s_transferCredits_request.Rcpt)

	s_login_credentials, err := db.GetAuthToken(dbConn, s_transferCredits_request.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {
		//verifica senha
		dk, err := scrypt.Key([]byte(s_transferCredits_request.Password), []byte(s_login_credentials.PasswordSalt), 16384, 8, 1, 32)
		if err != nil {
			s_transferCredits_response.StatusCode = 403
			s_transferCredits_response.ErrorMessage = "Login/Senha inválido."

			return s_transferCredits_response, nil
		}
		dkb64Encoded := b64.StdEncoding.EncodeToString([]byte(dk))

		fmt.Println(dkb64Encoded)

		if dkb64Encoded != s_login_credentials.Password {
			s_transferCredits_response.StatusCode = 403
			s_transferCredits_response.ErrorMessage = "Login/Senha inválido."

			return s_transferCredits_response, nil
		}
		// fim verifica senha
		transferResponse, err := ws.TransferCredits1(s_login_credentials, s_rcpt_info.Cel, s_rcpt_info.TerminalId, "15695", "1")
		if err != nil {
			s_transferCredits_response.StatusCode = 500
			s_transferCredits_response.ErrorMessage = "Internal server error"
		} else {
			strNominalsValue1 := transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp1
			strNominalsValue3 := transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp3
			V1 := strings.Split(strNominalsValue1, ";")
			V3 := strings.Split(strNominalsValue3, ";")

			for k, _ := range V1 {

				item := s_values_item_hdr{}
				amount := strings.Replace(V1[k], "|", "", -1)
				amountFloat, _ := strconv.ParseFloat(amount, 64)
				amountString := strconv.FormatFloat(amountFloat, 'E', -1, 64)

				item.Amount = amountString
				item.Id = V3[k]

				s_transferCredits_response.Data.Nominals.Items = append(s_transferCredits_response.Data.Nominals.Items, item)
			}
			s_transferCredits_response.Data.Session = transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp4
			s_transferCredits_response.Data.Id = transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Id
		}
	} else {
		s_transferCredits_response.StatusCode = 403
		s_transferCredits_response.ErrorMessage = "Login/Senha inválido"
	}
	return s_transferCredits_response, nil
}
func transferCredits2(s_transferCredits_request s_transferCredits_request_hdr) (s_transferCredits_response s_transferCredits2_response_hdr, err error) {
	//	defer func() {
	//		if r := recover(); r != nil {
	//			fmt.Println("PANIC - ", r)

	//			err = fmt.Errorf("panic")
	//		}
	//	}()

	s_transferCredits_response = s_transferCredits2_response_hdr{}

	dbConn := db.Connect()

	defer dbConn.Close()

	s_login_credentials, err := db.GetAuthToken(dbConn, s_transferCredits_request.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {
		// verifica senha
		dk, err := scrypt.Key([]byte(s_transferCredits_request.Password), []byte(s_login_credentials.PasswordSalt), 16384, 8, 1, 32)
		if err != nil {
			s_transferCredits_response.StatusCode = 403
			s_transferCredits_response.ErrorMessage = "Login/Senha inválido."

			return s_transferCredits_response, nil
		}
		dkb64Encoded := b64.StdEncoding.EncodeToString([]byte(dk))

		fmt.Println(dkb64Encoded)

		if dkb64Encoded != s_login_credentials.Password {
			s_transferCredits_response.StatusCode = 403
			s_transferCredits_response.ErrorMessage = "Login/Senha inválido."

			return s_transferCredits_response, nil
		}
		// fim senha

		s_rcpt, err := db.GetLoginInfoByCel(dbConn, s_transferCredits_request.Rcpt)
		if err != nil {
			s_transferCredits_response.StatusCode = 500
			s_transferCredits_response.ErrorMessage = "Internal server error"
			return s_transferCredits_response, nil
		} else {
			transferResponse, requestXML, responseXML, err := ws.TransferCredits2(s_login_credentials, s_transferCredits_request.Id, s_transferCredits_request.Session, s_rcpt.Cel, s_rcpt.TerminalId, "15695", s_transferCredits_request.Amount)
			if err != nil {
				s_transferCredits_response.StatusCode = 500
				s_transferCredits_response.ErrorMessage = "Internal server error"
				return s_transferCredits_response, nil
			}
			if transferResponse.XMLProvider.XMLPurchaseOnline != nil {
				if transferResponse.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLVoucher != nil {
					fmt.Println(transferResponse.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLVoucher.Code)
					s_transferCredits_response.Data.Voucher = transferResponse.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLVoucher.Code
					s_transferCredits_response.Data.Amount = transferResponse.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLVoucher.Amount
				}
				if transferResponse.XMLProvider.XMLPurchaseOnline.XMLPayment.Result != "0" {
					s_transferCredits_response.StatusCode = 400
					s_transferCredits_response.ErrorMessage = "Erro na transferencia"
					return s_transferCredits_response, nil
				}
				db.InsertPaymentHistory(dbConn, s_login_credentials.Id, "transfer", "15695", s_transferCredits_request, s_transferCredits_response, requestXML, responseXML, 1)

			} else {
				s_transferCredits_response.StatusCode = 500
				s_transferCredits_response.ErrorMessage = "Internal server error"
				return s_transferCredits_response, nil

			}
		}
	} else {
		s_transferCredits_response.StatusCode = 403
		s_transferCredits_response.ErrorMessage = "Login/Senha inválido"
	}
	return s_transferCredits_response, nil
}
