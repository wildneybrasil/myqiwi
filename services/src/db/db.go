// main
package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	_ "github.com/lib/pq"
)

type Login_credentials_hdr struct {
	Id               int
	TerminalId       string
	TerminalSerial   string
	TerminalPassword string
	TerminalLogin    string
	Cel              string
	Photo            string
	Name             string
	Document         string
	Password         string
	Email            string
	PasswordSalt     string
	FailedLoginCount int
	Status           int
}
type Services_hdr struct {
	Id          int
	Name        string
	LongName    string
	PaymentType string
	ServicoId   int
	RvId        string
	Type        string
}
type History_hdr struct {
	Id           int
	Timestamp    int
	PaymentType  string
	CategoryName string
	ServiceName  string
	ServiceId    int
	Rcpt         string
	Amount       string
}

func Connect() *sql.DB {
	db, err := sql.Open("postgres", "user=postgres password=invq1w2e3r4 host=postgres dbname=qiwi sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
func InsertToken(db *sql.DB, userId int, token string) error {
	stmt, err := db.Prepare(`insert into user_tokens (  user_credential_id, token, last_accessed  ) values ( $1,$2, CURRENT_TIMESTAMP )`)
	defer stmt.Close()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	_, err = stmt.Exec(userId, token)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil

}
func InsertPaymentHistory(db *sql.DB, userId int, paymentType string, serviceId string, structRequest interface{}, structResponse interface{}, xmlRequest *string, xmlReponse *string, status int) (int, error) {

	jsonRequest, err := json.Marshal(structRequest)
	if err != nil {
		jsonRequest = []byte("ERROR")
	}
	jsonResponse, err := json.Marshal(structResponse)
	if err != nil {
		jsonResponse = []byte("ERROR")
	}

	stmt, err := db.Prepare(`insert into payment_history ( user_credential_id,payment_type,service_id, json_request, json_response, xml_request, xml_response, status ) values ( $1,$2, $3, $4, $5, $6, $7, $8 )  RETURNING id`)
	defer stmt.Close()

	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}
	id := 0
	err = stmt.QueryRow(userId, paymentType, serviceId, jsonRequest, jsonResponse, *xmlRequest, *xmlReponse, status).Scan(&id)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}
	return id, nil
}

func InsertPaymentHistoryXML(db *sql.DB, xmlRequest string, xmlReponse string, status int) (int, error) {
	stmt, err := db.Prepare(`insert into payment_history ( xml_request, xml_reponse, status ) values ( $1,$2, $3 )  RETURNING id`)
	defer stmt.Close()

	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}
	id := 0
	err = stmt.QueryRow(xmlRequest, xmlReponse, status).Scan(&id)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}
	return id, nil
}
func InsertPaymentHistoryJSON(db *sql.DB, id int, jsonRequest string, jsonReponse string) error {
	stmt, err := db.Prepare(`update  payment_history set json_request=$1, json_response=$2 where id=$3`)
	defer stmt.Close()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	_, err = stmt.Exec(jsonRequest, jsonReponse, id)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func CreateAccount(db *sql.DB, email string, cel string, password string, salt string, name string, photo string, terminal_login string, terminal_password string, terminal_id string) (int, error) {
	stmt, err := db.Prepare(`insert into user_credentials (  email, cel, password, password_salt, name,status,photo, terminal_login, terminal_password, terminal_id ) values ( $1,$2, $3,$4,$5,1,$6, $7,$8, $9 )  RETURNING id`)
	defer stmt.Close()

	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}
	id := 0
	err = stmt.QueryRow(email, cel, password, salt, name, photo, terminal_login, terminal_password, terminal_id).Scan(&id)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}
	return id, nil

}
func VerifyAuth(db *sql.DB, authToken string) int {
	result := 0

	stmt, err := db.Prepare(`select 1 from user_credentials where authToken=$1 and status=1`)
	defer stmt.Close()

	if err != nil {
		fmt.Println(err.Error())
		return 0
	}
	err = stmt.QueryRow(authToken).Scan(&result)

	if err != nil {
		fmt.Println(err.Error())
		return 0
	}

	if result == 1 {
		return 1
	}
	return 0

}
func ResetFailedLoginOfEmail(db *sql.DB, email string) error {
	stmt, err := db.Prepare(`update user_credentials set login_failed_count = 0 where email=$1`)
	defer stmt.Close()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	_, err = stmt.Exec(email)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
func ActivateUser(db *sql.DB, token string) error {
	stmt, err := db.Prepare(`update user_credentials set status = 1 where password_salt=$1`)
	defer stmt.Close()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	_, err = stmt.Exec(token)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil

}
func ChangePassword(db *sql.DB, email string, password string, salt string) error {
	stmt, err := db.Prepare(`update user_credentials set password = $2, password_salt=$3 where email=$1`)
	defer stmt.Close()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	_, err = stmt.Exec(email, password, salt)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil

}

func UpdateUser(db *sql.DB, id int, photo string, password string) error {
	if password != "" {
		fmt.Println("UPDATE WITH PASSWORD: " + strconv.Itoa(id))
		stmt, err := db.Prepare(`update user_credentials set password=$1,  photo=$2, terminal_password=$3 where email=$4`)
		defer stmt.Close()
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		_, err = stmt.Exec(password, photo, password, id)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	} else {
		fmt.Println("UPDATE SIMPLE: " + strconv.Itoa(id))
		stmt, err := db.Prepare(`update user_credentials set photo=$1 where id=$2`)
		defer stmt.Close()
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		_, err = stmt.Exec(photo, id)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

	}
	return nil

}
func IncreaseFailedLoginOfEmail(db *sql.DB, email string) error {
	stmt, err := db.Prepare(`update user_credentials set login_failed_count = login_failed_count + 1 where email=$1`)
	defer stmt.Close()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	_, err = stmt.Exec(email)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil

}
func GetLoginInfoByEmail(db *sql.DB, email string) (*Login_credentials_hdr, error) {
	stmt, err := db.Prepare("select  id, password, password_salt, terminal_login, terminal_id , terminal_serial, terminal_password, login_failed_count, status, cel  from user_credentials where email=$1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	s_login_credentials := Login_credentials_hdr{}

	err = stmt.QueryRow(email).Scan(&s_login_credentials.Id, &s_login_credentials.Password, &s_login_credentials.PasswordSalt, &s_login_credentials.TerminalLogin, &s_login_credentials.TerminalId, &s_login_credentials.TerminalSerial, &s_login_credentials.TerminalPassword, &s_login_credentials.FailedLoginCount, &s_login_credentials.Status, &s_login_credentials.Cel)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &s_login_credentials, nil
}
func GetLoginInfoById(db *sql.DB, id int) (*Login_credentials_hdr, error) {
	stmt, err := db.Prepare("select  id, password, password_salt, terminal_login, terminal_id , terminal_serial, terminal_password, login_failed_count, status, cel, email,name,photo, document  from user_credentials where id=$1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	s_login_credentials := Login_credentials_hdr{}

	var photo []byte
	err = stmt.QueryRow(id).Scan(&s_login_credentials.Id, &s_login_credentials.Password, &s_login_credentials.PasswordSalt, &s_login_credentials.TerminalLogin, &s_login_credentials.TerminalId, &s_login_credentials.TerminalSerial, &s_login_credentials.TerminalPassword, &s_login_credentials.FailedLoginCount, &s_login_credentials.Status, &s_login_credentials.Cel, &s_login_credentials.Email, &s_login_credentials.Name, &photo, &s_login_credentials.Document)

	if len(photo) > 0 {
		s_login_credentials.Photo = string(photo)
	}
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &s_login_credentials, nil
}
func GetLoginInfoByCel(db *sql.DB, cel string) (*Login_credentials_hdr, error) {
	stmt, err := db.Prepare("select  id, password, password_salt, terminal_login, terminal_id , terminal_serial, terminal_password, login_failed_count, status, cel,name  from user_credentials where cel=$1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	s_login_credentials := Login_credentials_hdr{}

	err = stmt.QueryRow(cel).Scan(&s_login_credentials.Id, &s_login_credentials.Password, &s_login_credentials.PasswordSalt, &s_login_credentials.TerminalLogin, &s_login_credentials.TerminalId, &s_login_credentials.TerminalSerial, &s_login_credentials.TerminalPassword, &s_login_credentials.FailedLoginCount, &s_login_credentials.Status, &s_login_credentials.Cel, &s_login_credentials.Name)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &s_login_credentials, nil
}
func GetPublicLoginInfoByCel(db *sql.DB, cel string) (*Login_credentials_hdr, error) {
	stmt, err := db.Prepare("select  id, cel, name, photo  from user_credentials where cel=$1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	s_login_credentials := Login_credentials_hdr{}

	var photo []byte

	s_login_credentials.Photo = string(photo)
	err = stmt.QueryRow(cel).Scan(&s_login_credentials.Id, &s_login_credentials.Cel, &s_login_credentials.Name, &photo)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &s_login_credentials, nil
}
func GetLoginInfoBySalt(db *sql.DB, token string) (*Login_credentials_hdr, error) {
	stmt, err := db.Prepare("select  id, terminal_login, terminal_id , terminal_serial, terminal_password  from user_credentials where password_salt=$1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	s_login_credentials := Login_credentials_hdr{}

	err = stmt.QueryRow(token).Scan(&s_login_credentials.Id, &s_login_credentials.TerminalLogin, &s_login_credentials.TerminalId, &s_login_credentials.TerminalSerial, &s_login_credentials.TerminalPassword)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &s_login_credentials, nil
}
func LoginUsername(db *sql.DB, username string, password string) (*Login_credentials_hdr, error) {
	stmt, err := db.Prepare("select  id, terminal_login, terminal_id , terminal_serial, terminal_password  from user_credentials where email=$1 and password=$2")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	s_login_credentials := Login_credentials_hdr{}

	err = stmt.QueryRow(username, password).Scan(&s_login_credentials.Id, &s_login_credentials.TerminalLogin, &s_login_credentials.TerminalId, &s_login_credentials.TerminalSerial, &s_login_credentials.TerminalPassword)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &s_login_credentials, nil
}
func LoginCel(db *sql.DB, cel string, password string) (*Login_credentials_hdr, error) {
	stmt, err := db.Prepare("select  id, terminal_login, terminal_id , terminal_serial ,terminal_password,cel from user_credentials where cel=$1 and password=$2")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	s_login_credentials := Login_credentials_hdr{}
	err = stmt.QueryRow(cel, password).Scan(&s_login_credentials.Id, &s_login_credentials.TerminalLogin, &s_login_credentials.TerminalId, &s_login_credentials.TerminalSerial, &s_login_credentials.TerminalPassword, &s_login_credentials.Cel)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &s_login_credentials, nil
}
func GetAuthToken(db *sql.DB, authToken string) (*Login_credentials_hdr, error) {
	s_login_credentials := Login_credentials_hdr{}
	fmt.Println("AUTH" + authToken)

	stmt, err := db.Prepare("select u.password_salt, u.password, u.id, u.terminal_login, u.terminal_id , u.terminal_serial,terminal_password ,cel from user_tokens t, user_credentials u where t.user_credential_id=u.id and t.token=$1 and u.status=1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(authToken).Scan(&s_login_credentials.PasswordSalt, &s_login_credentials.Password, &s_login_credentials.Id, &s_login_credentials.TerminalLogin, &s_login_credentials.TerminalId, &s_login_credentials.TerminalSerial, &s_login_credentials.TerminalPassword, &s_login_credentials.Cel)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &s_login_credentials, nil
}

func FindServiceByLongName(services *[]Services_hdr, longName string) *Services_hdr {
	for _, v := range *services {
		fmt.Println("CMPFIND: [" + v.LongName + "] [" + longName + "]")
		if v.LongName == longName {
			fmt.Println("FOUNDDDDDD: " + v.LongName + " " + longName)
			return &v
		}
	}
	return nil
}
func ListServicos(db *sql.DB) (*[]Services_hdr, error) {
	result := make([]Services_hdr, 0)

	stmt, err := db.Prepare("select i.longName, i.rv_id, s.name, s.type from servicos s, servicos_items i where i.servico_id=s.id order by i.longName")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for rows.Next() {
		item := Services_hdr{}
		err = rows.Scan(&item.LongName, &item.RvId, &item.Name, &item.Type)
		if err != nil {
			return nil, err
		}
		fmt.Println("DB " + item.LongName)

		result = append(result, item)
	}
	return &result, nil
}
func GetServiceByPrid(db *sql.DB, prvId int) (*Services_hdr, error) {

	stmt, err := db.Prepare("select i.longName, i.rv_id, s.name, s.type, i.payment_type from servicos s, servicos_items i where i.servico_id=s.id and i.rv_id=$1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	result := Services_hdr{}
	err = stmt.QueryRow(prvId).Scan(&result.LongName, &result.RvId, &result.Name, &result.Type, &result.PaymentType)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &result, nil
}
func ListHistory(db *sql.DB, userId int) (*[]History_hdr, error) {
	result := make([]History_hdr, 0)

	stmt, err := db.Prepare("select i.longname, s.name, extract(epoch from h.tv at time zone 'utc')::integer as tv,h.id, h.payment_type, h.service_id, h.json_request->>'rcpt',h.json_request->>'amount'  from payment_history h, servicos_items i, servicos s where i.rv_id=h.service_id and i.servico_id=s.id and h.user_credential_id=$1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(userId)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for rows.Next() {
		item := History_hdr{}
		err = rows.Scan(&item.ServiceName, &item.CategoryName, &item.Timestamp, &item.Id, &item.PaymentType, &item.ServiceId, &item.Rcpt, &item.Amount)
		if err != nil {
			return nil, err
		}

		result = append(result, item)
	}
	return &result, nil
}
