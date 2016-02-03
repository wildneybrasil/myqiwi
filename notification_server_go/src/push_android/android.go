package android

import (
	"fmt"
	"github.com/googollee/go-gcm"
)


func Send(token string, title string, message string, params map[string]interface{} )( error ) {
	client := gcm.New("AIzaSyCj0ErrIcy0WqQL5KIYtq_jLmP7V4TCD08")

	
	fmt.Println("GCM: " + token)
	
	
	load := gcm.NewMessage()
	load.AddRecipient(token)
	load.SetPayload("data","alert")
	load.SetPayload("notificationTitle",title)
	load.SetPayload("notificationText",message)
	
	for k,v := range params  {
		load.SetPayload(k,v.(string))
		fmt.Println("PARAMETER: "  + k + " VALUE: "+ v.(string) )
	}
	load.CollapseKey = "Iris"
	load.DelayWhileIdle = true
	load.TimeToLive = 10


	resp, err := client.Send(load)

	fmt.Printf("id: %+v\n", resp)
	fmt.Println("err:", err)
//	fmt.Println("err index:", resp.ErrorIndexes())
//	fmt.Println("reg index:", resp.RefreshIndexes())
	
	return nil;
}
