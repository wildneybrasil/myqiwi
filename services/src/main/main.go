// main
package main

import (
	"db"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"redis"

	"net/http"
	"net/url"
	"notification"
	"random"
	"strings"
	"time"
	"ws"

	"golang.org/x/crypto/scrypt"
)

var (
	MAXLOGIN = 99
)

type s_redis_create_account_hdr struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	Name           string `json:"name"`
	Cel            string `json:"cel"`
	Salt           string `json:"salt"`
	ActivationCode string `json:"activationCode"`
	AuthToken      string `json:"authToken"`
}
type s_redis_lost_password_hdr struct {
	Email        string `json:"email"`
	RecoverToken string `json:"recoverToken"`
}

type s_login_request_hdr struct {
	Email    string `json:"email"`
	Cel      string `json:"cel"`
	Password string `json:"password"`
}
type s_activate_request_hdr struct {
	AuthToken      string `json:"authToken"`
	ActivationCode string `json:"activationCode"`
}
type s_login_create_request_hdr struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Cel      string `json:"cel"`
}
type s_login_response_data_hdr struct {
	AuthToken string `json:"authToken,omitempty"`
}
type s_lost_password_hdr struct {
	Email         string `json:"email"`
	Password      string `json:"password"`
	RecoveryToken string `json:"recoveryToken"`
}
type s_login_response_hdr struct {
	s_status
	Data s_login_response_data_hdr `json:"data,omitempty"`
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
type s_transferCredits_request_hdr struct {
	AuthToken string `json:"authToken"`
	Amount    string `json:"amount"`
	Account   string `json:"account"`
	Terminal  string `json:"terminal"`
	Service   string `json:"service"`
	Id        string `json:"id"`
	session   string `json:"session"`
}

type s_createBill_request_hdr struct {
	AuthToken string `json:"authToken"`
	Amount    string `json:"amount"`
}
type s_getBill_request_hdr struct {
	AuthToken string `json:"authToken"`
	BoletoId  string `json:"boletoId"`
}
type s_provider_request_hdr struct {
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
	Provider    []s_provider_response_row_hdr `json:"providers"`
}
type s_provider_response_data_hdr struct {
	Services []s_provider_response_service_hdr `json:"services"`
}

type s_provider_response_hdr struct {
	s_status
	Data *s_provider_response_data_hdr `json:"data,omitempty"`
}

//
type s_createBill_response_data_hdr struct {
	Amount string `json:"amount"`
	Id     string `json:"id,omitempty"`
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
type s_login_create_response_hdr struct {
	s_status
	Data *s_login_response_data_hdr `json:"data,omitempty"`
}
type s_values_item_hdr struct {
	Id     string `json:"id,attr,omitempty"`
	Amount string `json:"amount,attr,omitempty"`
}
type s_values_response_hdr struct {
	Items []s_values_item_hdr `json:"items,omitempty"`
}
type s_transferCredits_response_hdr struct {
	s_status
	Nominals s_values_response_hdr `json:"nominals,omitempty"`
}
type s_transferCredits2_response_hdr struct {
	s_status
	Voucher string `json:"voucher,omitempty"`
	Amount  string `json:"amount,omitempty"`
}

type s_status struct {
	Status       string `json:"status,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	StatusCode   int    `json:"statusCode"`
}

func parseContent(source io.Reader, dest interface{}) error {
	decoder := json.NewDecoder(source)

	err := decoder.Decode(dest)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//func sendJson(w http.ResponseWriter, json interface{}) {
//	resultString, err := json.Marshal(json)

//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	w.Write(resultString)

//}
func main() {

	notification.Connect()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		path := html.EscapeString(r.URL.Path)

		fmt.Println("PATH: " + path)

		switch path {
		case "/ws/login":
			s_login_request := s_login_request_hdr{}

			err := parseContent(r.Body, &s_login_request)
			if err != nil {
				fmt.Println(err.Error())
				w.Write([]byte(err.Error()))
			}

			s_result, err := login(s_login_request)
			if err != nil {
				fmt.Println(err)
				return
			}
			resultString, err := json.Marshal(s_result)

			if err != nil {
				fmt.Println(err)
				return
			}
			w.Write(resultString)

			break
		case "/ws/activateAccount":
			s_activate_request := s_activate_request_hdr{}

			err := parseContent(r.Body, &s_activate_request)
			if err != nil {
				fmt.Println(err.Error())
				w.Write([]byte(err.Error()))
			}

			result, err := activateAccount(s_activate_request)
			if err != nil {
				fmt.Println(err)
				return
			}
			resultString, _ := json.Marshal(result)
			w.Write(resultString)

			break
		case "/ws/resendActivationCode":
			s_activate_request := s_activate_request_hdr{}

			err := parseContent(r.Body, &s_activate_request)
			if err != nil {
				fmt.Println(err.Error())
				w.Write([]byte(err.Error()))
			}

			result, err := resendActivationCode(s_activate_request)
			if err != nil {
				fmt.Println(err)
				return
			}
			resultString, _ := json.Marshal(result)
			w.Write(resultString)

			break

		case "/ws/activateGETAccount":
			//	db.CreateAccount(dbConn, s_login_create_request.Email, s_login_create_request.Cel, dkb64Encoded, salt, s_login_create_request.Name)

			url, _ := url.Parse(r.URL.String())
			dbConn := db.Connect()

			token := url.Query().Get("TOKEN")

			fmt.Println(token)
			status := s_status{}
			status.Status = "failed"
			status.StatusCode = 500
			status.ErrorMessage = "Parametros'"

			if len(token) > 0 {
				userInfo, err := db.GetLoginInfoBySalt(dbConn, token)
				if err != nil {
					status.StatusCode = 403
					status.ErrorMessage = "Token não existe'"
				} else if userInfo.Status == 0 {
					db.ActivateUser(dbConn, token)
					status.StatusCode = 0
					status.Status = "success"
				}
			}
			resultString, err := json.Marshal(status)

			if err != nil {
				fmt.Println(err)
				return
			}
			w.Write(resultString)

			dbConn.Close()
			break
		case "/ws/listServicos":
			s_login_request := s_login_request_hdr{}

			err := parseContent(r.Body, &s_login_request)
			if err != nil {
				fmt.Println(err.Error())
				w.Write([]byte(err.Error()))
			}

			s_result, err := login(s_login_request)
			if err != nil {
				fmt.Println(err)

				resultString, _ := json.Marshal(s_status{"failed", "System error", 500})
				w.Write(resultString)

				return
			}
			resultString, err := json.Marshal(s_result)

			if err != nil {
				fmt.Println(err)
				return
			}
			w.Write(resultString)

			break
		case "/ws/lostPassword":
			s_lost_password := s_lost_password_hdr{}

			err := parseContent(r.Body, &s_lost_password)
			if err != nil {
				fmt.Println(err.Error())
				w.Write([]byte(err.Error()))
			}

			lostPassword(s_lost_password)

			resultString, _ := json.Marshal(s_status{"success", "", 0})

			if err != nil {
				fmt.Println(err)
				return
			}
			w.Write(resultString)

			break
		case "/ws/verifyLPToken":
			s_lost_password := s_lost_password_hdr{}

			err := parseContent(r.Body, &s_lost_password)
			if err != nil {
				fmt.Println(err.Error())
				w.Write([]byte(err.Error()))
			}

			result := verifyLPToken(s_lost_password)

			resultString, _ := json.Marshal(result)

			if err != nil {
				fmt.Println(err)
				return
			}
			w.Write(resultString)

			break
		case "/ws/changeLP":
			s_lost_password := s_lost_password_hdr{}

			err := parseContent(r.Body, &s_lost_password)
			if err != nil {
				fmt.Println(err.Error())
				w.Write([]byte(err.Error()))
			}

			result := changeLPPassword(s_lost_password)

			resultString, _ := json.Marshal(result)

			if err != nil {
				fmt.Println(err)
				return
			}
			w.Write(resultString)

			break
		case "/ws/getBalance":
			s_balance_request := s_balance_request_hdr{}

			err := parseContent(r.Body, &s_balance_request)
			if err != nil {
				fmt.Println(err.Error())
				w.Write([]byte(err.Error()))
			}

			s_result, err := getBalance(s_balance_request)
			if err != nil {
				fmt.Println(err)

				resultString, _ := json.Marshal(s_status{"failed", "System error", 500})
				w.Write(resultString)

				return
			}
			resultString, err := json.Marshal(s_result)

			if err != nil {
				fmt.Println(err)
				return
			}
			w.Write(resultString)

			break
		case "/ws/getProvider":
			s_provider_request := s_provider_request_hdr{}

			err := parseContent(r.Body, &s_provider_request)
			if err != nil {
				fmt.Println(err.Error())
				w.Write([]byte(err.Error()))
			}

			s_result, err := getProvider(s_provider_request)
			if err != nil {
				fmt.Println(err)

				resultString, _ := json.Marshal(s_status{"failed", "System error", 500})
				w.Write(resultString)

				return
			}
			resultString, err := json.Marshal(s_result)

			if err != nil {
				fmt.Println(err)
				return
			}
			w.Write(resultString)

			break
		case "/ws/transferCredits1":
			s_transferCredits_request := s_transferCredits_request_hdr{}

			err := parseContent(r.Body, &s_transferCredits_request)
			if err != nil {
				fmt.Println(err.Error())
				w.Write([]byte(err.Error()))
			}

			s_result, err := transferCredits1(s_transferCredits_request)
			if err != nil {
				fmt.Println(err)

				resultString, _ := json.Marshal(s_status{"failed", "System error", 500})
				w.Write(resultString)

				return
			}
			resultString, err := json.Marshal(s_result)

			if err != nil {
				fmt.Println(err)
				return
			}
			w.Write(resultString)

			break
		case "/ws/transferCredits2":
			s_transferCredits_request := s_transferCredits_request_hdr{}

			err := parseContent(r.Body, &s_transferCredits_request)
			if err != nil {
				fmt.Println(err.Error())
				w.Write([]byte(err.Error()))
			}

			s_result, err := transferCredits2(s_transferCredits_request)
			if err != nil {
				fmt.Println(err)

				resultString, _ := json.Marshal(s_status{"failed", "System error", 500})
				w.Write(resultString)

				return
			}
			resultString, err := json.Marshal(s_result)

			if err != nil {
				fmt.Println(err)
				return
			}
			w.Write(resultString)

			break

		case "/ws/createBill":
			s_createBill_request := s_createBill_request_hdr{}

			err := parseContent(r.Body, &s_createBill_request)
			if err != nil {
				fmt.Println(err.Error())
				w.Write([]byte(err.Error()))
			}

			s_result, err := createBill(s_createBill_request)
			if err != nil {
				fmt.Println(err)
				return
			}
			resultString, err := json.Marshal(s_result)

			if err != nil {
				fmt.Println(err)
				return
			}
			w.Write(resultString)

			break
		case "/ws/createAccount":
			s_login_create_request := s_login_create_request_hdr{}

			err := parseContent(r.Body, &s_login_create_request)
			if err != nil {
				fmt.Println(err.Error())
				w.Write([]byte(err.Error()))
			}

			s_result, err := createLogin(s_login_create_request)
			if err != nil {
				fmt.Println(err)
				return
			}
			resultString, err := json.Marshal(s_result)

			if err != nil {
				fmt.Println(err)
				return
			}
			w.Write(resultString)

			break
		case "/ws/getBoletoImage":
			s_getBill_request := s_getBill_request_hdr{}

			err := parseContent(r.Body, &s_getBill_request)
			if err != nil {
				fmt.Println(err.Error())
				w.Write([]byte(err.Error()))
			}

			s_result, err := getBillImage(s_getBill_request)
			if err != nil {
				fmt.Println(err)
				return
			}
			resultString, err := json.Marshal(s_result)

			if err != nil {
				fmt.Println(err)
				return
			}
			w.Write(resultString)

			break

		default:
			w.WriteHeader(500)
			w.Write([]byte("Not Found"))
			break
		}
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
func createLogin(s_login_create_request s_login_create_request_hdr) (s_login_create_response s_login_create_response_hdr, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("PANIC - ", r)

			err = fmt.Errorf("panic")
		}
	}()

	activationCode := random.RandomNumberString(6)
	authToken := random.RandomString(64)
	salt := random.RandomString(32)

	s_login_create_response = s_login_create_response_hdr{}

	dbConn := db.Connect()
	defer dbConn.Close()

	dk, err := scrypt.Key([]byte(s_login_create_request.Password), []byte(salt), 16384, 8, 1, 32)
	if err != nil {
		s_login_create_response.StatusCode = 500
		s_login_create_response.ErrorMessage = "Internal server error"

		return s_login_create_response, nil
	}
	dkb64Encoded := b64.StdEncoding.EncodeToString([]byte(dk))

	fmt.Println(dkb64Encoded)

	s_redis := s_redis_create_account_hdr{}
	s_redis.Name = s_login_create_request.Name
	s_redis.Cel = s_login_create_request.Cel
	s_redis.Email = s_login_create_request.Email
	s_redis.Password = dkb64Encoded
	s_redis.Salt = salt
	s_redis.AuthToken = authToken
	s_redis.ActivationCode = activationCode

	redisString, _ := json.Marshal(s_redis)
	redis.Set(s_redis.AuthToken, string(redisString), 10*time.Minute)

	notification.Send(notification.NotificationMessage{"sms", s_login_create_request.Cel, "QIWI - Seu codigo de ativação é: " + activationCode})

	s_login_create_response.Data = &s_login_response_data_hdr{}
	s_login_create_response.Data.AuthToken = authToken

	return s_login_create_response, nil
}

func login(s_login_request s_login_request_hdr) (s_login_response s_login_response_hdr, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("PANIC - ", r)

			err = fmt.Errorf("panic")
		}
	}()
	authToken := random.RandomString(64)

	s_login_response = s_login_response_hdr{}
	s_login_response.Status = "failed"

	dbConn := db.Connect()
	defer dbConn.Close()

	s_login_credentials, err := db.GetLoginInfoByEmail(dbConn, s_login_request.Email)
	if err != nil {
		s_login_response.StatusCode = 403
		s_login_response.ErrorMessage = "Login/Senha inválido"

		return s_login_response, nil
	}
	if s_login_credentials.FailedLoginCount > MAXLOGIN {
		s_login_response.StatusCode = 403
		s_login_response.ErrorMessage = "Tentativas excedidas"

		return s_login_response, nil
	}
	if s_login_credentials.Status == 0 {
		s_login_response.StatusCode = 403
		s_login_response.ErrorMessage = "Usuario não ativado"

		return s_login_response, nil
	}

	dk, err := scrypt.Key([]byte(s_login_request.Password), []byte(s_login_credentials.PasswordSalt), 16384, 8, 1, 32)
	if err != nil {
		s_login_response.StatusCode = 403
		s_login_response.ErrorMessage = "Login/Senha inválido."

		return s_login_response, nil
	}
	dkb64Encoded := b64.StdEncoding.EncodeToString([]byte(dk))

	fmt.Println(dkb64Encoded)

	if dkb64Encoded == s_login_credentials.Password {
		db.ResetFailedLoginOfEmail(dbConn, s_login_request.Email)
		db.InsertToken(dbConn, s_login_credentials.Id, authToken)

		s_login_response.Status = "success"
		s_login_response.Data.AuthToken = authToken
	} else {
		db.IncreaseFailedLoginOfEmail(dbConn, s_login_request.Email)

		s_login_response.StatusCode = 403
		s_login_response.ErrorMessage = "Login/Senha inválido"

	}

	return s_login_response, nil
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
func getProvider(s_provider_request s_provider_request_hdr) (s_provider_response s_provider_response_hdr, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("PANIC - ", r)

			err = fmt.Errorf("panic")
		}
	}()

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

			for k, _ := range result.XMLProvider.XMLGetProvider.Row {
				item := &result.XMLProvider.XMLGetProvider.Row[k]
				serviceInfo := db.FindServiceByLongName(servicos, strings.TrimSpace(item.LongName))
				if serviceInfo != nil {
					item.ServiceName = strings.TrimSpace(serviceInfo.Name)
				}
			}
			services := []s_provider_response_service_hdr{}

			for k, _ := range result.XMLProvider.XMLGetProvider.Row {
				item := &result.XMLProvider.XMLGetProvider.Row[k]

				if item.ServiceName != "" {
					providers := []s_provider_response_row_hdr{}

					serviceName := item.ServiceName
					for k2, _ := range result.XMLProvider.XMLGetProvider.Row {
						item2 := &result.XMLProvider.XMLGetProvider.Row[k2]

						if item.ServiceName != "" && strings.TrimSpace(item.ServiceName) == strings.TrimSpace(item2.ServiceName) {
							row := s_provider_response_row_hdr{}

							row.FiscalName = item2.FiscalName
							row.LongName = item2.LongName
							row.ReceiptName = item2.ReceiptName
							row.ShortName = item2.ShortName
							row.PrvId = item2.PrvId

							providers = append(providers, row)
							//							item2.ServiceName = ""
						}
					}
					if len(providers) > 0 {
						s_service := s_provider_response_service_hdr{}
						s_service.ServiceName = serviceName
						s_service.Provider = providers

						services = append(services, s_service)
					}
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

	s_login_credentials, err := db.GetAuthToken(dbConn, s_transferCredits_request.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {
		transferResponse, err := ws.TransferCredits1(s_login_credentials, s_transferCredits_request.Account, s_transferCredits_request.Terminal, s_transferCredits_request.Service, s_transferCredits_request.Amount)
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

				s_transferCredits_response.Nominals.Items = append(s_transferCredits_response.Nominals.Items, item)
			}
		}
	} else {
		s_transferCredits_response.StatusCode = 403
		s_transferCredits_response.ErrorMessage = "Login/Senha inválido"
	}
	return s_transferCredits_response, nil
}
func activateAccount(s_request s_activate_request_hdr) (result s_status, err error) {
	defer func() {
		if r := recover(); r != nil {
			result.Status = "failed"
			result.StatusCode = 404
			result.ErrorMessage = "AuthToken invalido"
			err = nil
		}
	}()

	dbConn := db.Connect()

	defer dbConn.Close()

	s_redis := s_redis_create_account_hdr{}

	redisString, err := redis.Get(s_request.AuthToken)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(*redisString), &s_redis)
	if err != nil {
		panic(err)
	}
	if s_request.ActivationCode == s_redis.ActivationCode {
		id, err := db.CreateAccount(dbConn, s_redis.Email, s_redis.Cel, s_redis.Password, s_redis.Salt, s_redis.Name)
		if err != nil {
			panic(err)
		}
		result.Status = "success"
		result.StatusCode = 0

		db.InsertToken(dbConn, id, s_request.AuthToken)

		redis.Del(s_request.AuthToken)

	} else {
		result.Status = "failed"
		result.StatusCode = 404
		result.ErrorMessage = "Código de ativação invalido"
	}
	return result, nil
}

func resendActivationCode(s_request s_activate_request_hdr) (result s_status, err error) {
	defer func() {
		if r := recover(); r != nil {
			result.Status = "failed"
			result.StatusCode = 404
			result.ErrorMessage = "AuthToken invalido"
			err = nil
		}
	}()

	dbConn := db.Connect()

	defer dbConn.Close()

	s_redis := s_redis_create_account_hdr{}

	redisString, err := redis.Get(s_request.AuthToken)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(*redisString), &s_redis)
	if err != nil {
		panic(err)
	}

	redisBytes, _ := json.Marshal(s_redis)

	redis.Set(s_redis.AuthToken, string(redisBytes), 10*time.Minute)

	notification.Send(notification.NotificationMessage{"sms", s_redis.Cel, "QIWI - Seu codigo de ativação é: " + s_redis.ActivationCode})

	return result, nil
}

func transferCredits2(s_transferCredits_request s_transferCredits_request_hdr) (s_transferCredits_response s_transferCredits2_response_hdr, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("PANIC - ", r)

			err = fmt.Errorf("panic")
		}
	}()

	s_transferCredits_response = s_transferCredits2_response_hdr{}

	dbConn := db.Connect()

	defer dbConn.Close()

	s_login_credentials, err := db.GetAuthToken(dbConn, s_transferCredits_request.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {
		transferResponse, err := ws.TransferCredits2(s_login_credentials, s_transferCredits_request.Id, s_transferCredits_request.session, s_transferCredits_request.Account, s_transferCredits_request.Terminal, s_transferCredits_request.Service, s_transferCredits_request.Amount)
		if err != nil {
			s_transferCredits_response.StatusCode = 500
			s_transferCredits_response.ErrorMessage = "Internal server error"
		} else {
			fmt.Println(transferResponse.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLVoucher.Code)
			s_transferCredits_response.Voucher = transferResponse.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLVoucher.Code
			s_transferCredits_response.Amount = transferResponse.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLVoucher.Amount
		}
	} else {
		s_transferCredits_response.StatusCode = 403
		s_transferCredits_response.ErrorMessage = "Login/Senha inválido"
	}
	return s_transferCredits_response, nil
}
func lostPassword(s_lost_password s_lost_password_hdr) {
	recoverToken := random.RandomString(64)

	s_redis := s_redis_lost_password_hdr{}
	s_redis.Email = s_lost_password.Email
	s_redis.RecoverToken = recoverToken

	redisString, _ := json.Marshal(s_redis)
	redis.Set(s_redis.Email, string(redisString), 10*time.Minute)

	notification.Send(notification.NotificationMessage{"email", s_lost_password.Email, "\n\nhttp://ec2-54-207-24-178.sa-east-1.compute.amazonaws.com/password/#/password/" + s_lost_password.Email + "/" + recoverToken})

}
func changeLPPassword(s_lost_password s_lost_password_hdr) (result s_status) {
	defer func() {
		if r := recover(); r != nil {
			result.Status = "failed"
			result.StatusCode = 0
			result.ErrorMessage = "AuthToken invalido"
		}
	}()
	salt := random.RandomString(32)
	s_redis_lost_password := s_redis_lost_password_hdr{}

	redisString, err := redis.Get(s_lost_password.Email)

	if err != nil || redisString == nil {
		return s_status{"failed", "Token não encontrado", 404}
	}
	err = json.Unmarshal([]byte(*redisString), &s_redis_lost_password)
	if err != nil {
		panic(err)
	}
	if s_redis_lost_password.RecoverToken != s_lost_password.RecoveryToken {
		return s_status{"failed", "Token não encontrado", 404}
	}

	dk, err := scrypt.Key([]byte(s_lost_password.Password), []byte(salt), 16384, 8, 1, 32)
	if err != nil {
		return s_status{"failed", "Internal server Error", 500}
	}
	dkb64Encoded := b64.StdEncoding.EncodeToString([]byte(dk))

	dbConn := db.Connect()
	defer dbConn.Close()
	db.ChangePassword(dbConn, s_redis_lost_password.Email, dkb64Encoded, salt)

	redis.Del(s_lost_password.Email)

	return s_status{"success", "", 0}

}

func verifyLPToken(s_lost_password s_lost_password_hdr) (result s_status) {
	defer func() {
		if r := recover(); r != nil {
			result.Status = "failed"
			result.StatusCode = 0
			result.ErrorMessage = "AuthToken invalido"
		}
	}()

	s_redis_lost_password := s_redis_lost_password_hdr{}

	redisString, err := redis.Get(s_lost_password.Email)
	if err != nil {
		panic(err)
	}
	if redisString != nil {
		return s_status{"success", "", 0}
	} else {
		return s_status{"failed", "Token não encontrado", 404}

	}

	err = json.Unmarshal([]byte(*redisString), &s_redis_lost_password)
	if err != nil {
		panic(err)
	}
	if s_redis_lost_password.RecoverToken != s_lost_password.RecoveryToken {
		return s_status{"failed", "Token não encontrado", 404}
	}

	return s_status{"success", "", 0}

}
