// main
package main

import (
	"db"
	"fmt"
	"strings"
	"ws"
)

type s_boletoResponse_data_hdr struct {
	Cedente    string `json:"cendente,omitempty"`
	Expiration string `json:"validate,omitempty"`
	Amount     string `json:"amount,omitempty"`
	Diff       string `json:"diff,omitempty"`
	Flag       string `json:"flag,omitempty"`
	Session    string `json:"session,omitempty"`
	Id         string `json:"id,omitempty"`
}
type s_boletoResponse_hdr struct {
	s_status
	Data s_boletoResponse_data_hdr `json:"data,omitempty"`
}
type s_boletoInfo_request_hdr struct {
	AuthToken string `json:"authToken"`
	Boleto    string `json:"boleto"`
	Scanned   string `json:"scanned"`
}

func getBoletoInfo(s_boletoInfo_request s_boletoInfo_request_hdr) (s_boletoResponse s_boletoResponse_hdr, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("PANIC - ", r)

			err = fmt.Errorf("panic")
		}
	}()

	dbConn := db.Connect()

	defer dbConn.Close()

	s_login_credentials, err := db.GetAuthToken(dbConn, s_boletoInfo_request.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {
		transferResponse, err := ws.GetBoletoInfo(s_login_credentials, s_boletoInfo_request.Boleto, s_login_credentials.Cel, s_boletoInfo_request.Scanned)
		if err != nil {
			s_boletoResponse.StatusCode = 500
			s_boletoResponse.ErrorMessage = "Internal server error"
		} else {
			boletoInfo := transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp2
			if boletoInfo != "" {
				boletoMetadata := strings.Split(boletoInfo, "|")

				if len(boletoMetadata) == 2 {
					s_boletoResponse.Data.Cedente = boletoMetadata[0]
					s_boletoResponse.Data.Amount = boletoMetadata[1]
				}
				if len(boletoMetadata) == 4 {
					s_boletoResponse.Data.Cedente = boletoMetadata[0]
					s_boletoResponse.Data.Expiration = boletoMetadata[1]
					s_boletoResponse.Data.Amount = boletoMetadata[2]
					s_boletoResponse.Data.Diff = boletoMetadata[3]
				}
				if len(boletoMetadata) == 3 {
					s_boletoResponse.Data.Cedente = boletoMetadata[0]
					s_boletoResponse.Data.Expiration = boletoMetadata[1]
					s_boletoResponse.Data.Amount = boletoMetadata[2]
				}
			}
			disp3 := transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp3

			if disp3 != "" {
				flags := strings.Split(disp3, "|")
				errorMessage := ""

				if len(flags) > 1 {
					s_boletoResponse.Data.Flag = flags[0]
					errorMessage = flags[1]
				} else {
					s_boletoResponse.Data.Flag = disp3
				}
				if s_boletoResponse.Data.Flag == "false" {
					s_boletoResponse.Status = "failed"
					s_boletoResponse.StatusCode = 400
					s_boletoResponse.ErrorMessage = errorMessage
				}
			}
			s_boletoResponse.Data.Session = transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp1
			s_boletoResponse.Data.Id = transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Id
		}
	} else {
		s_boletoResponse.StatusCode = 403
		s_boletoResponse.ErrorMessage = "Login/Senha inválido"
	}
	return s_boletoResponse, nil
}
