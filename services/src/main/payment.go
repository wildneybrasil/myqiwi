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

type s_payment_request_hdr struct {
	AuthToken      string `json:"authToken"`
	Rcpt           string `json:"rcpt"`
	Service        string `json:"service"`
	Id             string `json:"id"`
	Session        string `json:"session"`
	Amount         string `json:"amount"`
	Type           string `json:"type"`
	CardNumber     string `json:"cardNumber"`
	SelectedAmount string `json:"selectedAmount"`
	Password       string `json:"password"`
}
type s_payment_response_data_hdr struct {
	Session  string                 `json:"session,omitempty"`
	Id       string                 `json:"id,omitempty"`
	Nominals *s_values_response_hdr `json:"nominals,omitempty"`
}
type s_payment_response_hdr struct {
	s_status
	Data *s_transferCredits_response_data_hdr `json:"data,omitempty"`
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
		if err != nil || transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Status != "0" || transferResponse.XMLProvider.XMLCheckPaymentRequisites.Status != "0" {
			s_boletoResponse.StatusCode = 500
			s_boletoResponse.ErrorMessage = "Internal server error"
			return s_boletoResponse, nil
		}
		if transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Status != "0" || transferResponse.XMLProvider.XMLCheckPaymentRequisites.Status != "0" {
			s_boletoResponse.StatusCode = 400
			s_boletoResponse.ErrorMessage = "Internal server error"
			return s_boletoResponse, nil
		}
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

	} else {
		s_boletoResponse.StatusCode = 403
		s_boletoResponse.ErrorMessage = "Login/Senha inválido"
	}
	return s_boletoResponse, nil
}

func payment1(s_payment_request s_payment_request_hdr) (s_payment_response s_payment_response_hdr, err error) {

	s_payment_response = s_payment_response_hdr{}

	dbConn := db.Connect()

	defer dbConn.Close()

	s_login_credentials, err := db.GetAuthToken(dbConn, s_payment_request.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {
		if !CheckPassword(s_payment_request.Password, s_login_credentials.Password, s_login_credentials.PasswordSalt) {
			s_payment_response.StatusCode = 403
			s_payment_response.ErrorMessage = "Login/Senha inválido."

			return s_payment_response, nil
		}

		if s_payment_request.Type == "" {
			s_payment_response.StatusCode = 500
			s_payment_response.ErrorMessage = "Missing type."

			return s_payment_response, nil
		}

		// telefonia
		if s_payment_request.Type == "telefonia" {
			transferResponse, err := ws.DoPaymentTel1(s_login_credentials, s_payment_request.Rcpt, s_payment_request.Service)

			if err != nil || transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Status != "0" || transferResponse.XMLProvider.XMLCheckPaymentRequisites.Status != "0" {
				s_payment_response.StatusCode = 500
				s_payment_response.ErrorMessage = "Internal server error"
				return s_payment_response, nil
			}
			if transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Status != "0" || transferResponse.XMLProvider.XMLCheckPaymentRequisites.Status != "0" {
				s_payment_response.StatusCode = 400
				s_payment_response.ErrorMessage = "Internal server error"
				return s_payment_response, nil
			}

			strNominalsValue1 := transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp1
			strNominalsValue3 := transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp3
			V1 := strings.Split(strNominalsValue1, ";")
			V3 := strings.Split(strNominalsValue3, ";")

			s_payment_response.Data = &s_transferCredits_response_data_hdr{}
			s_payment_response.Data.Nominals = s_values_response_hdr{}

			for k, _ := range V1 {

				item := s_values_item_hdr{}
				item.Amount = V1[k]
				item.Id = V3[k]

				s_payment_response.Data.Nominals.Items = append(s_payment_response.Data.Nominals.Items, item)
			}
			s_payment_response.Data.Session = transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp4
			s_payment_response.Data.Id = transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Id
		}
		// transporte // NAO FUNCIONAL

		if s_payment_request.Type == "transporte" {
			transferResponse, err := ws.DoPaymentTrans1(s_login_credentials, s_payment_request.CardNumber, s_payment_request.Service)

			if err != nil {
				s_payment_response.StatusCode = 500
				s_payment_response.ErrorMessage = "Internal server error"
				return s_payment_response, nil
			}
			if transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Status != "0" || transferResponse.XMLProvider.XMLCheckPaymentRequisites.Status != "0" {
				s_payment_response.StatusCode = 400
				s_payment_response.ErrorMessage = "Internal server error"
				return s_payment_response, nil
			}

			strNominalsValue1 := transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp1
			strNominalsValue3 := transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp3
			V1 := strings.Split(strNominalsValue1, ";")
			V3 := strings.Split(strNominalsValue3, ";")

			s_payment_response.Data = &s_transferCredits_response_data_hdr{}
			s_payment_response.Data.Nominals = s_values_response_hdr{}
			for k, _ := range V1 {

				item := s_values_item_hdr{}
				item.Amount = V1[k]
				item.Id = V3[k]

				s_payment_response.Data.Nominals.Items = append(s_payment_response.Data.Nominals.Items, item)
			}
			s_payment_response.Data.Session = transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp4
			s_payment_response.Data.Id = transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Id
		}

		// games
		if s_payment_request.Type == "games" {
			transferResponse, err := ws.DoPaymentGames1(s_login_credentials, s_payment_request.CardNumber, s_payment_request.Service)

			if err != nil {
				s_payment_response.StatusCode = 500
				s_payment_response.ErrorMessage = "Internal server error"
				return s_payment_response, nil
			}
			if transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Status != "0" || transferResponse.XMLProvider.XMLCheckPaymentRequisites.Status != "0" {
				s_payment_response.StatusCode = 400
				s_payment_response.ErrorMessage = "Internal server error"
				return s_payment_response, nil
			}

			strNominalsValue1 := transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp1
			strNominalsValue3 := transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp3
			V1 := strings.Split(strNominalsValue1, ";")
			V3 := strings.Split(strNominalsValue3, ";")

			s_payment_response.Data = &s_transferCredits_response_data_hdr{}
			s_payment_response.Data.Nominals = s_values_response_hdr{}
			for k, _ := range V1 {

				item := s_values_item_hdr{}
				item.Amount = V1[k]
				item.Id = V3[k]

				s_payment_response.Data.Nominals.Items = append(s_payment_response.Data.Nominals.Items, item)
			}
			s_payment_response.Data.Session = transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp4
			s_payment_response.Data.Id = transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Id
		}

	} else {
		s_payment_response.StatusCode = 403
		s_payment_response.ErrorMessage = "Login/Senha inválido"
	}
	return s_payment_response, nil
}
func payment2(s_payment_request s_payment_request_hdr) (s_payment_response s_status, err error) {
	if s_payment_request.Type == "" {
		s_payment_response.StatusCode = 500
		s_payment_response.ErrorMessage = "Missing type."

		return s_payment_response, nil
	}

	s_payment_response = s_status{}

	dbConn := db.Connect()

	defer dbConn.Close()

	s_login_credentials, err := db.GetAuthToken(dbConn, s_payment_request.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {
		if !CheckPassword(s_payment_request.Password, s_login_credentials.Password, s_login_credentials.PasswordSalt) {
			s_payment_response.StatusCode = 403
			s_payment_response.ErrorMessage = "Login/Senha inválido."

			return s_payment_response, nil
		}

		session := ""

		transferResponse1, err := ws.DoPaymentTel2(s_login_credentials, s_payment_request.Id, s_payment_request.Rcpt, s_payment_request.Service, s_payment_request.SelectedAmount)
		if err != nil {
			s_payment_response.StatusCode = 500
			s_payment_response.ErrorMessage = "Internal server error"
			return s_payment_response, nil
		} else {
			session = transferResponse1.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp1
			//			pid = transferResponse1.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Id
		}

		// fim verifica senha
		transferResponse2, requestXML, responseXML, err := ws.DoPaymentTel3(s_login_credentials, s_payment_request.Id, session, s_payment_request.Rcpt, s_payment_request.Service, s_payment_request.SelectedAmount)
		if err != nil {
			s_payment_response.StatusCode = 500
			s_payment_response.ErrorMessage = "Internal server error"
			return s_payment_response, nil
		} else {
			//			if transferResponse2.XMLProvider.XMLPurchaseOnline.Result != "0" {
			//				s_payment_response.StatusCode = 500
			//				s_payment_response.ErrorMessage = "Erro código level 1: " + transferResponse2.XMLProvider.XMLPurchaseOnline.Result
			//				return s_payment_response, nil
			//			}
			if transferResponse2.XMLProvider.XMLPurchaseOnline.XMLPayment.Result != "0" {
				s_payment_response.StatusCode = 500
				s_payment_response.ErrorMessage = "Erro código level 2: " + transferResponse2.XMLProvider.XMLPurchaseOnline.XMLPayment.Result
				return s_payment_response, nil
			}
			db.InsertPaymentHistory(dbConn, s_login_credentials.Id, s_payment_request.Type, s_payment_request.Service, s_payment_request, s_payment_response, requestXML, responseXML, 1)

		}
	} else {
		s_payment_response.StatusCode = 403
		s_payment_response.ErrorMessage = "Login/Senha inválido"
	}
	return s_payment_response, nil
}
