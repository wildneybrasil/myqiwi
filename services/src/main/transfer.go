// main
package main

import (
	"db"
	"fmt"
	"ws"
)

type s_values_item_hdr struct {
	Id     string `json:"id,attr,omitempty"`
	Amount string `json:"amount,attr,omitempty"`
	Text   string `json:"text,attr,omitempty"`
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
	if err != nil || s_rcpt_info == nil {
		s_transferCredits_response.Status = "failed"
		s_transferCredits_response.StatusCode = 404
		s_transferCredits_response.ErrorMessage = "Telefone não cadastrado no QIWI"
		return s_transferCredits_response, nil
	}

	s_login_credentials, err := db.GetAuthToken(dbConn, s_transferCredits_request.AuthToken)

	if !CheckPassword(s_transferCredits_request.Password, s_login_credentials.Password, s_login_credentials.PasswordSalt) {
		s_transferCredits_response.StatusCode = 403
		s_transferCredits_response.ErrorMessage = "Login/Senha inválido."

		return s_transferCredits_response, nil
	}

	if err == nil && s_login_credentials.Id > 0 {
		transferResponse, err := ws.TransferCredits1(s_login_credentials, s_rcpt_info.Email, s_transferCredits_request.Amount)
		if err != nil {
			s_transferCredits_response.StatusCode = 500
			s_transferCredits_response.ErrorMessage = "Internal server error"
		}
		if transferResponse.Result == "0" {
			if transferResponse.XMLPerson.CreditTransfer != nil {
				if transferResponse.XMLPerson.CreditTransfer.Result == "0" {
					s_transferCredits_response.StatusCode = 0
				} else {
					s_transferCredits_response.StatusCode = 500
					s_transferCredits_response.ErrorMessage = transferResponse.XMLPerson.CreditTransfer.ResultDescription
				}
			} else {
				s_transferCredits_response.StatusCode = 500
				s_transferCredits_response.ErrorMessage = "Internal server error"

			}
		} else {
			s_transferCredits_response.StatusCode = 500
			s_transferCredits_response.ErrorMessage = "Internal server error"
		}
	} else {
		s_transferCredits_response.StatusCode = 403
		s_transferCredits_response.ErrorMessage = "Login/Senha inválido"
	}
	return s_transferCredits_response, nil
}
