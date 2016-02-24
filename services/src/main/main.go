// main
package main

import (
	"db"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"

	"net/http"
	"net/url"
	"notification"
)

var (
	MAXLOGIN = 99
)

//

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
		case "/ws/accountInfo":
			s_cel_info := s_cel_info_hdr{}

			err := parseContent(r.Body, &s_cel_info)
			if err != nil {
				fmt.Println(err.Error())
				w.Write([]byte(err.Error()))
			}

			s_result, err := getPublicLoginInfoByCel(s_cel_info)
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
		case "/ws/profile":
			s_cel_info := s_cel_info_hdr{}

			err := parseContent(r.Body, &s_cel_info)
			if err != nil {
				fmt.Println(err.Error())
				w.Write([]byte(err.Error()))
			}

			s_result, err := getMyInfo(s_cel_info)
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
		case "/ws/getHistory":
			s_history_request := s_history_request_hdr{}

			err := parseContent(r.Body, &s_history_request)
			if err != nil {
				fmt.Println(err.Error())
				w.Write([]byte(err.Error()))
			}

			s_result, err := getHistoryOfUser(s_history_request)
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
		case "/ws/getBoletoInfo":
			s_boletoInfo_request := s_boletoInfo_request_hdr{}

			err := parseContent(r.Body, &s_boletoInfo_request)
			if err != nil {
				fmt.Println(err.Error())
				w.Write([]byte(err.Error()))
			}

			s_result, err := getBoletoInfo(s_boletoInfo_request)
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

		case "/ws/pay1":
			s_payment_request := s_payment_request_hdr{}

			err := parseContent(r.Body, &s_payment_request)
			if err != nil {
				fmt.Println(err.Error())
				w.Write([]byte(err.Error()))
			}

			s_result, err := payment1(s_payment_request)
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

		case "/ws/pay2":
			s_payment_request := s_payment_request_hdr{}

			err := parseContent(r.Body, &s_payment_request)
			if err != nil {
				fmt.Println(err.Error())
				w.Write([]byte(err.Error()))
			}

			s_result, err := payment2(s_payment_request)
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
		case "/ws/updateAccount":
			s_login_update_request := s_login_update_request_hdr{}

			err := parseContent(r.Body, &s_login_update_request)
			if err != nil {
				fmt.Println(err.Error())
				w.Write([]byte(err.Error()))
			}

			s_result, err := updateUser(s_login_update_request)
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
