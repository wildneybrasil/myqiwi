package main

import (
	"db"
	"fmt"
	"strconv"
	"ws"
)

type s_cadastraPlaca_request_hdr struct {
	AuthToken string `json:"authToken"`
	Placa     string `json:"placa"`
	Nome      string `json:"nome"`
	Type      string `json:"type"`
}
type s_listaPlaca_request_hdr struct {
	AuthToken string `json:"authToken"`
}
type s_renamePlaca_request_hdr struct {
	AuthToken string `json:"authToken"`
	Placa     string `json:"placa"`
	Nome      string `json:"nome"`
}
type s_cet_extrato_request_hdr struct {
	AuthToken string `json:"authToken"`
	Days      string `json:"days"`
	Page      int    `json:"page"`
}
type s_cet_local_request_hdr struct {
	AuthToken string `json:"authToken"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type s_cet_buy_cad_request_hdr struct {
	AuthToken string `json:"authToken"`
	Qtde      int    `json:"qtde"`
}
type s_cet_ativate_cad_request_hdr struct {
	AuthToken    string `json:"authToken"`
	AddressId    int    `json:"addressId"`
	Ev_dateSend  string `json:"ev_dateSend"`
	Ev_tipoCET   int    `json:"ev_tipoCET"`
	Ev_document  string `json:"ev_document"`
	Ev_imei      string `json:"ev_imei"`
	Ev_latitude  string `json:"ev_latitude"`
	Ev_longitude string `json:"ev_longitude"`
	Ev_placa     string `json:"ev_placa"`
	Ev_temp      int    `json:"ev_temp"`
	Ev_qtdCAD    int    `json:"ev_qtdCAD"`
	Ev_check     int    `json:"ev_check"`
}

//

type s_cet_extrato_response_data_hdr struct {
	JSON string `json:"json"`
}
type s_cet_extrato_response_hdr struct {
	s_status
	Data *s_cet_extrato_response_data_hdr `json:"data,omitempty"`
}
type s_placa_list_response_data_hdr struct {
	Placas *[]db.Placa_hdr `json:"placas"`
}
type s_placa_list_response_hdr struct {
	s_status
	Data *s_placa_list_response_data_hdr `json:"data,omitempty"`
}
type s_placa_response_hdr struct {
	s_status
	JSON    string `json:"json,omitempty"`
	Date    string `json:"date,omitempty"`
	CETAuth string `json:"CETAuth,omitempty"`
}

/* PLACA */
func CadastraPlaca(s_cadastraPlaca_request s_cadastraPlaca_request_hdr) (s_placa_response s_placa_response_hdr, err error) {
	s_placa_response = s_placa_response_hdr{}

	fmt.Println("CADASTRA PLACA\n")
	dbConn := db.Connect()

	defer dbConn.Close()

	servico := "160930"

	s_login_credentials, err := db.GetAuthToken(dbConn, s_cadastraPlaca_request.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {
		if s_cadastraPlaca_request.Type == "" {
			s_placa_response.StatusCode = 500
			s_placa_response.Status = "failed"
			s_placa_response.ErrorMessage = "Missing type."

			return s_placa_response, nil
		}
		wsResponse, err := ws.CadastraPlaca(s_login_credentials, s_cadastraPlaca_request.Placa, servico, s_login_credentials.Document)

		if wsResponse.XMLProvider.XMLCheckPaymentRequisites.Result != "0" {
			s_placa_response.StatusCode = 400
			s_placa_response.ErrorMessage, s_placa_response.StatusCode = ws.GetErrorMessage(wsResponse.XMLProvider.XMLCheckPaymentRequisites.Result)
			return s_placa_response, nil
		}
		if wsResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp3 != "200" {
			s_placa_response.StatusCode = 400
			s_placa_response.ErrorMessage = wsResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp2
			return s_placa_response, nil
		}

		_, err = db.InsertPlaca(dbConn, s_cadastraPlaca_request.Nome, s_cadastraPlaca_request.Placa, s_cadastraPlaca_request.Type, s_login_credentials.Id)
		if err != nil {
			s_placa_response.StatusCode = 500
			s_placa_response.Status = "failed"
			s_placa_response.ErrorMessage = "Internal server error"
			return s_placa_response, nil
		}

	} else {
		s_placa_response.StatusCode = 403
		s_placa_response.ErrorMessage = "Login/Senha inválido"
		s_placa_response.Status = "failed"
	}
	return s_placa_response, nil
}

func RemovePlaca(s_cadastraPlaca_request s_cadastraPlaca_request_hdr) (s_placa_response s_placa_response_hdr, err error) {
	s_placa_response = s_placa_response_hdr{}

	dbConn := db.Connect()

	defer dbConn.Close()

	servico := "160930"

	s_login_credentials, err := db.GetAuthToken(dbConn, s_cadastraPlaca_request.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {
		wsResponse, err := ws.RemovePlaca(s_login_credentials, s_cadastraPlaca_request.Placa, s_login_credentials.Document, servico)

		if wsResponse.XMLProvider.XMLCheckPaymentRequisites.Result != "0" {
			s_placa_response.StatusCode = 400
			s_placa_response.ErrorMessage, s_placa_response.StatusCode = ws.GetErrorMessage(wsResponse.XMLProvider.XMLCheckPaymentRequisites.Result)
			return s_placa_response, nil
		}
		if wsResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp3 != "200" {
			s_placa_response.StatusCode = 400
			s_placa_response.ErrorMessage = wsResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp2
			return s_placa_response, nil
		}

		_, err = db.DeletePlaca(dbConn, s_cadastraPlaca_request.Placa, s_login_credentials.Id)
		if err != nil {
			s_placa_response.StatusCode = 500
			s_placa_response.Status = "failed"
			s_placa_response.ErrorMessage = "Internal server error"
			return s_placa_response, nil
		}

	} else {
		s_placa_response.StatusCode = 403
		s_placa_response.ErrorMessage = "Login/Senha inválido"
		s_placa_response.Status = "failed"
	}
	return s_placa_response, nil
}
func ListPlacas(s_listaPlaca_request s_listaPlaca_request_hdr) (s_placa_response s_placa_list_response_hdr, err error) {
	s_placa_response = s_placa_list_response_hdr{}

	fmt.Println("LIST PLACA\n")
	dbConn := db.Connect()

	defer dbConn.Close()

	s_login_credentials, err := db.GetAuthToken(dbConn, s_listaPlaca_request.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {
		placas, err := db.ListPlaca(dbConn, s_login_credentials.Id)

		if err != nil {
			s_placa_response.StatusCode = 500
			s_placa_response.Status = "failed"
			s_placa_response.ErrorMessage = "DB ERROR"

			return s_placa_response, nil
		}
		fmt.Printf("FOUND: %d\n", len(*placas))

		s_placa_response.Data = &s_placa_list_response_data_hdr{}
		s_placa_response.Data.Placas = placas
	} else {
		s_placa_response.StatusCode = 403
		s_placa_response.ErrorMessage = "Login/Senha inválido"
		s_placa_response.Status = "failed"
	}
	return s_placa_response, nil
}
func RenamePlaca(s_renamePlaca_request s_renamePlaca_request_hdr) (s_placa_response s_placa_list_response_hdr, err error) {
	s_placa_response = s_placa_list_response_hdr{}

	fmt.Println("RENAME PLACA\n")
	dbConn := db.Connect()

	defer dbConn.Close()

	s_login_credentials, err := db.GetAuthToken(dbConn, s_renamePlaca_request.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {
		err := db.RenamePlaca(dbConn, s_renamePlaca_request.Nome, s_renamePlaca_request.Placa, s_login_credentials.Id)

		if err != nil {
			s_placa_response.StatusCode = 500
			s_placa_response.Status = "failed"
			s_placa_response.ErrorMessage = "DB ERROR"

			return s_placa_response, nil
		}

	} else {
		s_placa_response.StatusCode = 403
		s_placa_response.ErrorMessage = "Login/Senha inválido"
		s_placa_response.Status = "failed"
	}
	return s_placa_response, nil
}
func CETExtrato(s_cet_extrato_request s_cet_extrato_request_hdr) (s_placa_response s_placa_response_hdr, err error) {
	s_placa_response = s_placa_response_hdr{}

	dbConn := db.Connect()

	defer dbConn.Close()

	servico := "160930"

	s_login_credentials, err := db.GetAuthToken(dbConn, s_cet_extrato_request.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {
		if s_cet_extrato_request.Days == "" {
			s_placa_response.StatusCode = 500
			s_placa_response.Status = "failed"
			s_placa_response.ErrorMessage = "Missing days."

			return s_placa_response, nil
		}

		wsResponse, err := ws.ListaExtrato(s_login_credentials, s_cet_extrato_request.Page, s_cet_extrato_request.Days, s_login_credentials.Document, servico)

		if wsResponse.XMLProvider.XMLCheckPaymentRequisites.Result != "0" {
			s_placa_response.StatusCode = 400
			s_placa_response.ErrorMessage, s_placa_response.StatusCode = ws.GetErrorMessage(wsResponse.XMLProvider.XMLCheckPaymentRequisites.Result)
			return s_placa_response, nil
		}
		if wsResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Result != "0" {
			s_placa_response.StatusCode = 400
			s_placa_response.ErrorMessage, s_placa_response.StatusCode = ws.GetErrorMessage(wsResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Result)
			return s_placa_response, nil
		}
		if wsResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp3 != "200" {
			s_placa_response.StatusCode = 400
			s_placa_response.ErrorMessage = wsResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp2
			return s_placa_response, nil
		}
		s_placa_response.JSON = wsResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp2

		if err != nil {
			s_placa_response.StatusCode = 500
			s_placa_response.Status = "failed"
			s_placa_response.ErrorMessage = "Internal server error"
			return s_placa_response, nil
		}

	} else {
		s_placa_response.StatusCode = 403
		s_placa_response.ErrorMessage = "Login/Senha inválido"
		s_placa_response.Status = "failed"
	}
	return s_placa_response, nil
}
func CETCompraCad(s_cet_buy_cad_request s_cet_buy_cad_request_hdr) (s_placa_response s_placa_response_hdr, err error) {
	s_placa_response = s_placa_response_hdr{}

	dbConn := db.Connect()

	defer dbConn.Close()

	servico := "160930"

	s_login_credentials, err := db.GetAuthToken(dbConn, s_cet_buy_cad_request.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {
		amount := 5.0 * s_cet_buy_cad_request.Qtde
		desconto := s_cet_buy_cad_request.Qtde / 10.0
		amount = amount - (desconto * 5.0)

		strQtde := strconv.Itoa(s_cet_buy_cad_request.Qtde)
		strAmount := strconv.FormatFloat(float64(amount), 'f', 2, 64)

		fmt.Println("AMOUNT: ", strAmount)
		date := ""
		wsResponse, err := ws.CompraCAD(s_login_credentials, strQtde, servico, s_login_credentials.Cel, date, s_login_credentials.Document, strAmount)
		if err != nil {
			s_placa_response.StatusCode = 500
			s_placa_response.Status = "failed"
			s_placa_response.ErrorMessage = err.Error()
			return s_placa_response, nil
		}

		if wsResponse.XMLProvider.XMLPurchaseOnline.Result != "0" {
			s_placa_response.StatusCode = 400
			s_placa_response.ErrorMessage, s_placa_response.StatusCode = ws.GetErrorMessage(wsResponse.XMLProvider.XMLCheckPaymentRequisites.Result)
			return s_placa_response, nil
		}
		if wsResponse.XMLProvider.XMLPurchaseOnline.XMLPayment.Result != "0" {
			s_placa_response.StatusCode = 400
			s_placa_response.ErrorMessage, s_placa_response.StatusCode = ws.GetErrorMessage(wsResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Result)
			return s_placa_response, nil
		}

		if wsResponse.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentExtras.Disp3 != "200" {
			s_placa_response.StatusCode = 400
			s_placa_response.ErrorMessage = wsResponse.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentExtras.Disp2
			return s_placa_response, nil
		}
		s_placa_response.StatusCode = 0

	} else {
		s_placa_response.StatusCode = 403
		s_placa_response.ErrorMessage = "Login/Senha inválido"
		s_placa_response.Status = "failed"
	}
	return s_placa_response, nil
}
func CETAtivaCad(s_cet_ativate_cad_request s_cet_ativate_cad_request_hdr) (s_placa_response s_placa_response_hdr, err error) {
	s_placa_response = s_placa_response_hdr{}

	fmt.Println("Ativa CAD 3\n")

	dbConn := db.Connect()

	defer dbConn.Close()

	servico := "160930"

	s_login_credentials, err := db.GetAuthToken(dbConn, s_cet_ativate_cad_request.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {
		wsResponse, err := ws.AtivaCAD(s_login_credentials, s_cet_ativate_cad_request.Ev_placa, s_login_credentials.Document, servico, s_cet_ativate_cad_request.Ev_dateSend, s_cet_ativate_cad_request.Ev_imei, s_cet_ativate_cad_request.Ev_latitude, s_cet_ativate_cad_request.Ev_longitude, s_cet_ativate_cad_request.Ev_placa, s_cet_ativate_cad_request.Ev_temp, s_cet_ativate_cad_request.Ev_qtdCAD, s_cet_ativate_cad_request.Ev_check)
		if err != nil {
			s_placa_response.StatusCode = 500
			s_placa_response.Status = "failed"
			s_placa_response.ErrorMessage = err.Error()
			return s_placa_response, nil
		}

		if wsResponse.XMLProvider.XMLCheckPaymentRequisites.Result != "0" {
			s_placa_response.StatusCode = 400
			s_placa_response.ErrorMessage, s_placa_response.StatusCode = ws.GetErrorMessage(wsResponse.XMLProvider.XMLCheckPaymentRequisites.Result)
			return s_placa_response, nil
		}

		if wsResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp3 != "200" {
			number, _ := strconv.ParseInt(wsResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp3, 10, 0)
			s_placa_response.StatusCode = int(number)
			s_placa_response.ErrorMessage = wsResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp2
			return s_placa_response, nil
		}
		s_placa_response.StatusCode = 0
		s_placa_response.JSON = wsResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp2
		s_placa_response.CETAuth = wsResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.PrtData6
		s_placa_response.Date = wsResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.PrtData5

	} else {
		s_placa_response.StatusCode = 403
		s_placa_response.ErrorMessage = "Login/Senha inválido"
		s_placa_response.Status = "failed"
	}
	return s_placa_response, nil
}
func CETLocal(s_cet_local_request s_cet_local_request_hdr) (s_placa_response s_placa_response_hdr, err error) {
	s_placa_response = s_placa_response_hdr{}

	dbConn := db.Connect()

	defer dbConn.Close()

	servico := "160930"

	s_login_credentials, err := db.GetAuthToken(dbConn, s_cet_local_request.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {
		if s_cet_local_request.Latitude == "" {
			s_placa_response.StatusCode = 500
			s_placa_response.Status = "failed"
			s_placa_response.ErrorMessage = "Missing latitude."

			return s_placa_response, nil
		}
		if s_cet_local_request.Longitude == "" {
			s_placa_response.StatusCode = 500
			s_placa_response.Status = "failed"
			s_placa_response.ErrorMessage = "Missing lonitude."

			return s_placa_response, nil
		}
		wsResponse, err := ws.ListaLocal(s_login_credentials, s_cet_local_request.Latitude, s_cet_local_request.Longitude, servico)

		if wsResponse.XMLProvider.XMLCheckPaymentRequisites.Result != "0" {
			s_placa_response.StatusCode = 400
			s_placa_response.ErrorMessage, s_placa_response.StatusCode = ws.GetErrorMessage(wsResponse.XMLProvider.XMLCheckPaymentRequisites.Result)
			return s_placa_response, nil
		}
		if wsResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp3 != "200" {
			s_placa_response.StatusCode = 400
			s_placa_response.ErrorMessage = wsResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp2
			return s_placa_response, nil
		}

		if err != nil {
			s_placa_response.StatusCode = 500
			s_placa_response.Status = "failed"
			s_placa_response.ErrorMessage = "Internal server error"
			return s_placa_response, nil
		}

		s_placa_response.JSON = wsResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp2

	} else {
		s_placa_response.StatusCode = 403
		s_placa_response.ErrorMessage = "Login/Senha inválido"
		s_placa_response.Status = "failed"
	}
	return s_placa_response, nil
}
