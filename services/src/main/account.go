// main
package main

import (
	"db"
	"fmt"
	"ws"
)

type s_provider_request_hdr struct {
	AuthToken string `json:"authToken"`
}
type s_history_request_hdr struct {
	AuthToken string `json:"authToken"`
}
type s_provider_response_row_hdr struct {
	FiscalName  string `json:"fiscalName"`
	LongName    string `json:"longName"`
	PrvId       string `json:"prvId"`
	ReceiptName string `json:"receiptName"`
	ShortName   string `json:"shortName"`
}
type s_provider_response_service_hdr struct {
	ServiceName string                        `json:"serviceName"`
	Type        string                        `json:"type"`
	Provider    []s_provider_response_row_hdr `json:"providers"`
}
type s_provider_response_data_hdr struct {
	Services []s_provider_response_service_hdr `json:"services"`
}

type s_provider_response_hdr struct {
	s_status
	Data *s_provider_response_data_hdr `json:"data,omitempty"`
}

type s_history_data_item_hdr struct {
	Timestamp    int    `json:"timestamp"`
	PaymentType  string `json:"paymentType"`
	CategoryName string `json:"categoryName"`
	ServiceName  string `json:"serviceName"`
	Rcpt         string `json:"rcpt"`
	Amount       string `json:"amount"`
}

type s_history_response_data_hdr struct {
	Item []s_history_data_item_hdr `json:"item,omitempty"`
}
type s_history_response_hdr struct {
	s_status
	Data s_history_response_data_hdr `json:"data"`
}

func getHistoryOfUser(s_history_request s_history_request_hdr) (s_history_response_hdr, error) {
	dbConn := db.Connect()
	defer dbConn.Close()

	s_history_response := s_history_response_hdr{}

	s_login_credentials, err := db.GetAuthToken(dbConn, s_history_request.AuthToken)
	if err != nil || s_login_credentials.Id <= 0 {
		s_history_response.StatusCode = 401
		s_history_response.ErrorMessage = "Login Invalido"

		return s_history_response, nil
	}
	list, err := db.ListHistory(dbConn, s_login_credentials.Id)
	if err != nil {
		s_history_response.StatusCode = 500
		s_history_response.ErrorMessage = "Internal server error"

		return s_history_response, nil
	}

	for _, v := range *list {
		item := s_history_data_item_hdr{}
		item.ServiceName = v.ServiceName
		item.Timestamp = v.Timestamp
		item.CategoryName = v.CategoryName
		item.Rcpt = v.Rcpt
		item.Amount = v.Amount
		item.PaymentType = v.PaymentType

		s_history_response.Data.Item = append(s_history_response.Data.Item, item)
	}
	return s_history_response, nil
}
func getProvider(s_provider_request s_provider_request_hdr) (s_provider_response s_provider_response_hdr, err error) {
	//	defer func() {
	//		if r := recover(); r != nil {
	//			fmt.Println("PANIC - ", r)

	//			err = fmt.Errorf("panic")
	//		}
	//	}()

	s_provider_response = s_provider_response_hdr{}
	fmt.Println("GET PROVIDER " + s_provider_request.AuthToken)

	dbConn := db.Connect()
	defer dbConn.Close()

	servicos, err := db.ListServicos(dbConn)

	s_login_credentials, err := db.GetAuthToken(dbConn, s_provider_request.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {
		result, err := ws.GetProvider(s_login_credentials)
		if err != nil {
			s_provider_response.StatusCode = 500
			s_provider_response.ErrorMessage = "Internal server error"
		} else {
			s_provider_response.Data = &s_provider_response_data_hdr{}

			services := []s_provider_response_service_hdr{}

			for k, _ := range result.XMLProvider.XMLGetProvider.Row {
				item := &result.XMLProvider.XMLGetProvider.Row[k]

				for _, s := range *servicos {
					if item.PrvId == s.RvId {
						item.ServiceName = s.Name
						item.Type = s.Type
						item.LongName = s.LongName
					}
				}
			}
			for k, _ := range result.XMLProvider.XMLGetProvider.Row {
				item := &result.XMLProvider.XMLGetProvider.Row[k]

				providers := []s_provider_response_row_hdr{}

				if item.PrvId != "0" && item.ServiceName != "" {
					for k2, _ := range result.XMLProvider.XMLGetProvider.Row {
						item2 := &result.XMLProvider.XMLGetProvider.Row[k2]

						if item.ServiceName == item2.ServiceName {
							fmt.Println("FOUND " + item2.ShortName)

							row := s_provider_response_row_hdr{}

							row.FiscalName = item2.FiscalName
							row.LongName = item2.LongName
							row.ReceiptName = item2.ReceiptName
							row.ShortName = item2.ShortName
							row.PrvId = item2.PrvId

							item2.PrvId = "0"
							providers = append(providers, row)
						}
					}
				}
				if len(providers) > 0 {
					s_service := s_provider_response_service_hdr{}
					s_service.Provider = providers

					s_service.ServiceName = item.ServiceName
					s_service.Type = item.Type
					services = append(services, s_service)

				}
			}

			s_provider_response.Data.Services = services
		}
	} else {
		s_provider_response.StatusCode = 403
		s_provider_response.ErrorMessage = "Login/Senha inválido"
	}
	return s_provider_response, nil
}
