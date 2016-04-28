// main
package main

import (
	"crypto/md5"
	"db"
	b64 "encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"notification"
	"random"
	"redis"
	"time"
	"ws"

	"golang.org/x/crypto/scrypt"
)

type s_login_create_response_hdr struct {
	s_status
	Data *s_login_response_data_hdr `json:"data,omitempty"`
}

type s_redis_create_account_hdr struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	Name           string `json:"name"`
	Cel            string `json:"cel"`
	Document       string `json:"document"`
	Photo          string `json:"photo"`
	Salt           string `json:"salt"`
	ActivationCode string `json:"activationCode"`
	AuthToken      string `json:"authToken"`
}
type s_redis_lost_password_hdr struct {
	Email        string `json:"email"`
	RecoverToken string `json:"recoverToken"`
}
type s_cel_info_hdr struct {
	AuthToken string `json:"authToken"`
	Cel       string `json:"cel"`
}
type s_contact_hdr struct {
	AuthToken string `json:"authToken"`
	Text      string `json:"text"`
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
type s_login_info_response_data_hdr struct {
	Name       string `json:"name",omitempty`
	Cel        string `json:"cel",omitempty`
	Photo      string `json:"photo",omitempty`
	Email      string `json:"email",omitempty`
	Document   string `json:"document",omitempty`
	TerminalId string `json:"terminalId",omitempty`
}
type s_login_info_response_hdr struct {
	s_status
	Data s_login_info_response_data_hdr `json:"data,omitempty"`
}
type s_login_update_request_hdr struct {
	AuthToken string `json:"authToken"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Name      string `json:"name"`
	Cel       string `json:"cel"`
	Photo     string `json:"photo"`
}
type s_login_create_request_hdr struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Document string `json:"document"`
	Cel      string `json:"cel"`
	Photo    string `json:"photo"`
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

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
func createLogin(s_login_create_request s_login_create_request_hdr) (s_login_create_response s_login_create_response_hdr, err error) {
	activationCode := random.RandomNumberString(6)
	authToken := random.RandomString(64)

	s_login_create_response = s_login_create_response_hdr{}

	if len(s_login_create_request.Cel) == 0 {
		s_login_create_response.Status = "failed"
		s_login_create_response.StatusCode = 400
		s_login_create_response.ErrorMessage = "Faltando numero de celular"

		return s_login_create_response, nil
	}
	if len(s_login_create_request.Document) == 0 {
		s_login_create_response.Status = "failed"
		s_login_create_response.StatusCode = 400
		s_login_create_response.ErrorMessage = "Faltando numero do documento"

		return s_login_create_response, nil
	}
	if len(s_login_create_request.Email) == 0 {
		s_login_create_response.Status = "failed"
		s_login_create_response.StatusCode = 400
		s_login_create_response.ErrorMessage = "Faltando o email"

		return s_login_create_response, nil
	}
	if len(s_login_create_request.Password) == 0 {
		s_login_create_response.Status = "failed"
		s_login_create_response.StatusCode = 400
		s_login_create_response.ErrorMessage = "Faltando a senha"

		return s_login_create_response, nil
	}
	if len(s_login_create_request.Name) == 0 {
		s_login_create_response.Status = "failed"
		s_login_create_response.StatusCode = 400
		s_login_create_response.ErrorMessage = "Faltando o nome"

		return s_login_create_response, nil
	}

	dbConn := db.Connect()
	defer dbConn.Close()

	// verifca se login ja existe
	vrfyLogin, err := db.GetLoginInfoByEmail(dbConn, s_login_create_request.Email)
	if vrfyLogin != nil {
		s_login_create_response.Status = "failed"
		s_login_create_response.StatusCode = 400
		s_login_create_response.ErrorMessage = "Email ja cadastrado."

		return s_login_create_response, nil
	}

	//	result, err := ws.CreateAccount(s_login_create_request.Name, s_login_create_request.Email, "", s_login_create_request.Cel, s_login_create_request.Password)
	//	if err != nil {
	//		s_login_create_response.Status = "failed"
	//		s_login_create_response.StatusCode = 500
	//		s_login_create_response.ErrorMessage = "Internal server error"

	//		return s_login_create_response, nil
	//	}

	s_redis := s_redis_create_account_hdr{}
	s_redis.Name = s_login_create_request.Name
	s_redis.Cel = s_login_create_request.Cel
	s_redis.Email = s_login_create_request.Email
	s_redis.Photo = s_login_create_request.Photo
	s_redis.Document = s_login_create_request.Document
	s_redis.Password = s_login_create_request.Password
	s_redis.Salt = ""
	s_redis.AuthToken = authToken
	s_redis.ActivationCode = activationCode

	redisString, _ := json.Marshal(s_redis)
	redis.Set(s_redis.AuthToken, string(redisString), 1000*time.Minute)

	notification.Send(notification.NotificationMessage{"sms", ws.GetTelRAW(s_login_create_request.Cel), "QIWI - Seu codigo de ativação é: " + activationCode})

	s_login_create_response.Data = &s_login_response_data_hdr{}
	s_login_create_response.Data.AuthToken = authToken

	return s_login_create_response, nil
}

func updateUser(s_login_update_request s_login_update_request_hdr) (s_login_create_response s_login_create_response_hdr, err error) {
	s_login_create_response = s_login_create_response_hdr{}

	dbConn := db.Connect()
	defer dbConn.Close()

	s_login_credentials, err := db.GetAuthToken(dbConn, s_login_update_request.AuthToken)
	if err != nil || s_login_credentials.Id == 0 {
		s_login_create_response.Status = "failed"
		s_login_create_response.StatusCode = 403
		s_login_create_response.ErrorMessage = "Invalid Token"

		return s_login_create_response, nil
	}
	dkb64Encoded := ""
	if s_login_update_request.Password != "" {
		dkb64Encoded = GetMD5Hash(s_login_update_request.Password)

		wsResult, err := ws.ChangePassword(s_login_credentials, dkb64Encoded)

		if err != nil || wsResult.Result != "0" {
			s_login_create_response.Status = "failed"
			s_login_create_response.StatusCode = 500
			s_login_create_response.ErrorMessage = "Error"

			return s_login_create_response, nil
		}
	}

	fmt.Println(dkb64Encoded)

	db.UpdateUser(dbConn, s_login_credentials.Id, s_login_update_request.Photo, dkb64Encoded)
	//	db.UpdateUser(dbConn, s_login_credentials.Id, s_login_update_request.Name, s_login_update_request.Photo, s_login_update_request.Email, dkb64Encoded, salt)

	return s_login_create_response, nil
}
func CheckPassword(source_password string, hash_password string, salt string) bool {
	if GetMD5Hash(source_password) == hash_password {
		return true
	}
	return false
}
func login(s_login_request s_login_request_hdr) (s_login_response s_login_response_hdr, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("PANIC - ", r)

			err = fmt.Errorf("panic")
		}
	}()

	if s_login_request.Email == "" || s_login_request.Password == "" {
		s_login_response.Status = "failed"
		s_login_response.StatusCode = 403
		s_login_response.ErrorMessage = "Faltando dados"

		return s_login_response, nil
	}

	authToken := random.RandomString(64)

	s_login_response = s_login_response_hdr{}
	s_login_response.Status = "failed"

	dbConn := db.Connect()
	defer dbConn.Close()

	s_login_credentials, err := db.GetLoginInfoByEmail(dbConn, s_login_request.Email)
	if err != nil {
		s_login_response.Status = "failed"
		s_login_response.StatusCode = 403
		s_login_response.ErrorMessage = "Login/Senha inválido"

		return s_login_response, nil
	}
	if s_login_credentials.FailedLoginCount > MAXLOGIN {
		s_login_response.Status = "failed"
		s_login_response.StatusCode = 403
		s_login_response.ErrorMessage = "Conta bloqueada, numero de tentativas excedidas. voce precisa redefinir sua senha."

		return s_login_response, nil
	}
	if s_login_credentials.Status == 0 {
		s_login_response.Status = "failed"
		s_login_response.StatusCode = 403
		s_login_response.ErrorMessage = "Usuario não ativado"

		return s_login_response, nil
	}

	s_login_credentials.Password = GetMD5Hash(s_login_request.Password)
	s_login_credentials.TerminalPassword = GetMD5Hash(s_login_request.Password)

	_, _, err = ws.GetBalance(s_login_credentials)
	if err != nil {
		db.IncreaseFailedLoginOfEmail(dbConn, s_login_request.Email)

		s_login_response.Status = "failed"
		s_login_response.StatusCode = 403
		s_login_response.ErrorMessage = "Login/Senha inválido"
	} else {
		db.ResetFailedLoginOfEmail(dbConn, s_login_request.Email)
		db.InsertToken(dbConn, s_login_credentials.Id, authToken)

		s_login_response.Status = "success"
		s_login_response.Data.AuthToken = authToken
	}

	return s_login_response, nil
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
	if s_request.ActivationCode == "" || s_request.AuthToken == "" {
		result.Status = "failed"
		result.StatusCode = 403
		result.ErrorMessage = "Faltando dados"

		return result, nil
	}

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
		wsResult, err := ws.CreateAccount(s_redis.Name, s_redis.Email, s_redis.Document, s_redis.Cel, s_redis.Password)

		if err != nil {
			result.Status = "failed"
			result.StatusCode = 500
			result.ErrorMessage = "Internal server error"

			return result, nil
		}
		if wsResult.Result != "0" {
			result.Status = "failed"
			result.StatusCode = 500
			result.ErrorMessage = "Internal server error"

			return result, nil

		}
		if wsResult.XMLPerson.CreateAccount.Result != "0" {
			result.Status = "failed"
			result.StatusCode = 500
			result.ErrorMessage = "Internal server error"

			return result, nil

		}

		id, err := db.CreateAccount(dbConn, s_redis.Email, s_redis.Cel, GetMD5Hash(s_redis.Password), "", s_redis.Name, s_redis.Photo, s_redis.Email, GetMD5Hash(s_redis.Password), wsResult.XMLPerson.CreateAccount.PointId)
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
	//	defer func() {
	//		if r := recover(); r != nil {
	//			result.Status = "failed"
	//			result.StatusCode = 404
	//			result.ErrorMessage = "AuthToken invalido"
	//			err = nil
	//		}
	//	}()
	if s_request.AuthToken == "" {
		result.Status = "failed"
		result.StatusCode = 403
		result.ErrorMessage = "Faltando dados"

		return result, nil
	}

	dbConn := db.Connect()

	defer dbConn.Close()

	s_redis := s_redis_create_account_hdr{}

	fmt.Println("AT" + s_request.AuthToken)
	redisString, err := redis.Get(s_request.AuthToken)
	if err != nil {
		panic(err)
	}
	fmt.Println("REDIS" + *redisString)

	err = json.Unmarshal([]byte(*redisString), &s_redis)
	if err != nil {
		panic(err)
	}

	redisBytes, _ := json.Marshal(s_redis)

	redis.Set(s_redis.AuthToken, string(redisBytes), 1000*time.Minute)

	notification.Send(notification.NotificationMessage{"sms", ws.GetTelRAW(s_redis.Cel), "QIWI - Seu codigo de ativação é: " + s_redis.ActivationCode})

	return result, nil
}

func lostPassword(s_lost_password s_lost_password_hdr) {
	dbConn := db.Connect()
	defer dbConn.Close()

	s_user_info, err := db.GetLoginInfoByEmail(dbConn, s_lost_password.Email)

	if err == nil {
		ws.ResetPassword(s_lost_password.Email, s_user_info.TerminalId, s_user_info.TerminalLogin)
	}
	//	recoverToken := random.RandomString(64)

	//	s_redis := s_redis_lost_password_hdr{}
	//	s_redis.Email = s_lost_password.Email
	//	s_redis.RecoverToken = recoverToken

	//	redisString, _ := json.Marshal(s_redis)
	//	redis.Set(s_redis.Email, string(redisString), 1000*time.Minute)

	//	notification.Send(notification.NotificationMessage{"email", s_lost_password.Email, "\n\nhttp://ec2-54-207-24-178.sa-east-1.compute.amazonaws.com/password/#/password/" + s_lost_password.Email + "/" + recoverToken})

}
func changeLPPassword(s_lost_password s_lost_password_hdr) (result s_status) {
	defer func() {
		if r := recover(); r != nil {
			result.Status = "failed"
			result.StatusCode = 0
			result.ErrorMessage = "AuthToken invalido"
		}
	}()
	if s_lost_password.Email == "" || s_lost_password.RecoveryToken == "" || s_lost_password.Password == "" {
		result.Status = "failed"
		result.StatusCode = 403
		result.ErrorMessage = "Faltando dados"

		return result
	}

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
	if s_lost_password.Email == "" || s_lost_password.RecoveryToken == "" {
		result.Status = "failed"
		result.StatusCode = 403
		result.ErrorMessage = "Faltando dados"

		return result
	}

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

//func getPublicLoginInfoByEmail(s_cel_info s_cel_info_hdr) (s_balance_response s_balance_response_hdr, err error) {
//	if s_cel_info.AuthToken == "" || s_cel_info.Cel == "" {
//		s_login_response.Status = "failed"
//		s_login_response.StatusCode = 403
//		s_login_response.ErrorMessage = "Faltando dados"

//		return s_login_response, nil
//	}

//	s_login_info := s_login_info_response_hdr{}

//	dbConn := db.Connect()
//	defer dbConn.Close()

//	s_login_credentials, err := db.GetAuthToken(dbConn, s_cel_info.AuthToken)
//	if err == nil && s_login_credentials.Id > 0 {
//		dbResult, err := db.GetPublicLoginInfoByCel(dbConn, s_cel_info.Cel)
//		if err != nil {
//			s_balance_response.StatusCode = 500
//			s_balance_response.ErrorMessage = "Internal server error"
//		} else {
//			s_balance_response.Status = "success"
//			s_balance_response.StatusCode = 0
//			s_login_info.Data.Cel = s_cel_info.Cel
//			s_login_info.Data.Photo = dbResult.Photo
//			s_login_info.Data.Name = dbResult.Name

//		}
//	} else {
//		s_balance_response.StatusCode = 403
//		s_balance_response.ErrorMessage = "Login/Senha inválido"
//	}
//	return s_balance_response, nil
//}
func getPublicLoginInfoByCel(s_cel_info s_cel_info_hdr) (s_login_info s_login_info_response_hdr, err error) {
	if s_cel_info.AuthToken == "" || s_cel_info.Cel == "" {
		s_login_info.Status = "failed"
		s_login_info.StatusCode = 403
		s_login_info.ErrorMessage = "Faltando dados"

		return s_login_info, nil
	}
	dbConn := db.Connect()
	defer dbConn.Close()

	s_login_credentials, err := db.GetAuthToken(dbConn, s_cel_info.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {
		dbResult, err := db.GetPublicLoginInfoByCel(dbConn, s_cel_info.Cel)
		if err != nil {
			s_login_info.Status = "failed"
			s_login_info.StatusCode = 500
			s_login_info.ErrorMessage = "Internal server error"
		} else {
			s_login_info.Status = "success"
			s_login_info.StatusCode = 0
			s_login_info.Data.Cel = s_cel_info.Cel
			s_login_info.Data.Photo = dbResult.Photo
			s_login_info.Data.Name = dbResult.Name

		}
	} else {
		s_login_info.Status = "failed"
		s_login_info.StatusCode = 403
		s_login_info.ErrorMessage = "Login/Senha inválido"
	}
	return s_login_info, nil
}
func getMyInfo(s_cel_info s_cel_info_hdr) (s_login_info s_login_info_response_hdr, err error) {
	if s_cel_info.AuthToken == "" {
		s_login_info.Status = "failed"
		s_login_info.StatusCode = 403
		s_login_info.ErrorMessage = "Faltando dados"

		return s_login_info, nil
	}

	dbConn := db.Connect()
	defer dbConn.Close()

	s_login_credentials, err := db.GetAuthToken(dbConn, s_cel_info.AuthToken)
	if err == nil && s_login_credentials.Id > 0 {
		userInfo, err := db.GetLoginInfoById(dbConn, s_login_credentials.Id)
		if err != nil {
			s_login_info.StatusCode = 500
			s_login_info.ErrorMessage = "Internal server error"
		} else {
			s_login_info.Status = "success"
			s_login_info.StatusCode = 0
			s_login_info.Data.Cel = userInfo.Cel
			s_login_info.Data.Photo = userInfo.Photo
			s_login_info.Data.Name = userInfo.Name
			s_login_info.Data.Email = userInfo.Email
			s_login_info.Data.Document = userInfo.Document
			s_login_info.Data.TerminalId = userInfo.TerminalId

		}
	} else {
		s_login_info.StatusCode = 403
		s_login_info.ErrorMessage = "Login/Senha inválido"
	}
	return s_login_info, nil
}
