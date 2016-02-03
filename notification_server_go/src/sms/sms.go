package sms

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bitly/go-simplejson"
)

type Data struct {
	SendSmsRequest SendSmsRequestHDR `json:"sendSmsRequest"`
}
type SendSmsRequestHDR struct {
	From           string `json:"from"`
	To             string `json:"to"`
	Msg            string `json:"msg"`
	CallbackOption string `json:"callbackOption"`
	//	AggregateId    string `json:"aggregateId"`
}

func toJson(i interface{}) []byte {
	data, err := json.Marshal(i)
	if err != nil {
		panic("???")
	}
	return data
}
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func Send(tel string, text string) (err error) {
	fmt.Println("SENDING SMS")
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic", r)
			err = fmt.Errorf("pkg: %v", r)
		}
	}()

	url := "https://api-rest.zenvia360.com.br/services/send-sms"

	data := Data{}
	data.SendSmsRequest.From = "QIWI"
	data.SendSmsRequest.To = "55" + tel
	data.SendSmsRequest.Msg = text
	data.SendSmsRequest.CallbackOption = "NONE"
	//	data.SendSmsRequest.AggregateId = "3434"
	fmt.Println("JSON:" + string(toJson(data)))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(toJson(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Add("Authorization", "Basic "+basicAuth("qiwi.api", "M6TQsawZ"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)

	js, err := simplejson.NewJson(body)

	resultStatus := js.Get("sendSmsResponse").Get("statusCode").MustString()

	fmt.Println("response Body:", string(body))
	fmt.Printf("response Status: %s\n", resultStatus)

	if resultStatus != "00" {
		return errors.New("Error")
	}

	return nil

}
