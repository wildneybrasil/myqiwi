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
}
type Services_hdr struct {
	Id        int
	Name      string
	LongName  string
	ServicoId int
}

func Connect() *sql.DB {
	db, err := sql.Open("postgres", "host=localhost dbname=qiwi sslmode=disable")
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

	stmt, err := db.Prepare("select  u.id, u.terminal_login, u.terminal_id , u.terminal_serial,terminal_password from user_tokens t, user_credentials u where t.user_credential_id=u.id and t.token=$1")
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
