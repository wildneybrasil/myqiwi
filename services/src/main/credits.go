// main
package main

import (
	"db"
	"fmt"
	"ws"
)

type s_createBill_request_hdr struct {
	AuthToken string `json:"authToken"`
	Amount    string `json:"amount"`
}
type s_getBill_request_hdr struct {
	AuthToken string `json:"authToken"`
	BoletoId  string `json:"boletoId"`
}

type s_createBill_response_data_hdr struct {
	Amount string `json:"amount"`
	Id     string `json:"id,omitempty"`
}

type s_balance_request_hdr struct {
	AuthToken string `json:"authToken"`
}
type s_balance_response_data_hdr struct {
	Balance   string `json:"balance"`
	Overdraft string `json:"overdraft"`
}

type s_balance_response_hdr struct {
	s_status
	Data *s_balance_response_data_hdr `json:"data,omitempty"`
}

type s_createBill_response_hdr struct {
	s_status
	Data *s_createBill_response_data_hdr `json:"data,omitempty"`
}
type s_geBill_image_data_response struct {
	Image string `json:"image,omitempty"`
}
type s_geBill_image_response_hdr struct {
	s_status
	Data *s_geBill_image_data_response `json:"data,omitempty"`
}

func createBill(s_createBill_request s_createBill_request_hdr) (s_createBill_response s_createBill_response_hdr, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("PANIC - ", r)

			err = fmt.Errorf("panic")
		}
	}()
	s_createBill_response = s_createBill_response_hdr{}
	fmt.Println("CREATE BILL " + s_createBill_request.AuthToken)

	dbConn := db.Connect()

	defer dbConn.Close()

	s_login_credentials, err := db.GetAuthToken(dbConn, s_createBill_request.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {
		WScreateBillResponse, err := ws.CreateBill(s_login_credentials, s_createBill_request.Amount)
		if err != nil {
			s_createBill_response.StatusCode = 500
			s_createBill_response.ErrorMessage = "Internal server error"
		} else {
			s_createBill_response.Data = &s_createBill_response_data_hdr{}
			s_createBill_response.Data.Amount = WScreateBillResponse.XMLAgents.CreateBill.BoletoBill.Amount
			s_createBill_response.Data.Id = WScreateBillResponse.XMLAgents.CreateBill.BoletoBill.Id
		}
	} else {
		s_createBill_response.StatusCode = 403
		s_createBill_response.ErrorMessage = "Login/Senha inválido"
	}
	return s_createBill_response, nil
}
func getBillImage(s_getBill_request s_getBill_request_hdr) (s_geBill_image_response s_geBill_image_response_hdr, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("PANIC - ", r)

			err = fmt.Errorf("panic")
		}
	}()
	s_geBill_image_response = s_geBill_image_response_hdr{}
	fmt.Println("CREATE BILL " + s_getBill_request.AuthToken)

	dbConn := db.Connect()

	defer dbConn.Close()

	s_login_credentials, err := db.GetAuthToken(dbConn, s_getBill_request.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {
		wsResponse, err := ws.GetBillImage(s_login_credentials, s_getBill_request.BoletoId)
		if err != nil {
			s_geBill_image_response.StatusCode = 500
			s_geBill_image_response.ErrorMessage = "Internal server error"
		} else {
			s_geBill_image_response.Data = &s_geBill_image_data_response{}
			s_geBill_image_response.Data.Image = wsResponse.XMLAgents.GetBillImage.Image
		}
	} else {
		s_geBill_image_response.StatusCode = 403
		s_geBill_image_response.ErrorMessage = "Login/Senha inválido"
	}
	return s_geBill_image_response, nil
}

func getBalance(s_balance_request s_balance_request_hdr) (s_balance_response s_balance_response_hdr, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("PANIC - ", r)

			err = fmt.Errorf("panic")
		}
	}()

	s_balance_response = s_balance_response_hdr{}
	fmt.Println("GET BALANCE " + s_balance_request.AuthToken)

	dbConn := db.Connect()
	defer dbConn.Close()

	s_login_credentials, err := db.GetAuthToken(dbConn, s_balance_request.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {
		balance, overdraft, err := ws.GetBalance(s_login_credentials)
		if err != nil {
			s_balance_response.StatusCode = 500
			s_balance_response.ErrorMessage = "Internal server error"
		} else {

			s_balance_response.Data = &s_balance_response_data_hdr{}
			s_balance_response.Data.Balance = *balance
			s_balance_response.Data.Overdraft = *overdraft
		}
	} else {
		s_balance_response.StatusCode = 403
		s_balance_response.ErrorMessage = "Login/Senha inválido"
	}
	return s_balance_response, nil
}
