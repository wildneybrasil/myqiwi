package ws

import (
	"bytes"
	"crypto/tls"
	"db"
	b64 "encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/bmuller/arrow/lib"
)

var (
	url = "https://189.36.23.69:8443/term2/xml.ashx"
)

type s_XMLAuth_hdr struct {
	Login   string `xml:"login,attr"`
	Sign    string `xml:"sign,attr"`
	SignAlg string `xml:"signAlg,attr"`
}
type s_XMLClient_hdr struct {
	Serial   string `xml:"serial,attr"`
	Software string `xml:"software,attr"`
	Terminal string `xml:"terminal,attr"`
}

type s_XMLAgents_hdr struct {
	GetBalance   *s_XMLGetBalance_hdr   `xml:"getBalance,omitempty"`
	CreateBill   *s_XMLCreateBill_hdr   `xml:"createBill,omitempty"`
	GetBillImage *s_XMLGetBillImage_hdr `xml:"getBillImage,omitempty"`
}
type s_XMLBoletoBill struct {
	Amount string `xml:"amount,attr,omitempty"`
	Id     string `xml:"id,attr,omitempty"`
}
type s_XMLCreateBill_hdr struct {
	Amount     string           `xml:"amount,omitempty"`
	BoletoBill *s_XMLBoletoBill `xml:"bill,omitempty"`
}
type s_XMLGetBillImage_hdr struct {
	BillId     string `xml:"bill-id,omitempty"`
	FileFormat string `xml:"file-format,omitempty"`
	Image      string `xml:"image,omitempty"`
	Result     string `xml:"result,attr,omitempty"`
}

type s_XMLGetBalance_hdr struct {
	Result      string `xml:"result,attr,omitempty"`
	Balance     string `xml:"balance,omitempty"`
	TreeBalance string `xml:"tree-balance,omitempty"`
	Overdraft   string `xml:"overdraft,omitempty"`
}
type s_XMLCheckPaymentRequisites struct {
	Result     string       `xml:"result,attr,omitempty"`
	XMLPayment s_XMLPayment `xml:"payment,omitempty"`
}
type s_XMLPaymentFrom struct {
	Amount   string `xml:"amount,attr,omitempty"`
	Currency string `xml:"currency,attr,omitempty"`
	Type     string `xml:"type,attr,omitempty"`
}
type s_XMLPaymentToGoods struct {
	Amount string `xml:"amount,attr,omitempty"`
	Count  string `xml:"count,attr,omitempty"`
	Id     string `xml:"id,attr,omitempty"`
}
type s_XMLPaymentTo struct {
	Account  string               `xml:"account,attr,omitempty"`
	Amount   string               `xml:"amount,attr,omitempty"`
	Currency string               `xml:"currency,attr,omitempty"`
	Service  string               `xml:"service,attr,omitempty"`
	Goods    *s_XMLPaymentToGoods `xml:"goods,omitempty"`
}
type s_XMLPaymentReceipt struct {
	Date string `xml:"date,attr,omitempty"`
	Id   string `xml:"id,attr,omitempty"`
}
type s_XMLPaymentExtras struct {
	Ev_id                string `xml:"ev_id,attr,omitempty"`
	Ev_isWeb             string `xml:"ev_isWeb,attr,omitempty"`
	Ev_ipte              string `xml:"ev_ipte,attr,omitempty"`
	Ev_isnom             string `xml:"ev_isnom,attr,omitempty"`
	Ev_reqtype           string `xml:"ev_reqtype,attr,omitempty"`
	Ev_scan              string `xml:"ev_scan,attr,omitempty"`
	Ev_nsum              string `xml:"ev_nsum,attr,omitempty"`
	Ev_nid               string `xml:"ev_nid,attr,omitempty"`
	Ev_force_amount      string `xml:"ev_force_amount,attr,omitempty"`
	Disp1                string `xml:"disp1,attr,omitempty"`
	Disp2                string `xml:"disp2,attr,omitempty"`
	Disp3                string `xml:"disp3,attr,omitempty"`
	Disp4                string `xml:"disp4,attr,omitempty"`
	PrtData1             string `xml:"prt-data1,attr,omitempty"`
	PrtData2             string `xml:"prt-data2,attr,omitempty"`
	PrtData3             string `xml:"prt-data3,attr,omitempty"`
	Ev_card_number       string `xml:"ev_card_number,attr,omitempty"`
	Ev_step              string `xml:"ev_step,attr,omitempty"`
	Ev_exact_amount      string `xml:"ev_exact_amount,attr,omitempty"`
	Ev_session_guid      string `xml:"ev_session_guid,attr,omitempty"`
	Ev_useExistsVouchers string `xml:"ev_useExistsVouchers,attr,omitempty"`
}
type s_XMLGoodsItemParam struct {
	Name  string `xml:"name,attr,omitempty"`
	Value string `xml:",chardata"`
}
type s_XMLGoodsItem struct {
	Param []s_XMLGoodsItemParam `xml:"param,omitempty"`
}
type s_XMLGoods struct {
	Item     s_XMLGoodsItem `xml:"item,attr,omitempty"`
	Amount   string         `xml:"amount,attr,omitempty"`
	Id       string         `xml:"id,attr,omitempty"`
	Name     string         `xml:"name,attr,omitempty"`
	Currency string         `xml:"currency,attr,omitempty"`
}
type s_XMLGoodsNomenclature struct {
	Amount   string `xml:"amount,attr,omitempty"`
	Id       string `xml:"id,attr,omitempty"`
	Name     string `xml:"name,attr,omitempty"`
	Currency string `xml:"currency,attr,omitempty"`
}
type s_XMLGoodsPayment struct {
	Amount string          `xml:"amount,attr,omitempty"`
	Id     string          `xml:"id,attr,omitempty"`
	Count  string          `xml:"count,attr,omitempty"`
	Item   *s_XMLGoodsItem `xml:"item,omitempty"`
}
type s_XMLVoucher struct {
	Amount      string            `xml:"amount,attr,omitempty"`
	Code        string            `xml:"code,attr,omitempty"`
	printHeader *s_XMLPrintHeader `xml:"printHeader,attr,omitempty"`
}
type s_XMLPrintHeader struct {
	Amount  string         `xml:"amount,attr,omitempty"`
	Code    string         `xml:"code,attr,omitempty"`
	PrvData []s_XMLPrvData `xml:"prvDate,attr,omitempty"`
}
type s_XMLPrvData struct {
	Value   string `xml:"value,attr,omitempty"`
	Content string `prvData,chardata`
}
type s_XMLPayment struct {
	Id                string              `xml:"id,attr,omitempty"`
	Comment           string              `xml:"comment,attr,omitempty"`
	Result            string              `xml:"result,attr,omitempty"`
	Status            string              `xml:"status,attr,omitempty"`
	XMLPaymentFrom    s_XMLPaymentFrom    `xml:"from,omitempty"`
	XMLPaymentTo      s_XMLPaymentTo      `xml:"to,omitempty"`
	XMLPaymentReceipt s_XMLPaymentReceipt `xml:"receipt,omitempty"`
	XMLPaymentExtras  s_XMLPaymentExtras  `xml:"extras,omitempty"`
	XMLVoucher        *s_XMLVoucher       `xml:"voucher,omitempty"`
	XMLGoods          *s_XMLGoodsPayment  `xml:"goods,omitempty"`
}
type s_XMLLastPayment_hdr struct {
	Id            string `xml:"id,attr,omitempty"`
	ReceiptNumber string `xml:"receipt-number,attr,omitempty"`
}
type s_XMLGetLastIDS_hdr struct {
	XMLLastPayment s_XMLLastPayment_hdr `xml:"last-payment,omitempty"`
}
type s_XMLTerminals_hdr struct {
	XMLGetLastIds *s_XMLGetLastIDS_hdr `xml:"getLastIds,omitempty"`
}
type s_XMLChangePassword struct {
	NewPassword string `xml:"new-password,omitempty"`
	Result      string `xml:"result,attr,omitempty"`
}
type s_XMLResetPassword struct {
	Step      string `xml:"step,omitempty"`
	Email     string `xml:"email,omitempty"`
	TermId    string `xml:"term-id,omitempty"`
	TermLogin string `xml:"term-login,omitempty"`
	Result    string `xml:"result,attr,omitempty"`
}
type s_XMLPersons_hdr struct {
	CreateAccount  *s_XMLCreateAccount  `xml:"createAccount"`
	ChangePassword *s_XMLChangePassword `xml:"changePassword"`
	ResetPassword  *s_XMLResetPassword  `xml:"resetPassword"`
}
type s_XMLCreateAccount struct {
	Name        string `xml:"name"`
	Email       string `xml:"e-mail"`
	Document    string `xml:"document"`
	PhoneNumber string `xml:"phone-number"`
	Password    string `xml:"password"`

	DealerId string `xml:"dealer-id,omitempty"`
	Result   string `xml:"result,attr,omitempty"`
	UserId   string `xml:"user-id,omitempty"`
	PointId  string `xml:"point-id,omitempty"`
}
type s_request_data struct {
	XMLName      xml.Name            `xml:"request"`
	XMLAuth      s_XMLAuth_hdr       `xml:"auth"`
	XMLClient    s_XMLClient_hdr     `xml:"client"`
	XMLAgents    *s_XMLAgents_hdr    `xml:"agents,omitempty"`
	XMLProvider  *s_XMLProvider_hdr  `xml:"providers,omitempty"`
	XMLTerminals *s_XMLTerminals_hdr `xml:"terminals,omitempty"`
	XMLPersons   *s_XMLPersons_hdr   `xml:"persons,omitempty"`
}

type WSResponse_getBalance_hdr struct {
	Result            string          `xml:"result,attr"`
	ResultDescription string          `xml:"result-description,attr"`
	XMLAgents         s_XMLAgents_hdr `xml:"agents"`
}
type WSResponse_createBill_hdr struct {
	Result            string          `xml:"result,attr"`
	ResultDescription string          `xml:"result-description,attr"`
	XMLAgents         s_XMLAgents_hdr `xml:"agents"`
}
type WSResponse_createAccount struct {
	Result            string           `xml:"result,attr"`
	ResultDescription string           `xml:"result-description,attr"`
	XMLPerson         s_XMLPersons_hdr `xml:"persons"`
}
type WSResponse_changePassword struct {
	Result            string           `xml:"result,attr"`
	ResultDescription string           `xml:"result-description,attr"`
	XMLPerson         s_XMLPersons_hdr `xml:"persons"`
}
type WSResponse_transferCredits_hdr struct {
	Result            string            `xml:"result,attr"`
	ResultDescription string            `xml:"result-description,attr"`
	XMLProvider       s_XMLProvider_hdr `xml:"providers"`
}
type WSResponse_lastGetID_hdr struct {
	Result            string             `xml:"result,attr"`
	ResultDescription string             `xml:"result-description,attr"`
	XMLTerminals      s_XMLTerminals_hdr `xml:"terminals"`
}

// providers
type s_XMLGetProviderROW_hdr struct {
	FiscalName  string `xml:"fiscal-name,attr"`
	LongName    string `xml:"long-name,attr"`
	PrvId       string `xml:"prv-id,attr"`
	ReceiptName string `xml:"receipt-name,attr"`
	ShortName   string `xml:"short-name,attr"`
	ServiceName string `xml:"service-name,attr, omitempty"`
	Type        string `xml:"type,attr, omitempty"`
}
type s_XMLGetProvider_hdr struct {
	Row []s_XMLGetProviderROW_hdr `xml:"row"`
}
type s_XMLPurchaseOnline struct {
	Result     string       `xml:"result,attr,omitempty"`
	Status     string       `xml:"status,attr,omitempty"`
	XMLPayment s_XMLPayment `xml:"payment,omitempty"`
}
type s_XMLProvider_hdr struct {
	XMLGetProvider            *s_XMLGetProvider_hdr        `xml:"getProviders,omitempty"`
	XMLCheckPaymentRequisites *s_XMLCheckPaymentRequisites `xml:"checkPaymentRequisites,omitempty"`
	XMLPurchaseOnline         *s_XMLPurchaseOnline         `xml:"purchaseOnline,omitempty"`
	XMLGetNomenclature        *s_XMLGetNomenclature        `xml:"getNomenclature,omitempty"`
}

type s_XMLGetNomenclature struct {
	Service string        `xml:"service"`
	Goods   *[]s_XMLGoods `xml:"goods"`
}
type WSResponse_getProvider_hdr struct {
	Result            string            `xml:"result,attr"`
	ResultDescription string            `xml:"result-description,attr"`
	XMLProvider       s_XMLProvider_hdr `xml:"providers"`
}

func send(s_credentials *db.Login_credentials_hdr, request *s_request_data) (*string, *string, error) {
	fmt.Println("SEND")
	buffer := bytes.NewBuffer([]byte{})
	request.XMLAuth.Login = s_credentials.TerminalLogin
	request.XMLAuth.Sign = s_credentials.TerminalPassword
	request.XMLAuth.SignAlg = "MD5"
	request.XMLClient.Serial = s_credentials.TerminalSerial
	request.XMLClient.Software = "X-Snake API v1.1"
	request.XMLClient.Terminal = s_credentials.TerminalId

	enc := xml.NewEncoder(buffer)
	enc.Indent("  ", "    ")

	if err := enc.Encode(request); err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Println(string(buffer.Bytes()))
	requestString := string(buffer.Bytes())

	req, err := http.NewRequest("POST", url, buffer)
	fmt.Println(url)
	req.Header.Set("Content-Type", "application/json")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)

	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	resultString := string(body)

	fmt.Println(resultString)

	return &resultString, &requestString, nil
}
func GetErrorMessage(code string) (string, int) {
	result := ""
	errorCode := 0

	switch code {
	case "3":
		result = "Provedor temporariamente fora de serviço"
		errorCode = 405
		break
	case "150":
		result = "Falha na autenticação"
		errorCode = 401
		break
	case "151":
		result = "Pagamento para este provedor foi negado"
		errorCode = 403
		break
	case "155":
		result = "Não é permitido usar a mesma senha"
		errorCode = 405
		break
	case "202":
		result = "Erro interno, tente novamente"
		errorCode = 405
		break
	case "216":
		result = "Falha ao confirmar transação"
		errorCode = 405
		break
	case "210":
		result = "Transação nao encontrada"
		errorCode = 405
		break
	case "246":
		result = "Terminal esta sendo utilizado por outro computador"
		errorCode = 405
		break
	case "300":
		result = "Erro inexperado"
		errorCode = 405
		break
	case "500":
		result = "Recebimento negado para este pagamento"
		errorCode = 405
		break
	case "700":
		result = "Terminal bloqueado"
		errorCode = 405
		break
	case "701":
		result = "Este servico atingiu o limite diario"
		errorCode = 405
		break
	case "702":
		result = "Este login ja existe"
		errorCode = 405
		break
	case "703":
		result = "Erro ao criar conta ( 703 )"
		errorCode = 405
		break
	case "704":
		result = "Erro ao criar conta ( 704 )"
		errorCode = 405
		break
	case "705":
		result = "Erro ao criar conta ( 705 )"
		errorCode = 405
		break
	case "220":
		result = "Erro na transferencia"
		errorCode = 500
		break
	default:
		result = "Internal server error ( Erro: " + code + ")"
		errorCode = 500
		break
	}
	return result, errorCode
}
func GetBalance(s_credentials *db.Login_credentials_hdr) (*string, *string, error) {
	s_response_getBalance := WSResponse_getBalance_hdr{}

	requestType := s_request_data{}
	requestType.XMLAgents = &s_XMLAgents_hdr{}
	requestType.XMLAgents.GetBalance = &s_XMLGetBalance_hdr{}

	//	s_response_getBalance.XMLAgents.GetBalance = s_XMLGetBalance_hdr{}

	result, _, err := send(s_credentials, &requestType)
	if err != nil {
		return nil, nil, err
	}

	//	fmt.Println(result)
	if err := xml.NewDecoder(strings.NewReader(*result)).Decode(&s_response_getBalance); err != nil {
		return nil, nil, err
	}
	fmt.Println("Codigo " + s_response_getBalance.Result)
	if s_response_getBalance.Result != "0" {

		return nil, nil, fmt.Errorf("Error")
	}
	fmt.Println("Balance " + s_response_getBalance.XMLAgents.GetBalance.Balance)

	return &s_response_getBalance.XMLAgents.GetBalance.Balance, &s_response_getBalance.XMLAgents.GetBalance.Overdraft, nil
}

func CreateAccount(name string, email string, document string, phone string, password string) (*WSResponse_createAccount, error) {
	s_response := WSResponse_createAccount{}

	requestType := s_request_data{}
	requestType.XMLPersons = &s_XMLPersons_hdr{}
	requestType.XMLPersons.CreateAccount = &s_XMLCreateAccount{}
	requestType.XMLPersons.CreateAccount.Name = name
	requestType.XMLPersons.CreateAccount.Email = email
	requestType.XMLPersons.CreateAccount.Document = document
	requestType.XMLPersons.CreateAccount.PhoneNumber = splitCelDDD(phone)
	requestType.XMLPersons.CreateAccount.Password = password

	s_credentials := db.Login_credentials_hdr{}
	s_credentials.TerminalLogin = "ttt2"
	s_credentials.TerminalPassword = "4995EA0596369F512A0334986E824C8A"
	s_credentials.TerminalId = "269"
	s_credentials.TerminalSerial = "2134"

	result, _, err := send(&s_credentials, &requestType)
	if err != nil {
		return nil, err
	}

	//	fmt.Println(result)
	if err := xml.NewDecoder(strings.NewReader(*result)).Decode(&s_response); err != nil {
		return nil, err
	}

	return &s_response, nil
}
func ResetPassword(email string, termId string, termLogin string) (*WSResponse_changePassword, error) {
	s_response := WSResponse_changePassword{}

	requestType := s_request_data{}
	requestType.XMLPersons = &s_XMLPersons_hdr{}
	requestType.XMLPersons.ResetPassword = &s_XMLResetPassword{}
	requestType.XMLPersons.ResetPassword.Email = email
	requestType.XMLPersons.ResetPassword.TermId = termId
	requestType.XMLPersons.ResetPassword.TermLogin = termLogin
	requestType.XMLPersons.ResetPassword.Step = "1"

	s_credentials2 := db.Login_credentials_hdr{}
	s_credentials2.TerminalLogin = "ttt2"
	s_credentials2.TerminalPassword = "4995EA0596369F512A0334986E824C8A"
	s_credentials2.TerminalId = "269"
	s_credentials2.TerminalSerial = "2134"

	result, _, err := send(&s_credentials2, &requestType)
	if err != nil {
		return nil, err
	}

	//	fmt.Println(result)
	if err := xml.NewDecoder(strings.NewReader(*result)).Decode(&s_response); err != nil {
		return nil, err
	}

	return &s_response, nil
}
func ChangePassword(s_credentials *db.Login_credentials_hdr, password string) (*WSResponse_changePassword, error) {
	s_response := WSResponse_changePassword{}

	requestType := s_request_data{}
	requestType.XMLPersons = &s_XMLPersons_hdr{}
	requestType.XMLPersons.ChangePassword = &s_XMLChangePassword{}
	requestType.XMLPersons.ChangePassword.NewPassword = password

	result, _, err := send(s_credentials, &requestType)
	if err != nil {
		return nil, err
	}

	//	fmt.Println(result)
	if err := xml.NewDecoder(strings.NewReader(*result)).Decode(&s_response); err != nil {
		return nil, err
	}

	return &s_response, nil
}
func GetProvider(s_credentials *db.Login_credentials_hdr) (*WSResponse_getProvider_hdr, error) {
	fmt.Println("GET PROVIDER")

	s_response_getProvider := WSResponse_getProvider_hdr{}

	requestType := s_request_data{}
	requestType.XMLProvider = &s_XMLProvider_hdr{}
	requestType.XMLProvider.XMLGetProvider = &s_XMLGetProvider_hdr{}

	result, _, err := send(s_credentials, &requestType)
	if err != nil {
		return nil, err
	}

	//	fmt.Println(result)

	if err := xml.NewDecoder(strings.NewReader(*result)).Decode(&s_response_getProvider); err != nil {
		return nil, err
	}

	fmt.Printf("COUNT %d\n", s_response_getProvider.XMLProvider.XMLGetProvider.Row[0].FiscalName)
	return &s_response_getProvider, nil

}

func CreateBill(s_credentials *db.Login_credentials_hdr, amount string) (*WSResponse_createBill_hdr, error) {
	fmt.Println("GET CREATE BILL")

	s_response_createBill := WSResponse_createBill_hdr{}

	requestType := s_request_data{}
	requestType.XMLAgents = &s_XMLAgents_hdr{}
	requestType.XMLAgents.CreateBill = &s_XMLCreateBill_hdr{}
	requestType.XMLAgents.CreateBill.Amount = amount

	result, _, err := send(s_credentials, &requestType)
	if err != nil {
		return nil, err
	}

	//	fmt.Println(result)

	if err := xml.NewDecoder(strings.NewReader(*result)).Decode(&s_response_createBill); err != nil {
		return nil, err
	}

	fmt.Println("BOLETO ID: " + s_response_createBill.XMLAgents.CreateBill.BoletoBill.Id)
	return &s_response_createBill, nil

}
func GetBillImage(s_credentials *db.Login_credentials_hdr, boletoId string) (*WSResponse_createBill_hdr, error) {
	fmt.Println("GET CREATE BILL")

	s_response_createBill := WSResponse_createBill_hdr{}

	requestType := s_request_data{}
	requestType.XMLAgents = &s_XMLAgents_hdr{}
	requestType.XMLAgents.GetBillImage = &s_XMLGetBillImage_hdr{}
	requestType.XMLAgents.GetBillImage.BillId = boletoId
	requestType.XMLAgents.GetBillImage.FileFormat = "pdf"

	result, _, err := send(s_credentials, &requestType)

	if err != nil {
		return nil, err
	}
	if err := xml.NewDecoder(strings.NewReader(*result)).Decode(&s_response_createBill); err != nil {
		return nil, err
	}

	//	fmt.Println(result)
	binResult := make([]byte, len(s_response_createBill.XMLAgents.GetBillImage.Image))

	_, err = hex.Decode(binResult, []byte(s_response_createBill.XMLAgents.GetBillImage.Image))

	if err != nil {
		return nil, err
	}
	b64Encoded := b64.StdEncoding.EncodeToString([]byte(binResult))
	fmt.Println(b64Encoded)

	s_response_createBill.XMLAgents.GetBillImage.Image = b64Encoded

	return &s_response_createBill, nil
}
func TransferCredits1(s_credentials *db.Login_credentials_hdr, toAccount string, toTerminal string, serviceId string, amount string) (*WSResponse_transferCredits_hdr, error) {
	fmt.Println("GET CREATE BILL")

	s_response_createBill := WSResponse_transferCredits_hdr{}

	currentDate := arrow.Now().CFormat("%Y-%m-%dT%H:%M:%S")

	lastId, _ := GetLastID(s_credentials)
	currentId, _ := strconv.ParseInt(lastId.XMLTerminals.XMLGetLastIds.XMLLastPayment.Id, 10, 0)
	currentId = currentId + 1

	requestType := s_request_data{}
	requestType.XMLProvider = &s_XMLProvider_hdr{}
	requestType.XMLProvider.XMLCheckPaymentRequisites = &s_XMLCheckPaymentRequisites{}
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Id = strconv.Itoa(int(currentId))
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentFrom.Amount = amount
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentFrom.Currency = "986"
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentTo.Amount = amount
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentTo.Currency = "986"
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentTo.Account = toAccount
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentTo.Service = "15695"
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentReceipt.Date = currentDate
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentReceipt.Id = strconv.Itoa(int(currentId))
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Ev_id = toTerminal
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Ev_isWeb = "1"

	result, _, err := send(s_credentials, &requestType)
	if err != nil {
		return nil, err
	}

	//	fmt.Println(result)

	if err := xml.NewDecoder(strings.NewReader(*result)).Decode(&s_response_createBill); err != nil {
		return nil, err
	}

	return &s_response_createBill, nil
}
func splitCelDDD(phone string) string {
	ddd := phone[0:2]
	cel := phone[2:]

	return "(" + ddd + ")" + cel
}
func DoPaymentTel1(s_credentials *db.Login_credentials_hdr, toAccount string, serviceId string, Ev_force string, Ev_step string, Ev_isnom string, Ev_reqtype string) (*WSResponse_transferCredits_hdr, error) {
	fmt.Println("GET CREATE BILL")

	s_response_createBill := WSResponse_transferCredits_hdr{}

	currentDate := arrow.Now().CFormat("%Y-%m-%dT%H:%M:%S")

	lastId, _ := GetLastID(s_credentials)
	currentId, _ := strconv.ParseInt(lastId.XMLTerminals.XMLGetLastIds.XMLLastPayment.Id, 10, 0)
	currentId = currentId + 1

	requestType := s_request_data{}
	requestType.XMLProvider = &s_XMLProvider_hdr{}
	requestType.XMLProvider.XMLCheckPaymentRequisites = &s_XMLCheckPaymentRequisites{}
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Id = strconv.Itoa(int(currentId))
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentFrom.Amount = "1.0"
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentFrom.Currency = "986"
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentTo.Amount = "1.0"
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentTo.Currency = "986"
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentTo.Account = splitCelDDD(toAccount)
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentTo.Service = serviceId
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentReceipt.Date = currentDate
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentReceipt.Id = strconv.Itoa(int(currentId))
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Ev_force_amount = Ev_force
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Ev_isnom = Ev_isnom
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Ev_reqtype = Ev_reqtype
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Ev_isWeb = "1"

	result, _, err := send(s_credentials, &requestType)
	if err != nil {
		return nil, err
	}

	//	fmt.Println(result)

	if err := xml.NewDecoder(strings.NewReader(*result)).Decode(&s_response_createBill); err != nil {
		return nil, err
	}

	return &s_response_createBill, nil
}
func DoPaymentTel2(s_credentials *db.Login_credentials_hdr, currentId string, toAccount string, serviceId string, selectedAmount string) (*WSResponse_transferCredits_hdr, error) {
	fmt.Println("GET CREATE BILL")

	s_response_createBill := WSResponse_transferCredits_hdr{}

	currentDate := arrow.Now().CFormat("%Y-%m-%dT%H:%M:%S")

	requestType := s_request_data{}
	requestType.XMLProvider = &s_XMLProvider_hdr{}
	requestType.XMLProvider.XMLCheckPaymentRequisites = &s_XMLCheckPaymentRequisites{}
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Id = currentId
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentFrom.Amount = selectedAmount
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentFrom.Currency = "986"
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentTo.Amount = selectedAmount
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentTo.Currency = "986"
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentTo.Account = splitCelDDD(toAccount)
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentTo.Service = serviceId
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentReceipt.Date = currentDate
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentReceipt.Id = currentId
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Ev_exact_amount = selectedAmount
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Ev_isWeb = "1"
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Ev_nid = "1"
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Ev_nsum = selectedAmount
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Ev_reqtype = "2"

	result, _, err := send(s_credentials, &requestType)
	if err != nil {
		return nil, err
	}

	//	fmt.Println(result)

	if err := xml.NewDecoder(strings.NewReader(*result)).Decode(&s_response_createBill); err != nil {
		return nil, err
	}

	return &s_response_createBill, nil
}
func DoPaymentGames2(s_credentials *db.Login_credentials_hdr, toAccount string, serviceId string, selectedAmount string, goodsId string) (*WSResponse_transferCredits_hdr, *string, *string, error) {
	fmt.Println("GET CREATE BILL")

	lastId, _ := GetLastID(s_credentials)
	currentId, _ := strconv.ParseInt(lastId.XMLTerminals.XMLGetLastIds.XMLLastPayment.Id, 10, 0)
	currentId = currentId + 1

	s_response_createBill := WSResponse_transferCredits_hdr{}

	currentDate := arrow.Now().CFormat("%Y-%m-%dT%H:%M:%S")

	requestType := s_request_data{}
	requestType.XMLProvider = &s_XMLProvider_hdr{}
	requestType.XMLProvider.XMLPurchaseOnline = &s_XMLPurchaseOnline{}
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLGoods = &s_XMLGoodsPayment{}
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentTo.Goods = &s_XMLPaymentToGoods{}

	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentFrom.Amount = selectedAmount
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentFrom.Currency = "986"
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.Comment = "Comment"

	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentTo.Amount = selectedAmount
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentTo.Currency = "986"
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentTo.Account = toAccount
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentTo.Service = serviceId
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentReceipt.Date = currentDate
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.Id = strconv.Itoa(int(currentId))
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentExtras.Ev_isWeb = "1"
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentExtras.Ev_useExistsVouchers = "true"
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentTo.Goods.Amount = selectedAmount
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentTo.Goods.Count = "1"
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentTo.Goods.Id = goodsId
	result, requestString, err := send(s_credentials, &requestType)
	if err != nil {
		return nil, nil, nil, err
	}

	//	fmt.Println(result)

	s_response_createBill.XMLProvider.XMLPurchaseOnline = &s_XMLPurchaseOnline{}
	s_response_createBill.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLVoucher = &s_XMLVoucher{}
	s_response_createBill.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLGoods = &s_XMLGoodsPayment{}
	s_response_createBill.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLGoods.Item = &s_XMLGoodsItem{}

	if err := xml.NewDecoder(strings.NewReader(*result)).Decode(&s_response_createBill); err != nil {
		return nil, nil, nil, err
	}

	return &s_response_createBill, requestString, result, nil
}
func DoPaymentTel3(s_credentials *db.Login_credentials_hdr, currentId string, session string, toAccount string, serviceId string, selectedAmount string) (*WSResponse_transferCredits_hdr, *string, *string, error) {
	fmt.Println("GET CREATE BILL")

	s_response_createBill := WSResponse_transferCredits_hdr{}

	currentDate := arrow.Now().CFormat("%Y-%m-%dT%H:%M:%S")

	if currentId == "" {
		lastId, _ := GetLastID(s_credentials)
		icurrentId, _ := strconv.ParseInt(lastId.XMLTerminals.XMLGetLastIds.XMLLastPayment.Id, 10, 0)
		icurrentId = icurrentId + 1
		currentId = strconv.Itoa(int(icurrentId))

	}
	requestType := s_request_data{}
	requestType.XMLProvider = &s_XMLProvider_hdr{}
	requestType.XMLProvider.XMLPurchaseOnline = &s_XMLPurchaseOnline{}
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.Id = currentId
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentFrom.Amount = selectedAmount
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentFrom.Currency = "986"
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.Comment = "Comment"

	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentTo.Amount = selectedAmount
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentTo.Currency = "986"
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentFrom.Type = "2"
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentTo.Account = toAccount
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentTo.Service = serviceId
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentReceipt.Date = currentDate
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentReceipt.Id = currentId
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentExtras.Ev_exact_amount = selectedAmount
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentExtras.Ev_isWeb = "1"
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentExtras.Ev_nid = "1"
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentExtras.Ev_nsum = selectedAmount
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentExtras.Ev_reqtype = "2"
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentExtras.Ev_useExistsVouchers = "true"
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentExtras.Ev_session_guid = session

	result, requestString, err := send(s_credentials, &requestType)
	if err != nil {
		return nil, nil, nil, err
	}

	//	fmt.Println(result)

	if err := xml.NewDecoder(strings.NewReader(*result)).Decode(&s_response_createBill); err != nil {
		return nil, nil, nil, err
	}

	return &s_response_createBill, requestString, result, nil
}
func DoPaymentTrans1(s_credentials *db.Login_credentials_hdr, cardNumber string, serviceId string, Ev_force string, Ev_step string, Ev_isnom string, Ev_reqtype string) (*WSResponse_transferCredits_hdr, error) {
	fmt.Println("GET CREATE BILL")

	s_response_createBill := WSResponse_transferCredits_hdr{}

	currentDate := arrow.Now().CFormat("%Y-%m-%dT%H:%M:%S")

	lastId, _ := GetLastID(s_credentials)
	currentId, _ := strconv.ParseInt(lastId.XMLTerminals.XMLGetLastIds.XMLLastPayment.Id, 10, 0)
	currentId = currentId + 1

	requestType := s_request_data{}
	requestType.XMLProvider = &s_XMLProvider_hdr{}
	requestType.XMLProvider.XMLCheckPaymentRequisites = &s_XMLCheckPaymentRequisites{}
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Id = strconv.Itoa(int(currentId))
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentFrom.Amount = "1.0"
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentFrom.Currency = "986"
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentTo.Amount = "1.0"
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentTo.Currency = "986"
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentTo.Account = "(99)999999999"
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentTo.Service = serviceId
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentReceipt.Date = currentDate
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentReceipt.Id = strconv.Itoa(int(currentId))
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Ev_step = Ev_step
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Ev_isWeb = "1"
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Ev_reqtype = Ev_reqtype
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Ev_isnom = Ev_isnom
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Ev_card_number = cardNumber

	result, _, err := send(s_credentials, &requestType)
	if err != nil {
		return nil, err
	}

	//	fmt.Println(result)

	if err := xml.NewDecoder(strings.NewReader(*result)).Decode(&s_response_createBill); err != nil {
		return nil, err
	}

	return &s_response_createBill, nil
}
func DoPaymentGames1(s_credentials *db.Login_credentials_hdr, serviceId string) (*WSResponse_transferCredits_hdr, error) {
	fmt.Println("GET CREATE BILL")

	s_response_createBill := WSResponse_transferCredits_hdr{}

	requestType := s_request_data{}
	requestType.XMLProvider = &s_XMLProvider_hdr{}
	requestType.XMLProvider.XMLGetNomenclature = &s_XMLGetNomenclature{}
	requestType.XMLProvider.XMLGetNomenclature.Service = serviceId

	result, _, err := send(s_credentials, &requestType)
	if err != nil {
		return nil, err
	}

	//	fmt.Println(result)

	if err := xml.NewDecoder(strings.NewReader(*result)).Decode(&s_response_createBill); err != nil {
		return nil, err
	}

	return &s_response_createBill, nil
}
func TransferCredits2(s_credentials *db.Login_credentials_hdr, id string, session string, toAccount string, toTerminal string, serviceId string, amount string) (*WSResponse_transferCredits_hdr, *string, *string, error) {
	fmt.Println("GET CREATE BILL")

	s_response_createBill := WSResponse_transferCredits_hdr{}

	currentDate := arrow.Now().CFormat("%Y-%m-%dT%H:%M:%S")

	requestType := s_request_data{}
	requestType.XMLProvider = &s_XMLProvider_hdr{}
	requestType.XMLProvider.XMLPurchaseOnline = &s_XMLPurchaseOnline{}
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.Id = id
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.Comment = "comment"
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentFrom.Amount = amount
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentFrom.Currency = "986"
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentTo.Amount = amount
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentTo.Currency = "986"
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentTo.Account = toAccount
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentTo.Service = serviceId
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentReceipt.Date = currentDate
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentReceipt.Id = id
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentExtras.Ev_isWeb = "1"
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentExtras.Ev_exact_amount = amount
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentExtras.Ev_session_guid = session
	requestType.XMLProvider.XMLPurchaseOnline.XMLPayment.XMLPaymentExtras.Ev_useExistsVouchers = "true"

	result, requestXML, err := send(s_credentials, &requestType)
	if err != nil {
		return nil, nil, nil, err
	}

	//	fmt.Println(result)

	s_response_createBill.XMLProvider.XMLPurchaseOnline = &s_XMLPurchaseOnline{}

	if err := xml.NewDecoder(strings.NewReader(*result)).Decode(&s_response_createBill); err != nil {
		return nil, nil, nil, err
	}

	return &s_response_createBill, requestXML, result, nil
}
func GetLastID(s_credentials *db.Login_credentials_hdr) (*WSResponse_lastGetID_hdr, error) {
	fmt.Println("GET CREATE BILL")

	s_response_createBill := WSResponse_lastGetID_hdr{}

	requestType := s_request_data{}
	requestType.XMLTerminals = &s_XMLTerminals_hdr{}
	requestType.XMLTerminals.XMLGetLastIds = &s_XMLGetLastIDS_hdr{}

	result, _, err := send(s_credentials, &requestType)
	if err != nil {
		return nil, err
	}

	if err := xml.NewDecoder(strings.NewReader(*result)).Decode(&s_response_createBill); err != nil {
		return nil, err
	}

	return &s_response_createBill, nil
}
func GetBoletoInfo(s_credentials *db.Login_credentials_hdr, boleto string, amount string, fromAccount string, scanned string, service string) (*WSResponse_transferCredits_hdr, error) {
	fmt.Println("GET CREATE BILL")

	s_response_createBill := WSResponse_transferCredits_hdr{}

	currentDate := arrow.Now().CFormat("%Y-%m-%dT%H:%M:%S")

	lastId, _ := GetLastID(s_credentials)
	currentId, _ := strconv.ParseInt(lastId.XMLTerminals.XMLGetLastIds.XMLLastPayment.Id, 10, 0)
	currentId = currentId + 1

	requestType := s_request_data{}
	requestType.XMLProvider = &s_XMLProvider_hdr{}
	requestType.XMLProvider.XMLCheckPaymentRequisites = &s_XMLCheckPaymentRequisites{}
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Id = strconv.Itoa(int(currentId))
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentFrom.Amount = "1"
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentFrom.Currency = "986"
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentTo.Amount = "1"
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentTo.Currency = "986"
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentTo.Account = fromAccount
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentTo.Service = service
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentReceipt.Date = currentDate
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentReceipt.Id = strconv.Itoa(int(currentId))
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Ev_isWeb = "1"
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Ev_ipte = boleto
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Ev_isnom = "1"
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Ev_reqtype = "2"
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Ev_scan = scanned
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Ev_exact_amount = amount

	result, _, err := send(s_credentials, &requestType)
	if err != nil {
		return nil, err
	}

	//	fmt.Println(result)

	if err := xml.NewDecoder(strings.NewReader(*result)).Decode(&s_response_createBill); err != nil {
		return nil, err
	}

	return &s_response_createBill, nil
}
func ChangeBoletoValue(s_credentials *db.Login_credentials_hdr, session string, id string, amount string, fromAccount string) (*WSResponse_transferCredits_hdr, error) {
	fmt.Println("GET CREATE BILL")

	s_response_createBill := WSResponse_transferCredits_hdr{}

	currentDate := arrow.Now().CFormat("%Y-%m-%dT%H:%M:%S")

	requestType := s_request_data{}
	requestType.XMLProvider = &s_XMLProvider_hdr{}
	requestType.XMLProvider.XMLCheckPaymentRequisites = &s_XMLCheckPaymentRequisites{}
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Id = id
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentFrom.Amount = amount
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentFrom.Currency = "986"
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentTo.Amount = "1"
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentTo.Currency = "986"
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentTo.Account = fromAccount
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentTo.Service = "151022"
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentReceipt.Date = currentDate
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentReceipt.Id = id
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Ev_isWeb = "1"
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Ev_exact_amount = amount
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Ev_session_guid = session
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.XMLPaymentExtras.Ev_reqtype = "2"

	result, _, err := send(s_credentials, &requestType)
	if err != nil {
		return nil, err
	}

	//	fmt.Println(result)

	if err := xml.NewDecoder(strings.NewReader(*result)).Decode(&s_response_createBill); err != nil {
		return nil, err
	}

	return &s_response_createBill, nil
}
