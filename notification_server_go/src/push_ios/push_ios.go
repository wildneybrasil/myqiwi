package ios

import (
	"fmt"
	apns "github.com/anachronistic/apns"

)

func Send(token string, badge int, message string, params map[string]interface{} ) error {
	payload := apns.NewPayload()
	payload.Alert = message
	payload.Badge = badge


	pn := apns.NewPushNotification()
	pn.DeviceToken = token
	pn.AddPayload(payload)

	for k,v := range params  {
		pn.Set( k, v )
		fmt.Println("PARAMETER: "  + k + " VALUE: "+ v.(string) )
	}
	fmt.Println("Sending Apple Push TO: " + token)

	client := apns.NewClient("gateway.push.apple.com:2195", "/IE/Iris/cert/ios_cert.pem", "/IE/Iris/cert/ios_key.pem")
	resp := client.Send(pn)

	pn.PayloadString()

//	fmt.Println("Sending Apple Push TO: " + token + " Result " + resp.Error.Error())

//	if( !Success ){
//		return errors.New( "Error" )
//	}
	return resp.Error;
}
