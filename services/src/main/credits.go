// main
package main

import (
	"db"
	"fmt"
	"strconv"
	"ws"
)

type s_createBill_request_hdr struct {
	AuthToken string `json:"authToken"`
	Amount    string `json:"amount"`
}
type s_getBill_request_hdr struct {
	AuthToken string `json:"authToken"`
	BoletoId  string `json:"boletoId"`
	Pdf       bool   `json:"pdf"`
	Info      bool   `json:"info"`
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
type s_geBill_image_data_response_items struct {
	Id    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}
type s_geBill_image_data_response struct {
	Image           string                               `json:"image,omitempty"`
	IssuerAgentId   string                               `json:"issuerAgentId,omitempty"`
	IssuerAgentName string                               `json:"issuerAgentName,omitempty"`
	Amount          string                               `json:"amount,omitempty"`
	Comission       string                               `json:"comission,omitempty"`
	BankName        string                               `json:"bankName,omitempty"`
	ExpireDate      string                               `json:"expireDate,omitempty"`
	CreateTime      string                               `json:"createTime,omitempty"`
	Instructions    string                               `json:"instructions,omitempty"`
	ReceiverAddress string                               `json:"receiverAddress,omitempty"`
	ReceiverInn     string                               `json:"ReceiverInn,omitempty"`
	IssuerAgentInn  string                               `json:"IssuerAgentInn,omitempty"`
	Ipte            string                               `json:"Ipte,omitempty"`
	TypeLine        string                               `json:"typeLine,omitempty"`
	OwnNumber       string                               `json:"ownNumber,omitempty"`
	CustomFields    []s_geBill_image_data_response_items `json:"customFields,omitempty"`
}
type s_geBill_image_response_hdr struct {
	s_status
	Data *s_geBill_image_data_response `json:"data,omitempty"`
}

func createBill(s_createBill_request s_createBill_request_hdr) (s_createBill_response s_createBill_response_hdr, err error) {
	defer func() {
		if r := recover(); r != nil {
			s_createBill_response.StatusCode = 500
			s_createBill_response.ErrorMessage = "Internal server error"

			fmt.Println("PANIC - ", r)

			err = fmt.Errorf("panic")
		}
	}()
	s_createBill_response = s_createBill_response_hdr{}
	fmt.Println("CREATE BILL " + s_createBill_request.AuthToken)

	dbConn := db.Connect()

	defer dbConn.Close()

	floatvalue, _ := strconv.ParseFloat(s_createBill_request.Amount, 64)
	if floatvalue == 0 || floatvalue < 10 {
		s_createBill_response.StatusCode = 400
		s_createBill_response.Status = "failed"
		s_createBill_response.ErrorMessage = "O valor mínimo para gerar créditos é de R$10,00"
		return s_createBill_response, nil
	}
	if floatvalue > 10000 {
		s_createBill_response.StatusCode = 400
		s_createBill_response.Status = "failed"
		s_createBill_response.ErrorMessage = "O valor máximo para gerar créditos é de R$10000,00"
		return s_createBill_response, nil
	}

	s_login_credentials, err := db.GetAuthToken(dbConn, s_createBill_request.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {

		WScreateBillResponse, err := ws.CreateBill(s_login_credentials, s_createBill_request.Amount)
		if err != nil {
			s_createBill_response.StatusCode = 500
			s_createBill_response.Status = "failed"
			s_createBill_response.ErrorMessage = "Internal server error"
		} else {
			if WScreateBillResponse.Result != "0" {
				s_createBill_response.ErrorMessage, s_createBill_response.StatusCode = ws.GetErrorMessage(WScreateBillResponse.Result)
				return s_createBill_response, nil
			}
			if WScreateBillResponse.XMLAgents.CreateBill.Result != "0" {
				s_createBill_response.ErrorMessage, s_createBill_response.StatusCode = ws.GetErrorMessage(WScreateBillResponse.XMLAgents.CreateBill.Result)
				return s_createBill_response, nil
			}
			s_createBill_response.Data = &s_createBill_response_data_hdr{}
			s_createBill_response.Data.Amount = WScreateBillResponse.XMLAgents.CreateBill.BoletoBill.Amount
			s_createBill_response.Data.Id = WScreateBillResponse.XMLAgents.CreateBill.BoletoBill.Id
			s_createBill_response.Status = "success"
		}
	} else {
		s_createBill_response.Status = "failed"
		s_createBill_response.StatusCode = 403
		s_createBill_response.ErrorMessage = "Login/Senha inválido"
	}
	return s_createBill_response, nil
}
func getBillInfo(s_getBill_request s_getBill_request_hdr) (s_geBill_image_response s_geBill_image_response_hdr, err error) {
	//	defer func() {
	//		if r := recover(); r != nil {
	//			fmt.Println("PANIC - ", r)

	//			err = fmt.Errorf("panic")
	//		}
	//	}()
	s_geBill_image_response = s_geBill_image_response_hdr{}
	fmt.Println("CREATE BILL " + s_getBill_request.AuthToken)

	dbConn := db.Connect()

	defer dbConn.Close()

	s_login_credentials, err := db.GetAuthToken(dbConn, s_getBill_request.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {
		var boletoInfo *ws.WSResponse_createBill_hdr
		var boletoImage *ws.WSResponse_createBill_hdr

		if s_getBill_request.Pdf {
			boletoImage, err = ws.GetBillImage(s_login_credentials, s_getBill_request.BoletoId)
			if err != nil {
				s_geBill_image_response.Status = "failed"
				s_geBill_image_response.StatusCode = 500
				s_geBill_image_response.ErrorMessage = "Internal server error"
				return s_geBill_image_response, nil
			}
			s_geBill_image_response.Data = &s_geBill_image_data_response{}
			fmt.Println("PDF: " + boletoImage.XMLAgents.GetBillImage.Image)
			s_geBill_image_response.Data.Image = boletoImage.XMLAgents.GetBillImage.Image
		}
		if s_getBill_request.Info {
			boletoInfo, err = ws.GetBillInfo(s_login_credentials, s_getBill_request.BoletoId)
			if err != nil {
				s_geBill_image_response.Status = "failed"
				s_geBill_image_response.StatusCode = 500
				s_geBill_image_response.ErrorMessage = "Internal server error"
				return s_geBill_image_response, nil
			}
			if s_geBill_image_response.Data == nil {
				s_geBill_image_response.Data = &s_geBill_image_data_response{}

			}
			s_geBill_image_response.Data.BankName = boletoInfo.XMLAgents.GetBillImage.Bill.BankName
			s_geBill_image_response.Data.Comission = boletoInfo.XMLAgents.GetBillImage.Bill.Comission
			s_geBill_image_response.Data.CreateTime = boletoInfo.XMLAgents.GetBillImage.Bill.CreateTime
			s_geBill_image_response.Data.ExpireDate = boletoInfo.XMLAgents.GetBillImage.Bill.ExpireDate
			s_geBill_image_response.Data.Instructions = boletoInfo.XMLAgents.GetBillImage.Bill.Instructions
			s_geBill_image_response.Data.Ipte = boletoInfo.XMLAgents.GetBillImage.Bill.Ipte
			s_geBill_image_response.Data.IssuerAgentId = boletoInfo.XMLAgents.GetBillImage.Bill.IssuerAgentId
			s_geBill_image_response.Data.IssuerAgentInn = boletoInfo.XMLAgents.GetBillImage.Bill.IssuerAgentInn
			s_geBill_image_response.Data.IssuerAgentName = boletoInfo.XMLAgents.GetBillImage.Bill.IssuerAgentName
			s_geBill_image_response.Data.OwnNumber = boletoInfo.XMLAgents.GetBillImage.Bill.OwnNumber
			s_geBill_image_response.Data.ReceiverAddress = boletoInfo.XMLAgents.GetBillImage.Bill.ReceiverAddress
			s_geBill_image_response.Data.ReceiverInn = boletoInfo.XMLAgents.GetBillImage.Bill.ReceiverInn
			s_geBill_image_response.Data.TypeLine = boletoInfo.XMLAgents.GetBillImage.Bill.TypeLine
			customValues := make([]s_geBill_image_data_response_items, 0)
			for _, v := range boletoInfo.XMLAgents.GetBillImage.Bill.CustomFields.Field {
				value := s_geBill_image_data_response_items{}
				value.Id = v.Id
				value.Name = v.Name
				value.Value = v.Value

				customValues = append(customValues, value)
			}
			s_geBill_image_response.Data.CustomFields = customValues
		}
		s_geBill_image_response.Status = "success"

		return s_geBill_image_response, nil
	} else {
		s_geBill_image_response.Status = "failed"
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
			s_balance_response.Status = "failed"
			s_balance_response.StatusCode = 500
			s_balance_response.ErrorMessage = "Internal server error"
		} else {

			s_balance_response.Data = &s_balance_response_data_hdr{}
			s_balance_response.Data.Balance = *balance
			s_balance_response.Data.Overdraft = *overdraft
		}
	} else {
		s_balance_response.Status = "failed"
		s_balance_response.StatusCode = 403
		s_balance_response.ErrorMessage = "Login/Senha inválido"
	}
	return s_balance_response, nil
}
