package main

import (
	"fmt"
	"strconv"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/twilio/twilio-go"
	"github.com/twilio/twilio-go/rest/api/v2010"
)

func readDB(seedData []string) {
    sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("TIMEZONE")),
	})
	svc := dynamodb.New(sess)
	for index, entry := range seedData {
		input := &dynamodb.GetItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"Link": {
					S: aws.String(entry),
				},
			},
			TableName: aws.String(os.Getenv("TABLE_NAME")),
		}
		result, err := svc.GetItem(input)
		if err != nil {
			fmt.Println("there was an error: ", err)
		}
		if(len(result.Item) == 0 && entry != "") {
			// alert me
			fmt.Println("new entry found: ", entry)
			alertMe(entry)
			updateDb(entry, strconv.Itoa(index))
		}
		
		fmt.Println("here is your result: ", result.Item)
	}
	
}

func updateDb(entry string, anIndex string) {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("TIMEZONE")),
	})
	svc := dynamodb.New(sess)

	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Link": {
				S: aws.String(entry),
			},
			"TheIndex": {
				S: aws.String(anIndex),
			},
		},
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}

	_, err := svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling UpdateItem: ", err)
	}
}

func alertMe(entry string) {
	fmt.Println("work in progress", entry)
	


	TWILIO_ACCOUNT_SID := os.Getenv("TWILIO_ACCOUNT_SID")
	TWILIO_AUTH_TOKEN := os.Getenv("TWILIO_AUTH_TOKEN")
	if(TWILIO_ACCOUNT_SID == "" || TWILIO_AUTH_TOKEN == "") {
		fmt.Println("you need either an account sid or auth token")
		os.Exit(1)
	}
	client := twilio.NewRestClient()
	params := &openapi.CreateMessageParams{}
	message := "New Hightower! " + entry
	params.SetTo(os.Getenv("PHONE_NUMBER"))
	params.SetFrom(os.Getenv("TWILIO_PHONE"))
	params.SetBody(message)

	_, err := client.ApiV2010.CreateMessage(params)
    if err != nil {
        fmt.Println(err.Error())
    } else {
        fmt.Println("SMS sent successfully!")
    }
}
	
// func seedDB(seedData []string) {
// 	sess, _ := session.NewSession(&aws.Config{
// 		Region: aws.String(Timezone),
// 	})
// 	svc := dynamodb.New(sess)
// 	seedData = seedData[1:]
// 	for index, entry := range seedData {
// 		anIndex := strconv.Itoa(index)
// 		fmt.Println("here is the index stringified: ", anIndex)
// 		input := &dynamodb.PutItemInput{
// 			Item: map[string]*dynamodb.AttributeValue{
// 				"Link": {
// 					S: aws.String(entry),
// 				},
// 				"TheIndex": {
// 					S: aws.String(anIndex),
// 				},
// 			},
// 			ReturnConsumedCapacity: aws.String("TOTAL"),
// 			TableName: aws.String(TableName),
// 		}
// 		result, err := svc.PutItem(input)
// 		if err != nil {
// 			fmt.Println("I have an error", err)
// 			os.Exit(1)
// 		}
// 		fmt.Println(result)
// 	}
// }
