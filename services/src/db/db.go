// main
package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Login_credentials_hdr struct {
	Id               int
	TerminalId       string
	TerminalSerial   string
	TerminalPassword string
	TerminalLogin    string
	Password         string
	PasswordSalt     string
	FailedLoginCount int
	Status           int
}
type Services_hdr struct {
	Id        int
	Name      string
	LongName  string
	ServicoId int
}

func Connect() *sql.DB {
	db, err := sql.Open("postgres", "host=postgres dbname=qiwi sslmode=disable")
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
func CreateAccount(db *sql.DB, email string, cel string, password string, salt string, name string) error {
	stmt, err := db.Prepare(`insert into user_credentials (  email, cel, password, password_salt, name ) values ( $1,$2, $3,$4,$5 )`)
	defer stmt.Close()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	_, err = stmt.Exec(email, cel, password, salt, name)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil

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
func GetLoginInfoByEmail(db *sql.DB, username string) (*Login_credentials_hdr, error) {
	stmt, err := db.Prepare("select  id, password, password_salt, terminal_login, terminal_id , terminal_serial, terminal_password, login_failed_count, status  from user_credentials where email=$1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	s_login_credentials := Login_credentials_hdr{}

	err = stmt.QueryRow(username).Scan(&s_login_credentials.Id, &s_login_credentials.Password, &s_login_credentials.PasswordSalt, &s_login_credentials.TerminalLogin, &s_login_credentials.TerminalId, &s_login_credentials.TerminalSerial, &s_login_credentials.TerminalPassword, &s_login_credentials.FailedLoginCount, &s_login_credentials.Status)

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
	stmt, err := db.Prepare("select  id, terminal_login, terminal_id , terminal_serial ,terminal_password from user_credentials where cel=$1 and password=$2")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	s_login_credentials := Login_credentials_hdr{}
	err = stmt.QueryRow(cel, password).Scan(&s_login_credentials.Id, &s_login_credentials.TerminalLogin, &s_login_credentials.TerminalId, &s_login_credentials.TerminalSerial, &s_login_credentials.TerminalPassword)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &s_login_credentials, nil
}
func GetAuthToken(db *sql.DB, authToken string) (*Login_credentials_hdr, error) {
	s_login_credentials := Login_credentials_hdr{}
	fmt.Println("AUTH" + authToken)

	stmt, err := db.Prepare("select  u.id, u.terminal_login, u.terminal_id , u.terminal_serial,terminal_password from user_tokens t, user_credentials u where t.user_credential_id=u.id and t.token=$1 and u.status=1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(authToken).Scan(&s_login_credentials.Id, &s_login_credentials.TerminalLogin, &s_login_credentials.TerminalId, &s_login_credentials.TerminalSerial, &s_login_credentials.TerminalPassword)

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

	stmt, err := db.Prepare("select i.longName, s.name from servicos s, servicos_items i where i.servico_id=s.id")
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
		err = rows.Scan(&item.LongName, &item.Name)
		if err != nil {
			return nil, err
		}
		fmt.Println("DB " + item.LongName)

		result = append(result, item)
	}
	return &result, nil
}
