// main
package main

import (
	"db"
	"fmt"
	"strings"
	"ws"
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
				item.Amount = V1[k]
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
		s_rcpt, err := db.GetLoginInfoByCel(dbConn, s_transferCredits_request.Rcpt)
		if err != nil {
			s_transferCredits_response.StatusCode = 500
			s_transferCredits_response.ErrorMessage = "Internal server error"
		} else {
			transferResponse, err := ws.TransferCredits2(s_login_credentials, s_transferCredits_request.Id, s_transferCredits_request.Session, s_rcpt.Cel, s_rcpt.TerminalId, "	", s_transferCredits_request.Amount)
			if err != nil {
				s_transferCredits_response.StatusCode = 500
				s_transferCredits_response.ErrorMessage = "Internal server error"
			} else {
				if transferResponse.XMLProvider.XMLPurchaseOnline != nil {
					if transferResponse.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLVoucher != nil {
						fmt.Println(transferResponse.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLVoucher.Code)
						s_transferCredits_response.Data.Voucher = transferResponse.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLVoucher.Code
						s_transferCredits_response.Data.Amount = transferResponse.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLVoucher.Amount
					}
					if transferResponse.XMLProvider.XMLPurchaseOnline.XMLPayment.Result != "0" {
						s_transferCredits_response.StatusCode = 400
						s_transferCredits_response.ErrorMessage = "Erro na transferencia"

					}
				} else {
					s_transferCredits_response.StatusCode = 500
					s_transferCredits_response.ErrorMessage = "Internal server error"

				}
			}
		}
	} else {
		s_transferCredits_response.StatusCode = 403
		s_transferCredits_response.ErrorMessage = "Login/Senha inválido"
	}
	return s_transferCredits_response, nil
}
