// main
package main

import (
	"db"
	"fmt"
	"notification"
	"strconv"
	"strings"
	"ws"

	"html"
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
	Amount    string `json:"amount"`
	Service   string `json:"service"`
}

type s_payment_request_hdr struct {
	AuthToken      string `json:"authToken"`
	Rcpt           string `json:"rcpt"`
	Service        string `json:"service"`
	Id             string `json:"id"`
	Session        string `json:"session"`
	Boleto         string `json:"boleto"`
	Amount         string `json:"amount"`
	SelectedAmount string `json:"selectedAmount"`
	Type           string `json:"type"`
	CardNumber     string `json:"cardNumber"`
	Password       string `json:"password"`
	EvCardType     string `json:"evCardType"`
	EvCardUID      string `json:"evCardUID"`
	EvCardData     string `json:"evCardData"`
	EvNominal      string `json:"evNominal"`
	Ev_wallet_code string `json:"ev_wallet_code"`
}

type s_payment_response_data_hdr struct {
	PIN         string                 `json:"pin,omitempty"`
	Serial      string                 `json:"serial,omitempty"`
	Session     string                 `json:"session,omitempty"`
	Info        string                 `json:"info,omitempty"`
	Id          string                 `json:"id,omitempty"`
	VoucherCode string                 `json:"id,omitempty"`
	PrtData1    string                 `json:"prtData1,omitempty"`
	RptData1    string                 `json:"rptData1,omitempty"`
	Nominals    *s_values_response_hdr `json:"nominals,omitempty"`
}
type s_payment_response_hdr struct {
	s_status
	Data *s_payment_response_data_hdr `json:"data,omitempty"`
}

func getBoletoInfo(s_boletoInfo_request s_boletoInfo_request_hdr) (s_boletoResponse s_boletoResponse_hdr, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("PANIC - ", r)

			err = fmt.Errorf("panic")
		}
	}()
	if s_boletoInfo_request.Service == "" || s_boletoInfo_request.Boleto == "" || s_boletoInfo_request.Scanned == "" {
		s_boletoResponse.StatusCode = 500
		s_boletoResponse.Status = "failed"
		s_boletoResponse.ErrorMessage = "Missing type."

		return s_boletoResponse, nil
	}

	dbConn := db.Connect()

	defer dbConn.Close()

	s_login_credentials, err := db.GetAuthToken(dbConn, s_boletoInfo_request.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {
		transferResponse, err := ws.GetBoletoInfo(s_login_credentials, s_boletoInfo_request.Boleto, s_boletoInfo_request.Amount, s_login_credentials.Cel, s_boletoInfo_request.Scanned, s_boletoInfo_request.Service)
		if err != nil {
			s_boletoResponse.StatusCode = 500
			s_boletoResponse.Status = "failed"
			s_boletoResponse.ErrorMessage = "Internal server error"
			return s_boletoResponse, nil
		}
		if transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Result != "0" || transferResponse.XMLProvider.XMLCheckPaymentRequisites.Result != "0" {
			s_boletoResponse.StatusCode = 400
			s_boletoResponse.ErrorMessage = "Internal server error"
			s_boletoResponse.Status = "failed"
			return s_boletoResponse, nil
		}
		if transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp3 != "true" {
			s_boletoResponse.StatusCode = 400
			s_boletoResponse.ErrorMessage = html.UnescapeString(transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp3)
			s_boletoResponse.Status = "failed"
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
		s_boletoResponse.Status = "failed"
	}
	return s_boletoResponse, nil
}

func payment1(s_payment_request s_payment_request_hdr) (s_payment_response s_payment_response_hdr, err error) {

	s_payment_response = s_payment_response_hdr{}

	dbConn := db.Connect()

	defer dbConn.Close()

	s_login_credentials, err := db.GetAuthToken(dbConn, s_payment_request.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {

		if s_payment_request.Type == "" {
			s_payment_response.StatusCode = 500
			s_payment_response.Status = "failed"
			s_payment_response.ErrorMessage = "Missing type."

			return s_payment_response, nil
		}
		serviceId, _ := strconv.ParseInt(s_payment_request.Service, 10, 0)
		serviceInfo, _ := db.GetServiceByPrid(dbConn, int(serviceId))

		reqType := "2"
		isNom := "1"
		forceAmount := "true"
		evStep := ""

		switch serviceInfo.PaymentType {
		case "Corban":
			break
		case "Credisan":
			break
		case "Software Express":
			reqType = "4"
			isNom = ""
			forceAmount = "true"
			break
		case "Pin Offline":
			break
		case "RV":
			reqType = "2"
			break
		case "QIWI":
			break
		case "Transporte":
			evStep = "1"
			forceAmount = ""
			isNom = ""
			reqType = ""
			break

		default:
			break
		}

		// telefonia
		if s_payment_request.Type == "telefonia" {
			transferResponse, err := ws.DoPaymentTel1(s_login_credentials, strings.Replace(s_payment_request.Rcpt, " ", "", -1), s_payment_request.Service, forceAmount, evStep, isNom, reqType)

			if err != nil {
				s_payment_response.StatusCode = 500
				s_payment_response.Status = "failed"
				s_payment_response.ErrorMessage = "Internal server error"
				return s_payment_response, nil
			}
			if transferResponse.XMLProvider.XMLCheckPaymentRequisites.Result != "0" {
				s_payment_response.StatusCode = 400
				s_payment_response.ErrorMessage, s_payment_response.StatusCode = ws.GetErrorMessage(transferResponse.XMLProvider.XMLCheckPaymentRequisites.Result)
				return s_payment_response, nil
			}
			if transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Result != "0" {
				s_payment_response.StatusCode = 400
				s_payment_response.ErrorMessage, s_payment_response.StatusCode = ws.GetErrorMessage(transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Result)
				return s_payment_response, nil
			}
			strNominalsValue1 := transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp1
			strNominalsValue3 := transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp3
			V1 := strings.Split(strNominalsValue1, ";")
			V3 := strings.Split(strNominalsValue3, ";")

			s_payment_response.Data = &s_payment_response_data_hdr{}
			s_payment_response.Data.Nominals = &s_values_response_hdr{}

			for k, _ := range V1 {

				item := s_values_item_hdr{}

				if len(V1[k]) > 0 {
					amount := strings.Replace(V1[k], "|", "", -1)

					amountString := ""
					if strings.Contains(amount, "-") {
						sV1 := strings.Split(amount, "-")

						amountFloat1, _ := strconv.ParseFloat(sV1[0], 64)
						amountFloat2, _ := strconv.ParseFloat(sV1[1], 64)
						amountString = fmt.Sprintf("%.02f-%.02f", amountFloat1/100, amountFloat2/100)
					} else {
						amountFloat, _ := strconv.ParseFloat(amount, 64)
						amountString = fmt.Sprintf("%.02f", amountFloat/100)
					}

					item.Amount = amountString
					item.Id = strings.Replace(V3[k], "|", "", -1)

					s_payment_response.Data.Nominals.Items = append(s_payment_response.Data.Nominals.Items, item)
				}
			}
			s_payment_response.Data.Id = transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Id
			switch serviceInfo.PaymentType {
			case "Corban":
				s_payment_response.Data.Session = transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp4
				break
			case "Credisan":
				s_payment_response.Data.Session = transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp4
				break
			case "Software Express":
				s_payment_response.Data.Session = ""
				strNominalsValue1 := transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Ev_combo1
				V1 := strings.Split(strNominalsValue1, "|")

				//				s_payment_response.Data = &s_payment_response_data_hdr{}
				//				s_payment_response.Data.Nominals = &s_values_response_hdr{}

				s_payment_response.Data.Nominals.Items = []s_values_item_hdr{}

				for k, _ := range V1 {
					V2 := strings.Split(V1[k], ";")

					item := s_values_item_hdr{}

					if len(V2) >= 2 {
						amount := V2[2]
						amountFloat, _ := strconv.ParseFloat(amount, 64)
						amountString := fmt.Sprintf("%.02f", amountFloat)

						item.Amount = amountString
						item.Id = V2[1]
						item.Text = V2[0]

						if amountFloat > 0.1 {
							s_payment_response.Data.Nominals.Items = append(s_payment_response.Data.Nominals.Items, item)
						}
					}
				}

				break
			case "Pin Offline":
				s_payment_response.Data.Session = transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp4
				break
			case "RV":
				s_payment_response.Data.Session = transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp4
				break
			case "QIWI":
				s_payment_response.Data.Session = transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp4
				break
			default:
				break
			}
		}
		// transporte

		if s_payment_request.Type == "transporte" {
			cardNumber := strings.Replace(s_payment_request.CardNumber, ".", "", -1)
			cardNumber = strings.Replace(cardNumber, "-", "", -1)

			transferResponse, err := ws.DoPaymentTrans1(s_login_credentials, cardNumber, s_payment_request.Service, forceAmount, evStep, isNom, reqType)

			if err != nil {
				s_payment_response.Status = "failed"
				s_payment_response.StatusCode = 500
				s_payment_response.ErrorMessage = "Internal server error"
				return s_payment_response, nil
			}
			if transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Result != "0" {
				s_payment_response.ErrorMessage, s_payment_response.StatusCode = ws.GetErrorMessage(transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Result)
				return s_payment_response, nil

			}
			if transferResponse.XMLProvider.XMLCheckPaymentRequisites.Result != "0" {
				s_payment_response.ErrorMessage, s_payment_response.StatusCode = ws.GetErrorMessage(transferResponse.XMLProvider.XMLCheckPaymentRequisites.Result)
				return s_payment_response, nil
			}
			if transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp2 != "true" {
				s_payment_response.StatusCode = 400
				s_payment_response.ErrorMessage = transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp2
				return s_payment_response, nil
			}

			s_payment_response.Data = &s_payment_response_data_hdr{}
			s_payment_response.Data.Nominals = &s_values_response_hdr{}
			// fake nominal
			item := s_values_item_hdr{}
			item.Amount = "20"
			item.Id = "1"

			s_payment_response.Data.Nominals.Items = []s_values_item_hdr{item}

			s_payment_response.Data.Session = transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp1
			s_payment_response.Data.Id = transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Id
		}

		// games
		if s_payment_request.Type == "games" {
			transferResponse, err := ws.DoPaymentGames1(s_login_credentials, s_payment_request.Service)

			if err != nil {
				s_payment_response.StatusCode = 500
				s_payment_response.Status = "failed"
				s_payment_response.ErrorMessage = "Internal server error"
				return s_payment_response, nil
			}
			s_payment_response.Data = &s_payment_response_data_hdr{}
			s_payment_response.Data.Nominals = &s_values_response_hdr{}

			if transferResponse.XMLProvider.XMLGetNomenclature.Goods == nil {
				s_payment_response.StatusCode = 500
				s_payment_response.Status = "failed"
				s_payment_response.ErrorMessage = "Internal server error"
				return s_payment_response, nil
			}
			for _, v := range *transferResponse.XMLProvider.XMLGetNomenclature.Goods {
				item := s_values_item_hdr{}
				item.Amount = v.Amount
				item.Id = v.Id

				s_payment_response.Data.Nominals.Items = append(s_payment_response.Data.Nominals.Items, item)
			}
		}

	} else {
		s_payment_response.StatusCode = 403
		s_payment_response.ErrorMessage = "Login/Senha inválido"
		s_payment_response.Status = "failed"
	}
	return s_payment_response, nil
}
func payment2(s_payment_request s_payment_request_hdr) (s_payment_response s_payment_response_hdr, err error) {
	if s_payment_request.Type == "" || s_payment_request.Service == "" {
		s_payment_response.StatusCode = 500
		s_payment_response.ErrorMessage = "Missing type."
		s_payment_response.Status = "failed"

		return s_payment_response, nil
	}

	dbConn := db.Connect()

	serviceId, _ := strconv.ParseInt(s_payment_request.Service, 10, 0)
	serviceInfo, err := db.GetServiceByPrid(dbConn, int(serviceId))
	if err != nil {
		s_payment_response.StatusCode = 500
		s_payment_response.Status = "failed"
		s_payment_response.ErrorMessage = "Invalid service"

		return s_payment_response, nil
	}

	defer dbConn.Close()

	s_login_credentials, err := db.GetAuthToken(dbConn, s_payment_request.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {
		if !CheckPassword(s_payment_request.Password, s_login_credentials.TerminalPassword, s_login_credentials.PasswordSalt) {
			s_payment_response.StatusCode = 403
			s_payment_response.Status = "failed"
			s_payment_response.ErrorMessage = "Login/Senha inválido."

			return s_payment_response, nil
		}

		session := ""
		fmt.Printf("PAYMENT TYPE [%s]\n", serviceInfo.PaymentType)

		switch serviceInfo.PaymentType {
		case "Corban":
			session = s_payment_request.Session
			s_payment_request.Rcpt = s_login_credentials.Cel
			break
		case "Credisan":
			session = s_payment_request.Session
			s_payment_request.Rcpt = s_login_credentials.Cel
			break
		case "Transporte":
			session = s_payment_request.Session
			break
		// case "Software Express":
		// case "Pin Offline":
		// case "RV":
		// case "QIWI":

		default:
			fmt.Println("TEL 2")
			transferResponse1, err := ws.DoPaymentTel2(s_login_credentials, s_payment_request.Id, strings.Replace(s_payment_request.Rcpt, " ", "", -1), s_payment_request.Service, s_payment_request.Amount)
			if err != nil {
				s_payment_response.StatusCode = 500
				s_payment_response.Status = "failed"
				s_payment_response.ErrorMessage = "Internal server error"

				return s_payment_response, nil
			} else {
				session = transferResponse1.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp1
			}
			break
		}
		transferResponse2 := &ws.WSResponse_transferCredits_hdr{}
		var requestXML *string
		var responseXML *string

		switch serviceInfo.PaymentType {
		case "Pin Offline":
			transferResponse2, requestXML, responseXML, err = ws.DoPaymentGames2(s_login_credentials, s_payment_request.Rcpt, s_payment_request.Service, s_payment_request.Amount, s_payment_request.SelectedAmount)
			break
		case "Transporte":
			fmt.Println("TRANSRPOTE 2")
			s_payment_request.Rcpt = "99999999999"
			transferResponse2, requestXML, responseXML, err = ws.DoPaymentTel3(s_login_credentials, s_payment_request.Id, session, s_payment_request.Rcpt, s_payment_request.Service, s_payment_request.Amount)
			break
		// case "Corban":
		// case "Credisan":
		// case "Software Express":
		// case "RV":
		// case "QIWI":
		default:
			fmt.Println("TEL 3")
			transferResponse2, requestXML, responseXML, err = ws.DoPaymentTel3(s_login_credentials, s_payment_request.Id, session, s_payment_request.Rcpt, s_payment_request.Service, s_payment_request.Amount)

			break
		}

		if err != nil {
			s_payment_response.StatusCode = 500
			s_payment_response.Status = "failed"
			s_payment_response.ErrorMessage = "Internal server error"
			return s_payment_response, nil
		} else {
			if transferResponse2.XMLProvider.XMLPurchaseOnline.Result != "0" {
				s_payment_response.ErrorMessage, s_payment_response.StatusCode = ws.GetErrorMessage(transferResponse2.XMLProvider.XMLPurchaseOnline.Result)
				return s_payment_response, nil
			}
			if transferResponse2.XMLProvider.XMLPurchaseOnline.XMLPayment.Result != "0" {
				s_payment_response.ErrorMessage, s_payment_response.StatusCode = ws.GetErrorMessage(transferResponse2.XMLProvider.XMLPurchaseOnline.XMLPayment.Result)
				return s_payment_response, nil
			}
			switch serviceInfo.PaymentType {
			case "Corban":
				break
			case "Credisan":
				break
			case "Software Express":
			case "Pin Offline":
				fmt.Println("PIN OFFLINE")
				cursor := transferResponse2.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLGoods.Item.Param
				fmt.Printf("C SIZE: %d\n", len(cursor))

				s_payment_response.Data = &s_payment_response_data_hdr{}
				for _, v := range cursor {
					fmt.Println("KEY: " + v.Name + " VALUE [" + v.Value + "]")
					if v.Name == "Pin" {
						s_payment_response.Data.PIN = v.Value
					}
					if v.Name == "%OSMP_UPG_PIN_SERIAL%" {
						s_payment_response.Data.Serial = v.Value

					}
				}
				cel := strings.Replace(s_payment_request.Rcpt, " ", "", -1)
				notification.Send(notification.NotificationMessage{"sms", cel, "PIN: " + s_payment_response.Data.PIN + " SERIAL: " + s_payment_response.Data.Serial})

				break
			case "RV":
			case "QIWI":
			default:
				break
			}
			db.InsertPaymentHistory(dbConn, s_login_credentials.Id, s_payment_request.Type, s_payment_request.Service, s_payment_request, s_payment_response, requestXML, responseXML, 1)

		}
	} else {
		s_payment_response.StatusCode = 403
		s_payment_response.Status = "failed"
		s_payment_response.ErrorMessage = "Login/Senha inválido"
	}
	return s_payment_response, nil
}

func paymentNFC1(s_payment_request s_payment_request_hdr) (s_payment_response s_payment_response_hdr, err error) {

	s_payment_response = s_payment_response_hdr{}

	dbConn := db.Connect()

	defer dbConn.Close()

	s_login_credentials, err := db.GetAuthToken(dbConn, s_payment_request.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {

		if s_payment_request.Type == "" {
			s_payment_response.StatusCode = 500
			s_payment_response.Status = "failed"
			s_payment_response.ErrorMessage = "Missing type."

			return s_payment_response, nil
		}

		transferResponse, err := ws.DoPaymentTransNFC1(s_login_credentials, s_payment_request.Service, s_payment_request.Amount, s_payment_request.EvCardType, s_payment_request.EvCardUID)

		if err != nil {
			s_payment_response.Status = "failed"
			s_payment_response.StatusCode = 500
			s_payment_response.ErrorMessage = "Internal server error"
			return s_payment_response, nil
		}
		if transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Result != "0" {
			s_payment_response.ErrorMessage, s_payment_response.StatusCode = ws.GetErrorMessage(transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Result)
			return s_payment_response, nil
		}
		if transferResponse.XMLProvider.XMLCheckPaymentRequisites.Result != "0" {
			s_payment_response.ErrorMessage, s_payment_response.StatusCode = ws.GetErrorMessage(transferResponse.XMLProvider.XMLCheckPaymentRequisites.Result)
			return s_payment_response, nil
		}

		s_payment_response.Data = &s_payment_response_data_hdr{}

		s_payment_response.Data.Session = transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp1
		s_payment_response.Data.Id = transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Id
		s_payment_response.Data.Info = transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp2
		s_payment_response.Status = "success"

	} else {
		s_payment_response.StatusCode = 403
		s_payment_response.ErrorMessage = "Login/Senha inválido"
		s_payment_response.Status = "failed"
	}
	return s_payment_response, nil
}

//func DoPaymentTransNFC2(s_credentials *db.Login_credentials_hdr, id string, session string, amount string, serviceId string, ev_card_data string) (*WSResponse_transferCredits_hdr, error) {

func paymentNFC2(s_payment_request s_payment_request_hdr) (s_payment_response s_payment_response_hdr, err error) {

	s_payment_response = s_payment_response_hdr{}

	dbConn := db.Connect()

	defer dbConn.Close()

	s_login_credentials, err := db.GetAuthToken(dbConn, s_payment_request.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {

		transferResponse, err := ws.DoPaymentTransNFC2(s_login_credentials, s_payment_request.Id, s_payment_request.Session, s_payment_request.Amount, s_payment_request.Service, s_payment_request.EvCardData)
		if err != nil {
			s_payment_response.Status = "failed"
			s_payment_response.StatusCode = 500
			s_payment_response.ErrorMessage = "Internal server error"
			return s_payment_response, nil
		}
		if transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Result != "0" {
			s_payment_response.ErrorMessage, s_payment_response.StatusCode = ws.GetErrorMessage(transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Result)
			return s_payment_response, nil
		}
		if transferResponse.XMLProvider.XMLCheckPaymentRequisites.Result != "0" {
			s_payment_response.ErrorMessage, s_payment_response.StatusCode = ws.GetErrorMessage(transferResponse.XMLProvider.XMLCheckPaymentRequisites.Result)
			return s_payment_response, nil
		}

		s_payment_response.Data = &s_payment_response_data_hdr{}
		s_payment_response.Data.Info = transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp1
		s_payment_response.Status = "success"

	} else {
		s_payment_response.StatusCode = 403
		s_payment_response.ErrorMessage = "Login/Senha inválido"
		s_payment_response.Status = "failed"
	}
	return s_payment_response, nil
}
func paymentNFC3(s_payment_request s_payment_request_hdr) (s_payment_response s_payment_response_hdr, err error) {

	s_payment_response = s_payment_response_hdr{}

	dbConn := db.Connect()

	defer dbConn.Close()

	s_login_credentials, err := db.GetAuthToken(dbConn, s_payment_request.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {

		transferResponse, err := ws.DoPaymentTransNFC3(s_login_credentials, s_payment_request.Id, s_payment_request.Session, s_payment_request.Amount, s_payment_request.Service, s_payment_request.EvCardData)

		if err != nil {
			s_payment_response.Status = "failed"
			s_payment_response.StatusCode = 500
			s_payment_response.ErrorMessage = "Internal server error"
			return s_payment_response, nil
		}
		if transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Result != "0" {
			s_payment_response.ErrorMessage, s_payment_response.StatusCode = ws.GetErrorMessage(transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Result)
			return s_payment_response, nil
		}
		if transferResponse.XMLProvider.XMLCheckPaymentRequisites.Result != "0" {
			s_payment_response.ErrorMessage, s_payment_response.StatusCode = ws.GetErrorMessage(transferResponse.XMLProvider.XMLCheckPaymentRequisites.Result)
			return s_payment_response, nil
		}

		s_payment_response.Data = &s_payment_response_data_hdr{}

		s_payment_response.Data.Info = transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp1
		s_payment_response.Status = "success"
	} else {
		s_payment_response.StatusCode = 403
		s_payment_response.ErrorMessage = "Login/Senha inválido"
		s_payment_response.Status = "failed"
	}
	return s_payment_response, nil
}
func paymentNFC4(s_payment_request s_payment_request_hdr) (s_payment_response s_payment_response_hdr, err error) {

	s_payment_response = s_payment_response_hdr{}

	dbConn := db.Connect()

	defer dbConn.Close()

	s_login_credentials, err := db.GetAuthToken(dbConn, s_payment_request.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {

		if s_payment_request.Type == "" {
			s_payment_response.StatusCode = 500
			s_payment_response.Status = "failed"
			s_payment_response.ErrorMessage = "Missing type."

			return s_payment_response, nil
		}
		transferResponse, err := ws.DoPaymentTransNFC4(s_login_credentials, s_payment_request.Id, s_payment_request.Session, s_payment_request.Amount, s_payment_request.Service, s_payment_request.EvNominal, s_payment_request.Ev_wallet_code)

		if err != nil {
			s_payment_response.Status = "failed"
			s_payment_response.StatusCode = 500
			s_payment_response.ErrorMessage = "Internal server error"
			return s_payment_response, nil
		}
		if transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Result != "0" {
			s_payment_response.ErrorMessage, s_payment_response.StatusCode = ws.GetErrorMessage(transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Result)
			return s_payment_response, nil
		}
		if transferResponse.XMLProvider.XMLCheckPaymentRequisites.Result != "0" {
			s_payment_response.ErrorMessage, s_payment_response.StatusCode = ws.GetErrorMessage(transferResponse.XMLProvider.XMLCheckPaymentRequisites.Result)
			return s_payment_response, nil
		}

		s_payment_response.Data = &s_payment_response_data_hdr{}

		s_payment_response.Data.Session = transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp1
		s_payment_response.Data.Id = transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Id
		s_payment_response.Data.Info = transferResponse.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Disp2
		s_payment_response.Status = "success"

	} else {
		s_payment_response.StatusCode = 403
		s_payment_response.ErrorMessage = "Login/Senha inválido"
		s_payment_response.Status = "failed"
	}
	return s_payment_response, nil
}
func paymentNFC5(s_payment_request s_payment_request_hdr) (s_payment_response s_payment_response_hdr, err error) {

	s_payment_response = s_payment_response_hdr{}

	dbConn := db.Connect()

	defer dbConn.Close()

	s_login_credentials, err := db.GetAuthToken(dbConn, s_payment_request.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {
		if s_payment_request.Type == "" {
			s_payment_response.StatusCode = 500
			s_payment_response.Status = "failed"
			s_payment_response.ErrorMessage = "Missing type."

			return s_payment_response, nil
		}
		transferResponse, err := ws.DoPaymentTransNFC5(s_login_credentials, s_payment_request.Id, s_payment_request.Session, s_payment_request.Amount, s_payment_request.Service)

		if err != nil {
			s_payment_response.Status = "failed"
			s_payment_response.StatusCode = 500
			s_payment_response.ErrorMessage = "Internal server error"
			return s_payment_response, nil
		}
		if transferResponse.XMLProvider.XMLPurchaseOnline.Result != "0" {
			s_payment_response.ErrorMessage, s_payment_response.StatusCode = ws.GetErrorMessage(transferResponse.XMLProvider.XMLPurchaseOnline.Result)

			return s_payment_response, nil
		}
		if transferResponse.XMLProvider.XMLPurchaseOnline.XMLPayment.Result != "0" {
			s_payment_response.ErrorMessage, s_payment_response.StatusCode = ws.GetErrorMessage(transferResponse.XMLProvider.XMLPurchaseOnline.XMLPayment.Result)

			return s_payment_response, nil
		}

		s_payment_response.Data = &s_payment_response_data_hdr{}

		s_payment_response.Data.VoucherCode = transferResponse.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLVoucher.Code
		s_payment_response.Status = "success"

	} else {
		s_payment_response.StatusCode = 403
		s_payment_response.ErrorMessage = "Login/Senha inválido"
		s_payment_response.Status = "failed"
	}
	return s_payment_response, nil
}
