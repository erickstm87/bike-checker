package main

import (
	"fmt"
	"os"

	"github.com/oliveroneill/exponent-server-sdk-golang/sdk"
)

func push(link string, model string, status string) {
	 // To check the token is valid 
	 pushToken, err := expo.NewExponentPushToken(os.Getenv("EXPOTOKEN"))
	 if err != nil {
		 panic(err)
	 }
 
	 // Create a new Expo SDK client
	 client := expo.NewPushClient(nil)
	 message := ""
	 theTitle := ""
	 if status == "available" {
		message = "New " + model + " available " + link
		theTitle = "New Bike Found! "
	 } else if status == "sold"{
		message = model + " has been sold and deleted from DB" + link
		theTitle = "Bike Sold"
	 }
	 // Publish message
	 response, err := client.Publish(
		 &expo.PushMessage{
			 To: []expo.ExponentPushToken{pushToken},
			 Body: message,
			 Data: map[string]string{"myLink": link},
			 Sound: "default",
			 Title: theTitle,
			 Priority: expo.HighPriority,
		 },
	 )
	 
	 // Check errors
	 if err != nil {
		 panic(err)
	 }
	 
	 // Validate responses
	 if response.ValidateResponse() != nil {
		 fmt.Println(response.PushMessage.To, "failed")
	 }
}