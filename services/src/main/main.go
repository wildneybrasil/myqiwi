// main
package main

import (
	"db"
	"encoding/json"
	"fmt"
	//	"html"
	"io"
	//	"log"

	"github.com/gin-gonic/gin"

	"net/http"
	//	"net/url"
	"notification"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/thirdparty/tollbooth_gin"
)

var (
	MAXLOGIN = 15
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
	r := gin.New()

	limiterGeneric := tollbooth.NewLimiter(1000, time.Minute)
	limiterGeneric.Methods = []string{"GET", "POST"}
	limiterGeneric.Message = "{ 'status':'failed', 'errorCode':400, 'errorMessage':'Aguarde alguns segundos e repita a operação' }"

	limiterAuth := tollbooth.NewLimiter(50, time.Minute)
	limiterAuth.Methods = []string{"GET", "POST"}
	limiterAuth.Message = "{ 'status':'failed', 'errorCode':400, 'errorMessage':'Aguarde alguns segundos e repita a operação' }"

	limiterNotification := tollbooth.NewLimiter(20, time.Minute)
	limiterNotification.Methods = []string{"GET", "POST"}
	limiterNotification.Message = "{ 'status':'failed', 'errorCode':400, 'errorMessage':'Aguarde alguns segundos e repita a operação' }"

	notification.Connect()

	r.POST("/ws/resendActivationCode", tollbooth_gin.LimitHandler(limiterNotification), func(c *gin.Context) {
		s_activate_request := s_activate_request_hdr{}
		if c.BindJSON(&s_activate_request) == nil {
			result, err := resendActivationCode(s_activate_request)
			if err != nil {
				fmt.Println(err)
				return
			}

			c.JSON(http.StatusOK, result)
		}
	})
	r.POST("/ws/createAccount", tollbooth_gin.LimitHandler(limiterNotification), func(c *gin.Context) {
		s_login_create_request := s_login_create_request_hdr{}
		if c.BindJSON(&s_login_create_request) == nil {
			s_result, err := createLogin(s_login_create_request)
			if err != nil {
				fmt.Println(err)
				return
			}
			c.JSON(http.StatusOK, s_result)
		}
	})
	r.POST("/ws/createPOSAccount", tollbooth_gin.LimitHandler(limiterNotification), func(c *gin.Context) {
		s_login_create_request := s_login_pos_create_request_hdr{}
		if c.BindJSON(&s_login_create_request) == nil {
			s_result, err := createPOSLogin(s_login_create_request)
			if err != nil {
				fmt.Println(err)
				return
			}
			c.JSON(http.StatusOK, s_result)
		}
	})
	r.POST("/ws/createBill", tollbooth_gin.LimitHandler(limiterNotification), func(c *gin.Context) {
		s_createBill_request := s_createBill_request_hdr{}
		if c.BindJSON(&s_createBill_request) == nil {
			s_result, err := createBill(s_createBill_request)
			if err != nil {
				fmt.Println(err)
				return
			}
			c.JSON(http.StatusOK, s_result)
		}
	})

	r.POST("/ws/lostPassword", tollbooth_gin.LimitHandler(limiterNotification), func(c *gin.Context) {
		s_lost_password := s_lost_password_hdr{}
		if c.BindJSON(&s_lost_password) == nil {
			lostPassword(s_lost_password)

			c.JSON(http.StatusOK, s_status{"success", "", 0})
		}
	})
	r.Static("/ws/image/logo", "/var/logos")

	r.POST("/ws/activateAccount", tollbooth_gin.LimitHandler(limiterNotification), func(c *gin.Context) {
		s_activate_request := s_activate_request_hdr{}
		if c.BindJSON(&s_activate_request) == nil {
			result, err := activateAccount(s_activate_request)
			if err != nil {
				fmt.Println(err)
				return
			}
			c.JSON(http.StatusOK, result)
		}
	})
	r.POST("/ws/login", tollbooth_gin.LimitHandler(limiterAuth), func(c *gin.Context) {
		s_login_request := s_login_request_hdr{}
		if c.BindJSON(&s_login_request) == nil {
			s_result, err := login(s_login_request)
			if err != nil {
				fmt.Println(err)
				return
			}
			c.JSON(http.StatusOK, s_result)
		}
	})
	r.POST("/ws/accountInfo", tollbooth_gin.LimitHandler(limiterGeneric), func(c *gin.Context) {
		s_cel_info := s_cel_info_hdr{}
		if c.BindJSON(&s_cel_info) == nil {
			s_result, err := getPublicLoginInfoByCel(s_cel_info)
			if err != nil {
				fmt.Println(err)
				return
			}
			c.JSON(http.StatusOK, s_result)
		}
	})
	r.POST("/ws/profile", tollbooth_gin.LimitHandler(limiterGeneric), func(c *gin.Context) {
		s_cel_info := s_cel_info_hdr{}
		if c.BindJSON(&s_cel_info) == nil {
			s_result, err := getMyInfo(s_cel_info)
			if err != nil {
				fmt.Println(err)
				return
			}

			c.JSON(http.StatusOK, s_result)
		}
	})
	r.POST("/ws/contact", tollbooth_gin.LimitHandler(limiterGeneric), func(c *gin.Context) {
		s_contact_request := s_contact_request_hdr{}
		if c.BindJSON(&s_contact_request) == nil {
			s_result := contactUs(s_contact_request)

			c.JSON(http.StatusOK, s_result)
		}
	})
	r.POST("/ws/getHistory", tollbooth_gin.LimitHandler(limiterGeneric), func(c *gin.Context) {
		s_history_request := s_history_request_hdr{}
		if c.BindJSON(&s_history_request) == nil {
			s_result, err := getHistoryOfUser(s_history_request)
			if err != nil {
				fmt.Println(err)
				return
			}
			c.JSON(http.StatusOK, s_result)
		}
	})
	r.GET("/ws/activateGETAccount", tollbooth_gin.LimitHandler(limiterGeneric), func(c *gin.Context) {
		dbConn := db.Connect()
		defer dbConn.Close()

		token := c.Query("TOKEN")

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
		c.JSON(http.StatusOK, status)

	})
	r.POST("/ws/listServicos", tollbooth_gin.LimitHandler(limiterGeneric), func(c *gin.Context) {
		s_login_request := s_login_request_hdr{}
		if c.BindJSON(&s_login_request) == nil {
			s_result, err := login(s_login_request)
			if err != nil {
				fmt.Println(err)

				c.JSON(http.StatusOK, s_status{"failed", "System error", 500})

				return
			}
			c.JSON(http.StatusOK, s_result)

		}
	})
	r.POST("/ws/verifyLPToken", tollbooth_gin.LimitHandler(limiterGeneric), func(c *gin.Context) {
		s_lost_password := s_lost_password_hdr{}
		if c.BindJSON(&s_lost_password) == nil {
			result := verifyLPToken(s_lost_password)

			c.JSON(http.StatusOK, result)

		}
	})
	r.POST("/ws/changeLP", tollbooth_gin.LimitHandler(limiterGeneric), func(c *gin.Context) {
		s_lost_password := s_lost_password_hdr{}
		if c.BindJSON(&s_lost_password) == nil {

			result := changeLPPassword(s_lost_password)

			c.JSON(http.StatusOK, result)
		}
	})
	r.POST("/ws/getBalance", tollbooth_gin.LimitHandler(limiterGeneric), func(c *gin.Context) {
		s_balance_request := s_balance_request_hdr{}

		if c.BindJSON(&s_balance_request) == nil {
			s_result, err := getBalance(s_balance_request)
			if err != nil {
				fmt.Println(err)

				c.JSON(http.StatusOK, s_status{"failed", "System error", 500})

				return
			}
			c.JSON(http.StatusOK, s_result)
		}
	})
	r.POST("/ws/getProvider", tollbooth_gin.LimitHandler(limiterGeneric), func(c *gin.Context) {
		s_provider_request := s_provider_request_hdr{}
		if c.BindJSON(&s_provider_request) == nil {
			s_result, err := getProvider(s_provider_request)
			if err != nil {
				fmt.Println(err)

				c.JSON(http.StatusOK, s_status{"failed", "System error", 500})
				return
			}

			c.JSON(http.StatusOK, s_result)
		}
	})
	r.POST("/ws/getBoletoInfo", tollbooth_gin.LimitHandler(limiterGeneric), func(c *gin.Context) {
		s_boletoInfo_request := s_boletoInfo_request_hdr{}
		if c.BindJSON(&s_boletoInfo_request) == nil {

			s_result, err := getBoletoInfo(s_boletoInfo_request)
			if err != nil {
				c.JSON(http.StatusOK, s_status{"failed", "System error", 500})

				return
			}
			c.JSON(http.StatusOK, s_result)
		}
	})
	r.POST("/ws/transferCredits1", tollbooth_gin.LimitHandler(limiterGeneric), func(c *gin.Context) {
		s_transferCredits_request := s_transferCredits_request_hdr{}
		if c.BindJSON(&s_transferCredits_request) == nil {
			s_result, err := transferCredits1(s_transferCredits_request)
			if err != nil {
				fmt.Println(err)

				c.JSON(http.StatusOK, s_status{"failed", "System error", 500})
				return
			}
			c.JSON(http.StatusOK, s_result)
		}
	})

	r.POST("/ws/pay1", tollbooth_gin.LimitHandler(limiterGeneric), func(c *gin.Context) {
		s_payment_request := s_payment_request_hdr{}
		if c.BindJSON(&s_payment_request) == nil {
			s_result, err := payment1(s_payment_request)
			if err != nil {
				fmt.Println(err)

				c.JSON(http.StatusOK, s_status{"failed", "System error", 500})

				return
			}
			c.JSON(http.StatusOK, s_result)
		}
	})
	r.POST("/ws/pay2", tollbooth_gin.LimitHandler(limiterGeneric), func(c *gin.Context) {
		s_payment_request := s_payment_request_hdr{}
		if c.BindJSON(&s_payment_request) == nil {
			s_result, err := payment2(s_payment_request)
			if err != nil {
				fmt.Println(err)

				c.JSON(http.StatusOK, s_status{"failed", "System error", 500})

				return
			}
			c.JSON(http.StatusOK, s_result)
		}
	})
	r.POST("/ws/updateAccount", tollbooth_gin.LimitHandler(limiterGeneric), func(c *gin.Context) {
		s_login_update_request := s_login_update_request_hdr{}
		if c.BindJSON(&s_login_update_request) == nil {
			s_result, err := updateUser(s_login_update_request)
			if err != nil {
				fmt.Println(err)
				c.JSON(http.StatusOK, s_status{"failed", "System error", 500})
				return
			}
			c.JSON(http.StatusOK, s_result)
		}
	})
	r.POST("/ws/paymentNFC/1", tollbooth_gin.LimitHandler(limiterGeneric), func(c *gin.Context) {
		s_payment_request := s_payment_request_hdr{}
		if c.BindJSON(&s_payment_request) == nil {
			s_result, err := paymentNFC1(s_payment_request)
			if err != nil {
				fmt.Println(err)
				c.JSON(http.StatusOK, s_status{"failed", "System error", 500})
				return
			}
			c.JSON(http.StatusOK, s_result)
		}
	})
	r.POST("/ws/paymentNFC/2", tollbooth_gin.LimitHandler(limiterGeneric), func(c *gin.Context) {
		s_payment_request := s_payment_request_hdr{}
		if c.BindJSON(&s_payment_request) == nil {
			s_result, err := paymentNFC2(s_payment_request)
			if err != nil {
				fmt.Println(err)
				c.JSON(http.StatusOK, s_status{"failed", "System error", 500})
				return
			}
			c.JSON(http.StatusOK, s_result)
		}
	})
	r.POST("/ws/paymentNFC/3", tollbooth_gin.LimitHandler(limiterGeneric), func(c *gin.Context) {
		s_payment_request := s_payment_request_hdr{}
		if c.BindJSON(&s_payment_request) == nil {
			s_result, err := paymentNFC3(s_payment_request)
			if err != nil {
				fmt.Println(err)
				c.JSON(http.StatusOK, s_status{"failed", "System error", 500})
				return
			}
			c.JSON(http.StatusOK, s_result)
		}
	})
	r.POST("/ws/paymentNFC/5", tollbooth_gin.LimitHandler(limiterGeneric), func(c *gin.Context) {
		s_payment_request := s_payment_request_hdr{}
		if c.BindJSON(&s_payment_request) == nil {
			s_result, err := paymentNFC4(s_payment_request)
			if err != nil {
				fmt.Println(err)
				c.JSON(http.StatusOK, s_status{"failed", "System error", 500})
				return
			}
			c.JSON(http.StatusOK, s_result)
		}
	})
	r.POST("/ws/paymentNFC/5", tollbooth_gin.LimitHandler(limiterGeneric), func(c *gin.Context) {
		s_payment_request := s_payment_request_hdr{}
		if c.BindJSON(&s_payment_request) == nil {
			s_result, err := paymentNFC5(s_payment_request)
			if err != nil {
				fmt.Println(err)
				c.JSON(http.StatusOK, s_status{"failed", "System error", 500})
				return
			}
			c.JSON(http.StatusOK, s_result)
		}
	})
	r.POST("/ws/getBillInfo", tollbooth_gin.LimitHandler(limiterGeneric), func(c *gin.Context) {
		s_getBill_request := s_getBill_request_hdr{}
		if c.BindJSON(&s_getBill_request) == nil {

			s_result, err := getBillInfo(s_getBill_request)
			if err != nil {
				fmt.Println(err)
				c.JSON(http.StatusOK, s_status{"failed", "System error", 500})
				return
			}
			c.JSON(http.StatusOK, s_result)
		}
	})

	r.Run(":8081")
}
