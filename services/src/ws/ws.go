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
	BillId     int    `xml:"bill-id,omitempty"`
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
	XMLPayment s_XMLPayment `xml:"payment,omitempty"`
}
type s_XMLPaymentFrom struct {
	Amount   string `xml:"amount,attr,omitempty"`
	Currency string `xml:"currency,attr,omitempty"`
}
type s_XMLPaymentTo struct {
	Account  string `xml:"account,attr,omitempty"`
	Amount   string `xml:"amount,attr,omitempty"`
	Currency string `xml:"currency,attr,omitempty"`
	Service  string `xml:"service,attr,omitempty"`
}
type s_XMLPaymentReceipt struct {
	Date string `xml:"date,attr,omitempty"`
	Id   string `xml:"id,attr,omitempty"`
}
type s_XMLPaymentExtras struct {
	Ev_id    string `xml:"ev_id,attr,omitempty"`
	Ev_isWeb string `xml:"ev_isWeb,attr,omitempty"`
}

type s_XMLPayment struct {
	Id                string              `xml:"id,attr,omitempty"`
	XMLPaymentFrom    s_XMLPaymentFrom    `xml:"from,omitempty"`
	XMLPaymentTo      s_XMLPaymentTo      `xml:"to,omitempty"`
	XMLPaymentReceipt s_XMLPaymentReceipt `xml:"receipt,omitempty"`
	XMLPaymentExtras  s_XMLPaymentExtras  `xml:"extras,omitempty"`
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
type s_request_data struct {
	XMLName      xml.Name            `xml:"request"`
	XMLAuth      s_XMLAuth_hdr       `xml:"auth"`
	XMLClient    s_XMLClient_hdr     `xml:"client"`
	XMLAgents    *s_XMLAgents_hdr    `xml:"agents,omitempty"`
	XMLProvider  *s_XMLProvider_hdr  `xml:"providers,omitempty"`
	XMLTerminals *s_XMLTerminals_hdr `xml:"terminals,omitempty"`
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
}
type s_XMLGetProvider_hdr struct {
	Row []s_XMLGetProviderROW_hdr `xml:"row"`
}
type s_XMLProvider_hdr struct {
	XMLGetProvider            *s_XMLGetProvider_hdr        `xml:"getProviders,omitempty"`
	XMLCheckPaymentRequisites *s_XMLCheckPaymentRequisites `xml:"checkPaymentRequisites,omitempty"`
}

type WSResponse_getProvider_hdr struct {
	Result            string            `xml:"result,attr"`
	ResultDescription string            `xml:"result-description,attr"`
	XMLProvider       s_XMLProvider_hdr `xml:"providers"`
}

func send(s_credentials *db.Login_credentials_hdr, request *s_request_data) (*string, error) {
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

	req, err := http.NewRequest("POST", url, buffer)
	fmt.Println(url)
	req.Header.Set("Content-Type", "application/json")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resultString := string(body)

	fmt.Println(resultString)

	return &resultString, nil
}

func GetBalance(s_credentials *db.Login_credentials_hdr) (*string, *string, error) {
	s_response_getBalance := WSResponse_getBalance_hdr{}

	requestType := s_request_data{}
	requestType.XMLAgents = &s_XMLAgents_hdr{}
	requestType.XMLAgents.GetBalance = &s_XMLGetBalance_hdr{}

	//	s_response_getBalance.XMLAgents.GetBalance = s_XMLGetBalance_hdr{}

	result, err := send(s_credentials, &requestType)
	if err != nil {
		return nil, nil, err
	}

	//	fmt.Println(result)
	if err := xml.NewDecoder(strings.NewReader(*result)).Decode(&s_response_getBalance); err != nil {
		return nil, nil, err
	}
	fmt.Println("Codigo " + s_response_getBalance.Result)
	fmt.Println("Balance " + s_response_getBalance.XMLAgents.GetBalance.Balance)

	return &s_response_getBalance.XMLAgents.GetBalance.Balance, &s_response_getBalance.XMLAgents.GetBalance.Overdraft, nil
}
func GetProvider(s_credentials *db.Login_credentials_hdr) (*WSResponse_getProvider_hdr, error) {
	fmt.Println("GET PROVIDER")

	s_response_getProvider := WSResponse_getProvider_hdr{}

	requestType := s_request_data{}
	requestType.XMLProvider = &s_XMLProvider_hdr{}
	requestType.XMLProvider.XMLGetProvider = &s_XMLGetProvider_hdr{}

	result, err := send(s_credentials, &requestType)
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

	result, err := send(s_credentials, &requestType)
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
func GetBillImage(s_credentials *db.Login_credentials_hdr, boletoId int) (*WSResponse_createBill_hdr, error) {
	fmt.Println("GET CREATE BILL")

	s_response_createBill := WSResponse_createBill_hdr{}

	requestType := s_request_data{}
	requestType.XMLAgents = &s_XMLAgents_hdr{}
	requestType.XMLAgents.GetBillImage = &s_XMLGetBillImage_hdr{}
	requestType.XMLAgents.GetBillImage.BillId = boletoId
	requestType.XMLAgents.GetBillImage.FileFormat = "pdf"

	result, err := send(s_credentials, &requestType)

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
func TransferCredits(s_credentials *db.Login_credentials_hdr, toAccount string, toTerminal string, serviceId string, amount string) (*WSResponse_transferCredits_hdr, error) {
	fmt.Println("GET CREATE BILL")

	s_response_createBill := WSResponse_transferCredits_hdr{}

	currentDate := arrow.Now().CFormat("%Y-%m-%dT%H:%M:%S")

	lastId, _ := GetLastID(s_credentials)
	currentId, _ := strconv.ParseInt(lastId.XMLTerminals.XMLGetLastIds.XMLLastPayment.Id, 10, 0)
	currentId = currentId + 1

	requestType := s_request_data{}
	requestType.XMLProvider = &s_XMLProvider_hdr{}
	requestType.XMLProvider.XMLCheckPaymentRequisites = &s_XMLCheckPaymentRequisites{}
	requestType.XMLProvider.XMLCheckPaymentRequisites.XMLPayment.Id = "1"
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

	result, err := send(s_credentials, &requestType)
	if err != nil {
		return nil, err
	}

	//	fmt.Println(result)

	if err := xml.NewDecoder(strings.NewReader(*result)).Decode(&s_response_createBill); err != nil {
		return nil, err
	}

	return &s_response_createBill, nil
}
func GetLastID(s_credentials *db.Login_credentials_hdr) (*WSResponse_lastGetID_hdr, error) {
	fmt.Println("GET CREATE BILL")

	s_response_createBill := WSResponse_lastGetID_hdr{}

	requestType := s_request_data{}
	requestType.XMLTerminals = &s_XMLTerminals_hdr{}
	requestType.XMLTerminals.XMLGetLastIds = &s_XMLGetLastIDS_hdr{}

	result, err := send(s_credentials, &requestType)
	if err != nil {
		return nil, err
	}

	if err := xml.NewDecoder(strings.NewReader(*result)).Decode(&s_response_createBill); err != nil {
		return nil, err
	}

	return &s_response_createBill, nil
}
